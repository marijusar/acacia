package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken     = errors.New("Invalid token")
	ErrExpiredToken     = errors.New("Token has expired")
	ErrInvalidSignature = errors.New("Invalid token signature")
)

// TokenClaims represents the JWT claims
type TokenClaims struct {
	UserID int64  `json:"user_id"`
	JTI    string `json:"jti,omitempty"` // Only for refresh tokens
	jwt.RegisteredClaims
}

// JWTManager handles JWT token generation and validation
type JWTManager struct {
	secretKey            []byte
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(secretKey string, accessTokenDuration, refreshTokenDuration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:            []byte(secretKey),
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}
}

// GenerateJTI generates a unique JWT ID
func GenerateJTI() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateAccessToken generates an access token (15 minutes)
func (m *JWTManager) GenerateAccessToken(userID int64) (string, time.Time, error) {
	expiresAt := time.Now().Add(m.accessTokenDuration)

	claims := &TokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(m.secretKey)
	return tokenString, expiresAt, err
}

// GenerateRefreshToken generates a refresh token with JTI (30 days)
func (m *JWTManager) GenerateRefreshToken(userID int64, jti string) (string, time.Time, error) {
	expiresAt := time.Now().Add(m.refreshTokenDuration)

	claims := &TokenClaims{
		UserID: userID,
		JTI:    jti,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(m.secretKey)
	return tokenString, expiresAt, err
}

// ValidateToken validates a token and returns its claims
func (m *JWTManager) ValidateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (any, error) {
		// Verify signing method is exactly HS256
		if token.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidSignature
		}
		return m.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
