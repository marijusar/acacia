package auth

import (
	"acacia/packages/db"
	"acacia/packages/httperr"
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type contextKey string

const UserIDKey contextKey = "user_id"

var (
	ErrMissingToken = errors.New("Unauthorized")
)

// RefreshTokenStore defines the interface for refresh token database operations
type RefreshTokenStore interface {
	GetRefreshTokenByUserAndJTI(ctx context.Context, arg db.GetRefreshTokenByUserAndJTIParams) (db.RefreshToken, error)
}

// AuthMiddleware validates JWT tokens and attaches user ID to request context
type AuthMiddleware struct {
	jwtManager *JWTManager
	queries    RefreshTokenStore
	logger     *logrus.Logger
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(jwtManager *JWTManager, queries RefreshTokenStore, logger *logrus.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
		queries:    queries,
		logger:     logger,
	}
}

// authenticateRequest handles the authentication logic and returns an error if authentication fails
func (m *AuthMiddleware) authenticateRequest(w http.ResponseWriter, r *http.Request, next http.Handler) error {
	// Try to get access token from cookie
	accessCookie, err := r.Cookie("access-token")
	if err == nil && accessCookie.Value != "" {
		// Validate access token
		claims, err := m.jwtManager.ValidateToken(accessCookie.Value)
		if err == nil {
			// Access token is valid, attach user ID to context and proceed
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
			return nil
		}

		// Access token is invalid or expired, continue to check refresh token
		if !errors.Is(err, ErrExpiredToken) {
			m.logger.WithError(err).Debug("Invalid access token")
		}
	}

	// Access token is missing or invalid, try refresh token
	refreshCookie, err := r.Cookie("refresh-token")
	if err != nil || refreshCookie.Value == "" {
		// No refresh token available
		return httperr.WithStatus(ErrMissingToken, http.StatusUnauthorized)
	}

	// Validate refresh token
	refreshClaims, err := m.jwtManager.ValidateToken(refreshCookie.Value)
	if err != nil {
		// Refresh token is invalid or expired
		if errors.Is(err, ErrExpiredToken) {
			m.logger.Debug("Refresh token expired")
		} else {
			m.logger.WithError(err).Debug("Invalid refresh token")
		}
		return httperr.WithStatus(ErrMissingToken, http.StatusUnauthorized)
	}

	// Check if refresh token exists in database and is not revoked
	_, err = m.queries.GetRefreshTokenByUserAndJTI(r.Context(), db.GetRefreshTokenByUserAndJTIParams{
		UserID: refreshClaims.UserID,
		Jti:    refreshClaims.JTI,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			m.logger.Debug("Refresh token not found or revoked")
		} else {
			m.logger.WithError(err).Error("Failed to validate refresh token in database")
		}
		return httperr.WithStatus(ErrMissingToken, http.StatusUnauthorized)
	}

	// Refresh token is valid, generate new access token
	newAccessToken, accessExpiresAt, err := m.jwtManager.GenerateAccessToken(refreshClaims.UserID)
	if err != nil {
		m.logger.WithError(err).Error("Failed to generate new access token")
		return httperr.WithStatus(errors.New("Internal server error"), http.StatusInternalServerError)
	}

	// Set new access token cookie
	accessMaxAge := int(time.Until(accessExpiresAt).Seconds())
	http.SetCookie(w, &http.Cookie{
		Name:     "access-token",
		Value:    newAccessToken,
		HttpOnly: true,
		Secure:   true, // Set to true in production with HTTPS
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   accessMaxAge,
	})

	// Attach user ID to context and proceed
	ctx := context.WithValue(r.Context(), UserIDKey, refreshClaims.UserID)
	next.ServeHTTP(w, r.WithContext(ctx))
	return nil
}

// Handle is the middleware handler that works with chi routers
func (m *AuthMiddleware) Handle(next http.Handler) http.Handler {
	return httperr.WithMiddlewareErrorHandler(func(w http.ResponseWriter, r *http.Request, next http.Handler) error {
		return m.authenticateRequest(w, r, next)
	})(next)
}

// GetUserID retrieves the user ID from the request context
func GetUserID(r *http.Request) (int64, bool) {
	userID, ok := r.Context().Value(UserIDKey).(int64)
	return userID, ok
}
