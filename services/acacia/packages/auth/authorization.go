package auth

import (
	"acacia/packages/db"
	"acacia/packages/httperr"
	"errors"
	"net/http"

	"github.com/sirupsen/logrus"
)

var (
	ErrForbidden = errors.New("Forbidden: insufficient permissions")
)

// AccessChecker is a callback function that checks if a user has access
// It receives the request and queries object
// The checker is responsible for extracting user ID, resource IDs, and verifying access
// Returns an error if access is denied, nil if access is granted
type AccessChecker func(r *http.Request, queries *db.Queries) error

// AuthorizationMiddleware handles resource authorization using callback-based access checkers
type AuthorizationMiddleware struct {
	queries *db.Queries
	logger  *logrus.Logger
}

// NewAuthorizationMiddleware creates a new authorization middleware
func NewAuthorizationMiddleware(queries *db.Queries, logger *logrus.Logger) *AuthorizationMiddleware {
	return &AuthorizationMiddleware{
		queries: queries,
		logger:  logger,
	}
}

// RequireAccess creates a middleware that checks if the user has access using the provided checker
// checker: callback function that verifies access
func (m *AuthorizationMiddleware) RequireAccess(checker AccessChecker) func(http.Handler) http.Handler {
	return httperr.WithMiddlewareErrorHandler(func(w http.ResponseWriter, r *http.Request, next http.Handler) error {
		// Get user ID from context (set by authentication middleware)
		// This is just for early termination if not authenticated
		_, ok := GetUserID(r)
		if !ok {
			m.logger.Debug("Missing user ID in context")
			return httperr.WithStatus(errors.New("Unauthorized"), http.StatusUnauthorized)
		}

		// Check authorization using the provided checker
		if err := checker(r, m.queries); err != nil {
			m.logger.WithError(err).Debug("Access denied")
			return httperr.WithStatus(ErrForbidden, http.StatusForbidden)
		}

		// Authorization successful, proceed to next handler
		next.ServeHTTP(w, r)
		return nil
	})
}
