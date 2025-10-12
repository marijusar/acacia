package auth

import (
	"acacia/packages/db"
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockQueries is a mock implementation of db.Queries for testing
type MockQueries struct {
	mock.Mock
}

func (m *MockQueries) GetRefreshTokenByUserAndJTI(ctx context.Context, arg db.GetRefreshTokenByUserAndJTIParams) (db.RefreshToken, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(db.RefreshToken), args.Error(1)
}

// testSetup contains all dependencies needed for testing the middleware
type testSetup struct {
	jwtManager *JWTManager
	queries    *MockQueries
	middleware *AuthMiddleware
	logger     *logrus.Logger
}

// withTestSetup creates a new test setup with default configuration
func withTestSetup() *testSetup {
	jwtManager := NewJWTManager("test-secret-key", 15*time.Minute, 30*24*time.Hour)
	mockQueries := new(MockQueries)
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel) // Silence logs during tests
	middleware := NewAuthMiddleware(jwtManager, mockQueries, logger)

	return &testSetup{
		jwtManager: jwtManager,
		queries:    mockQueries,
		middleware: middleware,
		logger:     logger,
	}
}

func TestHandleMiddleware(t *testing.T) {
	t.Run("should authenticate with valid access token", func(t *testing.T) {
		setup := withTestSetup()

		userID := int64(123)
		accessToken, _, _ := setup.jwtManager.GenerateAccessToken(userID)

		// Create request with valid access token
		req := httptest.NewRequest("GET", "/protected", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access-token",
			Value: accessToken,
		})

		// Create response recorder
		w := httptest.NewRecorder()

		// Create a test handler that checks if user ID is in context
		var capturedUserID int64
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			uid, ok := GetUserID(r)
			assert.True(t, ok)
			capturedUserID = uid
			w.WriteHeader(http.StatusOK)
		})

		// Execute middleware
		setup.middleware.Handle(nextHandler).ServeHTTP(w, req)

		// Assert
		assert.Equal(t, userID, capturedUserID)
		assert.Equal(t, http.StatusOK, w.Code)

		// Verify no database calls were made (access token was valid)
		setup.queries.AssertNotCalled(t, "GetRefreshTokenByUserAndJTI")
	})

	t.Run("should refresh access token when expired but refresh token is valid", func(t *testing.T) {
		setup := withTestSetup()

		userID := int64(123)
		jti := "test-jti"

		// Generate access token with very short expiration and wait for it to expire
		shortLivedManager := NewJWTManager("test-secret-key", 1*time.Millisecond, 30*24*time.Hour)
		accessToken, _, _ := shortLivedManager.GenerateAccessToken(userID)
		time.Sleep(2 * time.Millisecond)

		// Generate valid refresh token
		refreshToken, expiresAt, _ := setup.jwtManager.GenerateRefreshToken(userID, jti)

		// Mock database call to return valid refresh token
		setup.queries.On("GetRefreshTokenByUserAndJTI", mock.Anything, db.GetRefreshTokenByUserAndJTIParams{
			UserID: userID,
			Jti:    jti,
		}).Return(db.RefreshToken{
			ID:        1,
			UserID:    userID,
			Jti:       jti,
			ExpiresAt: expiresAt,
			CreatedAt: time.Now(),
		}, nil)

		// Create request with expired access token and valid refresh token
		req := httptest.NewRequest("GET", "/protected", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access-token",
			Value: accessToken,
		})
		req.AddCookie(&http.Cookie{
			Name:  "refresh-token",
			Value: refreshToken,
		})

		// Create response recorder
		w := httptest.NewRecorder()

		// Create a test handler
		var capturedUserID int64
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			uid, ok := GetUserID(r)
			assert.True(t, ok)
			capturedUserID = uid
			w.WriteHeader(http.StatusOK)
		})

		// Execute middleware
		setup.middleware.Handle(nextHandler).ServeHTTP(w, req)

		// Assert
		assert.Equal(t, userID, capturedUserID)
		assert.Equal(t, http.StatusOK, w.Code)

		// Check that new access token cookie was set
		cookies := w.Result().Cookies()
		foundNewAccessToken := false
		for _, cookie := range cookies {
			if cookie.Name == "access-token" {
				foundNewAccessToken = true
				assert.NotEmpty(t, cookie.Value)
				assert.NotEqual(t, accessToken, cookie.Value, "New access token should be different from expired one")
			}
		}
		assert.True(t, foundNewAccessToken, "New access token should be set")

		// Verify database was called to validate refresh token
		setup.queries.AssertExpectations(t)
	})

	t.Run("should return 401 when no tokens provided", func(t *testing.T) {
		setup := withTestSetup()

		// Create request without any tokens
		req := httptest.NewRequest("GET", "/protected", nil)
		w := httptest.NewRecorder()

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("Should not reach next handler")
		})

		// Execute middleware
		setup.middleware.Handle(nextHandler).ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Unauthorized")
	})

	t.Run("should return 401 when only access token is invalid and no refresh token", func(t *testing.T) {
		setup := withTestSetup()

		// Create request with invalid access token
		req := httptest.NewRequest("GET", "/protected", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access-token",
			Value: "invalid-token-string",
		})
		w := httptest.NewRecorder()

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("Should not reach next handler")
		})

		// Execute middleware
		setup.middleware.Handle(nextHandler).ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Unauthorized")
	})

	t.Run("should return 401 when refresh token is revoked", func(t *testing.T) {
		setup := withTestSetup()

		userID := int64(123)
		jti := "test-jti"

		// Generate expired access token
		shortLivedManager := NewJWTManager("test-secret-key", 1*time.Millisecond, 30*24*time.Hour)
		accessToken, _, _ := shortLivedManager.GenerateAccessToken(userID)
		time.Sleep(2 * time.Millisecond)

		// Generate refresh token
		refreshToken, _, _ := setup.jwtManager.GenerateRefreshToken(userID, jti)

		// Mock database call to return "not found" (revoked token)
		setup.queries.On("GetRefreshTokenByUserAndJTI", mock.Anything, db.GetRefreshTokenByUserAndJTIParams{
			UserID: userID,
			Jti:    jti,
		}).Return(db.RefreshToken{}, sql.ErrNoRows)

		// Create request
		req := httptest.NewRequest("GET", "/protected", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access-token",
			Value: accessToken,
		})
		req.AddCookie(&http.Cookie{
			Name:  "refresh-token",
			Value: refreshToken,
		})

		w := httptest.NewRecorder()

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("Should not reach next handler")
		})

		// Execute middleware
		setup.middleware.Handle(nextHandler).ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Unauthorized")

		setup.queries.AssertExpectations(t)
	})

	t.Run("should return 401 when refresh token is expired", func(t *testing.T) {
		setup := withTestSetup()

		userID := int64(123)
		jti := "test-jti"

		// Generate tokens with very short durations
		shortLivedManager := NewJWTManager("test-secret-key", 1*time.Millisecond, 2*time.Millisecond)
		accessToken, _, _ := shortLivedManager.GenerateAccessToken(userID)
		refreshToken, _, _ := shortLivedManager.GenerateRefreshToken(userID, jti)

		// Wait for both tokens to expire
		time.Sleep(3 * time.Millisecond)

		// Create request with expired tokens
		req := httptest.NewRequest("GET", "/protected", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access-token",
			Value: accessToken,
		})
		req.AddCookie(&http.Cookie{
			Name:  "refresh-token",
			Value: refreshToken,
		})

		w := httptest.NewRecorder()

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("Should not reach next handler")
		})

		// Execute middleware
		setup.middleware.Handle(nextHandler).ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Unauthorized")

		// Database should not be called if refresh token JWT validation fails
		setup.queries.AssertNotCalled(t, "GetRefreshTokenByUserAndJTI")
	})

	t.Run("should return 401 when refresh token has wrong signature", func(t *testing.T) {
		setup := withTestSetup()

		userID := int64(123)
		jti := "test-jti"

		// Generate expired access token with correct secret
		shortLivedManager := NewJWTManager("test-secret-key", 1*time.Millisecond, 30*24*time.Hour)
		accessToken, _, _ := shortLivedManager.GenerateAccessToken(userID)
		time.Sleep(2 * time.Millisecond)

		// Generate refresh token with WRONG secret
		wrongJwtManager := NewJWTManager("wrong-secret-key", 15*time.Minute, 30*24*time.Hour)
		refreshToken, _, _ := wrongJwtManager.GenerateRefreshToken(userID, jti)

		// Create request
		req := httptest.NewRequest("GET", "/protected", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access-token",
			Value: accessToken,
		})
		req.AddCookie(&http.Cookie{
			Name:  "refresh-token",
			Value: refreshToken,
		})

		w := httptest.NewRecorder()

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("Should not reach next handler")
		})

		// Execute middleware
		setup.middleware.Handle(nextHandler).ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Unauthorized")

		// Database should not be called if refresh token JWT validation fails
		setup.queries.AssertNotCalled(t, "GetRefreshTokenByUserAndJTI")
	})
}

func TestGetUserID(t *testing.T) {
	t.Run("should retrieve user ID from context", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		ctx := context.WithValue(req.Context(), UserIDKey, int64(456))
		req = req.WithContext(ctx)

		userID, ok := GetUserID(req)

		assert.True(t, ok)
		assert.Equal(t, int64(456), userID)
	})

	t.Run("should return false when user ID not in context", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)

		_, ok := GetUserID(req)

		assert.False(t, ok)
	})
}
