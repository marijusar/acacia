package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"acacia/packages/db"
	"acacia/packages/schemas"
	"acacia/packages/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestRegisterUser(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("should register user successfully", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()

		// Prepare request body
		registerReq := schemas.RegisterUserInput{
			Email:    "test@example.com",
			Name:     "Test User",
			Password: "password123",
		}
		reqBody, _ := json.Marshal(registerReq)

		// Make HTTP request to the server
		url := fmt.Sprintf("%s/users/register", setup.Server.GetURL())
		resp, _ := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
		defer resp.Body.Close()

		// Assert response status
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// Parse response body
		var userResp schemas.UserResponse
		_ = json.NewDecoder(resp.Body).Decode(&userResp)

		// Assert response body
		assert.NotZero(t, userResp.ID)
		assert.Equal(t, "test@example.com", userResp.Email)
		assert.Equal(t, "Test User", userResp.Name)
		assert.NotZero(t, userResp.CreatedAt)

		// Verify password was hashed (fetch from DB)
		var user db.User
		_ = setup.DB.DB.QueryRowContext(ctx, "SELECT * FROM users WHERE email = $1", "test@example.com").
			Scan(&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)

		// Verify password hash
		err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("password123"))
		assert.NoError(t, err, "Password should be properly hashed")
	})

	t.Run("should return 400 for invalid JSON", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()


		// Make HTTP request with invalid JSON
		url := fmt.Sprintf("%s/users/register", setup.Server.GetURL())
		resp, _ := http.Post(url, "application/json", bytes.NewBufferString("{invalid json"))
		defer resp.Body.Close()

		// Assert response status
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for invalid email", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()


		// Prepare request with invalid email
		registerReq := schemas.RegisterUserInput{
			Email:    "invalid-email",
			Name:     "Test User",
			Password: "password123",
		}
		reqBody, _ := json.Marshal(registerReq)

		url := fmt.Sprintf("%s/users/register", setup.Server.GetURL())
		resp, _ := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// Verify error message
		var errResp map[string]string
		json.NewDecoder(resp.Body).Decode(&errResp)
		assert.Contains(t, errResp["message"], "email")
	})

	t.Run("should return 400 for missing email", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()


		// Prepare request without email
		registerReq := schemas.RegisterUserInput{
			Email:    "",
			Name:     "Test User",
			Password: "password123",
		}
		reqBody, _ := json.Marshal(registerReq)

		url := fmt.Sprintf("%s/users/register", setup.Server.GetURL())
		resp, _ := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for empty name", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()


		// Prepare request with empty name
		registerReq := schemas.RegisterUserInput{
			Email:    "test@example.com",
			Name:     "",
			Password: "password123",
		}
		reqBody, _ := json.Marshal(registerReq)

		url := fmt.Sprintf("%s/users/register", setup.Server.GetURL())
		resp, _ := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for name too long", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()


		// Prepare request with name exceeding 100 characters
		longName := strings.Repeat("a", 101)
		registerReq := schemas.RegisterUserInput{
			Email:    "test@example.com",
			Name:     longName,
			Password: "password123",
		}
		reqBody, _ := json.Marshal(registerReq)

		url := fmt.Sprintf("%s/users/register", setup.Server.GetURL())
		resp, _ := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for password too short", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()


		// Prepare request with password less than 6 characters
		registerReq := schemas.RegisterUserInput{
			Email:    "test@example.com",
			Name:     "Test User",
			Password: "12345",
		}
		reqBody, _ := json.Marshal(registerReq)

		url := fmt.Sprintf("%s/users/register", setup.Server.GetURL())
		resp, _ := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// Verify error message
		var errResp map[string]string
		json.NewDecoder(resp.Body).Decode(&errResp)
		assert.Contains(t, errResp["message"], "at least 6")
	})

	t.Run("should return 400 for password too long", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()


		// Prepare request with password exceeding 50 characters
		longPassword := strings.Repeat("a", 51)
		registerReq := schemas.RegisterUserInput{
			Email:    "test@example.com",
			Name:     "Test User",
			Password: longPassword,
		}
		reqBody, _ := json.Marshal(registerReq)

		url := fmt.Sprintf("%s/users/register", setup.Server.GetURL())
		resp, _ := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// Verify error message
		var errResp map[string]string
		json.NewDecoder(resp.Body).Decode(&errResp)
		assert.Contains(t, errResp["message"], "at most 50")
	})

	t.Run("should return 409 for duplicate email", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()


		// Register first user
		registerReq := schemas.RegisterUserInput{
			Email:    "duplicate@example.com",
			Name:     "First User",
			Password: "password123",
		}
		reqBody, _ := json.Marshal(registerReq)

		url := fmt.Sprintf("%s/users/register", setup.Server.GetURL())
		resp, _ := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
		resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// Try to register second user with same email
		registerReq2 := schemas.RegisterUserInput{
			Email:    "duplicate@example.com",
			Name:     "Second User",
			Password: "password456",
		}
		reqBody2, _ := json.Marshal(registerReq2)

		resp2, _ := http.Post(url, "application/json", bytes.NewBuffer(reqBody2))
		defer resp2.Body.Close()

		assert.Equal(t, http.StatusConflict, resp2.StatusCode)

		// Verify error message
		var errResp map[string]string
		json.NewDecoder(resp2.Body).Decode(&errResp)
		assert.Contains(t, errResp["message"], "already registered")
	})
}

