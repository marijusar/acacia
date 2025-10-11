package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"acacia/packages/auth"
	"acacia/packages/db"
	"acacia/packages/httperr"
	"acacia/packages/schemas"

	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UsersController struct {
	queries    *db.Queries
	logger     *logrus.Logger
	validate   *validator.Validate
	jwtManager *auth.JWTManager
}

func NewUsersController(queries *db.Queries, logger *logrus.Logger, jwtManager *auth.JWTManager) *UsersController {
	return &UsersController{
		queries:    queries,
		logger:     logger,
		validate:   validator.New(),
		jwtManager: jwtManager,
	}
}

func (c *UsersController) Register(w http.ResponseWriter, r *http.Request) error {
	var req schemas.RegisterUserInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return httperr.WithStatus(errors.New("Invalid JSON"), http.StatusBadRequest)
	}

	// Validate input
	if err := c.validate.Struct(req); err != nil {
		return httperr.WithStatus(schemas.HandleUserValidationErrors(err), http.StatusBadRequest)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.logger.WithError(err).Error("Failed to hash password")
		return httperr.WithStatus(errors.New("Internal server error"), http.StatusInternalServerError)
	}

	// Create user
	params := db.CreateUserParams{
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: string(hashedPassword),
	}

	user, err := c.queries.CreateUser(r.Context(), params)
	if err != nil {
		// Check for duplicate email
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == db.PgErrUniqueViolation {
				return httperr.WithStatus(errors.New("Email already registered"), http.StatusConflict)
			}
		}
		c.logger.WithError(err).Error("Failed to create user")
		return httperr.WithStatus(errors.New("Internal server error"), http.StatusInternalServerError)
	}

	// Return user without password hash
	response := schemas.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
	return nil
}

func (c *UsersController) Login(w http.ResponseWriter, r *http.Request) error {
	var req schemas.LoginUserInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return httperr.WithStatus(errors.New("Invalid JSON"), http.StatusBadRequest)
	}

	// Validate input
	if err := c.validate.Struct(req); err != nil {
		return httperr.WithStatus(schemas.HandleUserValidationErrors(err), http.StatusBadRequest)
	}

	// Get user by email
	user, err := c.queries.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperr.WithStatus(errors.New("Invalid email or password"), http.StatusUnauthorized)
		}
		c.logger.WithError(err).Error("Failed to get user by email")
		return httperr.WithStatus(errors.New("Internal server error"), http.StatusInternalServerError)
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return httperr.WithStatus(errors.New("Invalid email or password"), http.StatusUnauthorized)
	}

	// Generate JTI for refresh token
	jti, err := auth.GenerateJTI()
	if err != nil {
		c.logger.WithError(err).Error("Failed to generate JTI")
		return httperr.WithStatus(errors.New("Internal server error"), http.StatusInternalServerError)
	}

	// Generate tokens
	accessToken, accessExpiresAt, err := c.jwtManager.GenerateAccessToken(user.ID)
	if err != nil {
		c.logger.WithError(err).Error("Failed to generate access token")
		return httperr.WithStatus(errors.New("Internal server error"), http.StatusInternalServerError)
	}

	refreshToken, refreshExpiresAt, err := c.jwtManager.GenerateRefreshToken(user.ID, jti)
	if err != nil {
		c.logger.WithError(err).Error("Failed to generate refresh token")
		return httperr.WithStatus(errors.New("Internal server error"), http.StatusInternalServerError)
	}

	// Store refresh token in database
	_, err = c.queries.CreateRefreshToken(r.Context(), db.CreateRefreshTokenParams{
		UserID:    user.ID,
		Jti:       jti,
		ExpiresAt: refreshExpiresAt,
	})
	if err != nil {
		c.logger.WithError(err).Error("Failed to store refresh token")
		return httperr.WithStatus(errors.New("Internal server error"), http.StatusInternalServerError)
	}

	// Calculate cookie MaxAge from expiration times
	accessMaxAge := int(time.Until(accessExpiresAt).Seconds())
	refreshMaxAge := int(time.Until(refreshExpiresAt).Seconds())

	// Set cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "access-token",
		Value:    accessToken,
		HttpOnly: true,
		Secure:   true, // Set to true in production with HTTPS
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   accessMaxAge,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh-token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true, // Set to true in production with HTTPS
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   refreshMaxAge,
	})

	// Return user info
	response := schemas.LoginResponse{
		User: schemas.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
		},
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	return nil
}