func TestLoginUser(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("should login successfully and set cookies", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()


		// First register a user
		registerReq := schemas.RegisterUserInput{
			Email:    "login@example.com",
			Name:     "Login User",
			Password: "password123",
		}
		reqBody, _ := json.Marshal(registerReq)

		registerURL := fmt.Sprintf("%s/users/register", setup.Server.GetURL())
		registerResp, _ := http.Post(registerURL, "application/json", bytes.NewBuffer(reqBody))
		registerResp.Body.Close()
		require.Equal(t, http.StatusCreated, registerResp.StatusCode)

		// Now login
		loginReq := schemas.LoginUserInput{
			Email:    "login@example.com",
			Password: "password123",
		}
		loginBody, _ := json.Marshal(loginReq)

		loginURL := fmt.Sprintf("%s/users/login", setup.Server.GetURL())
		loginResp, _ := http.Post(loginURL, "application/json", bytes.NewBuffer(loginBody))
		defer loginResp.Body.Close()

		// Assert response status
		assert.Equal(t, http.StatusOK, loginResp.StatusCode)

		// Parse response
		var loginResponse schemas.LoginResponse
		_ = json.NewDecoder(loginResp.Body).Decode(&loginResponse)

		// Assert user data
		assert.Equal(t, "login@example.com", loginResponse.User.Email)
		assert.Equal(t, "Login User", loginResponse.User.Name)

		// Assert cookies are set
		cookies := loginResp.Cookies()
		var accessToken, refreshToken *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "access-token" {
				accessToken = cookie
			}
			if cookie.Name == "refresh-token" {
				refreshToken = cookie
			}
		}

		require.NotNil(t, accessToken, "access-token cookie should be set")
		require.NotNil(t, refreshToken, "refresh-token cookie should be set")

		// Assert cookie properties
		assert.NotEmpty(t, accessToken.Value)
		assert.True(t, accessToken.HttpOnly)
		assert.Equal(t, "/", accessToken.Path)
		assert.Greater(t, accessToken.MaxAge, 0)
		assert.Less(t, accessToken.MaxAge, 16*60) // Less than 16 minutes

		assert.NotEmpty(t, refreshToken.Value)
		assert.True(t, refreshToken.HttpOnly)
		assert.Equal(t, "/", refreshToken.Path)
		assert.Greater(t, refreshToken.MaxAge, 29*24*60*60) // More than 29 days

		// Verify refresh token is stored in database
		var count int
		_ = setup.DB.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM refresh_tokens").Scan(&count)
		assert.Equal(t, 1, count, "Should have one refresh token in database")
	})

	t.Run("should return 400 for invalid JSON", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()


		url := fmt.Sprintf("%s/users/login", setup.Server.GetURL())
		resp, _ := http.Post(url, "application/json", bytes.NewBufferString("{invalid json"))
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for invalid email format", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()


		loginReq := schemas.LoginUserInput{
			Email:    "not-an-email",
			Password: "password123",
		}
		reqBody, _ := json.Marshal(loginReq)

		url := fmt.Sprintf("%s/users/login", setup.Server.GetURL())
		resp, _ := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 401 for non-existent user", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()


		loginReq := schemas.LoginUserInput{
			Email:    "nonexistent@example.com",
			Password: "password123",
		}
		reqBody, _ := json.Marshal(loginReq)

		url := fmt.Sprintf("%s/users/login", setup.Server.GetURL())
		resp, _ := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		var errResp map[string]string
		json.NewDecoder(resp.Body).Decode(&errResp)
		assert.Contains(t, errResp["message"], "Invalid email or password")
	})

	t.Run("should return 401 for wrong password", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()


		// First register a user
		registerReq := schemas.RegisterUserInput{
			Email:    "wrongpass@example.com",
			Name:     "Wrong Pass User",
			Password: "correctpassword",
		}
		reqBody, _ := json.Marshal(registerReq)

		registerURL := fmt.Sprintf("%s/users/register", setup.Server.GetURL())
		registerResp, _ := http.Post(registerURL, "application/json", bytes.NewBuffer(reqBody))
		registerResp.Body.Close()

		// Try to login with wrong password
		loginReq := schemas.LoginUserInput{
			Email:    "wrongpass@example.com",
			Password: "wrongpassword",
		}
		loginBody, _ := json.Marshal(loginReq)

		loginURL := fmt.Sprintf("%s/users/login", setup.Server.GetURL())
		loginResp, _ := http.Post(loginURL, "application/json", bytes.NewBuffer(loginBody))
		defer loginResp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, loginResp.StatusCode)

		var errResp map[string]string
		json.NewDecoder(loginResp.Body).Decode(&errResp)
		assert.Contains(t, errResp["message"], "Invalid email or password")
	})

	t.Run("should create multiple refresh tokens for multiple logins", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()


		// Register a user
		registerReq := schemas.RegisterUserInput{
			Email:    "multisession@example.com",
			Name:     "Multi Session User",
			Password: "password123",
		}
		reqBody, _ := json.Marshal(registerReq)

		registerURL := fmt.Sprintf("%s/users/register", setup.Server.GetURL())
		registerResp, _ := http.Post(registerURL, "application/json", bytes.NewBuffer(reqBody))
		registerResp.Body.Close()

		// Login multiple times
		loginReq := schemas.LoginUserInput{
			Email:    "multisession@example.com",
			Password: "password123",
		}
		loginBody, _ := json.Marshal(loginReq)

		loginURL := fmt.Sprintf("%s/users/login", setup.Server.GetURL())

		// First login
		resp1, _ := http.Post(loginURL, "application/json", bytes.NewBuffer(loginBody))
		resp1.Body.Close()
		assert.Equal(t, http.StatusOK, resp1.StatusCode)

		// Second login
		resp2, _ := http.Post(loginURL, "application/json", bytes.NewBuffer(loginBody))
		resp2.Body.Close()
		assert.Equal(t, http.StatusOK, resp2.StatusCode)

		// Third login
		resp3, _ := http.Post(loginURL, "application/json", bytes.NewBuffer(loginBody))
		resp3.Body.Close()
		assert.Equal(t, http.StatusOK, resp3.StatusCode)

		// Verify 3 refresh tokens in database
		var count int
		_ = setup.DB.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM refresh_tokens").Scan(&count)
		assert.Equal(t, 3, count, "Should have three refresh tokens for three logins")

		// Verify all have different JTIs
		rows, _ := setup.DB.DB.QueryContext(ctx, "SELECT jti FROM refresh_tokens")
		defer rows.Close()

		jtis := make(map[string]bool)
		for rows.Next() {
			var jti string
			_ = rows.Scan(&jti)
			jtis[jti] = true
		}

		assert.Len(t, jtis, 3, "All JTIs should be unique")
	})
}
