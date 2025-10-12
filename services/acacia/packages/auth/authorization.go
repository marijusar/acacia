package auth

import (
	"acacia/packages/db"
	"acacia/packages/httperr"
	"bytes"
	"context"
	"database/sql"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

var (
	ErrUnauthorized          = errors.New("Unauthorized: user does not have access to this resource")
	ErrResourceNotFound      = errors.New("Resource not found")
	ErrInvalidResourceID     = errors.New("Invalid resource ID")
	ErrMissingAuthentication = errors.New("Authentication required")
)

// ResourceType represents different types of resources that can be authorized
type ResourceType string

const (
	ResourceTypeProject             ResourceType = "project"
	ResourceTypeProjectStatusColumn ResourceType = "project_status_column"
	ResourceTypeIssue               ResourceType = "issue"
	ResourceTypeTeam                ResourceType = "team"
)

// AuthorizationMiddleware handles team-based authorization
type AuthorizationMiddleware struct {
	queries   *db.Queries
	logger    *logrus.Logger
	resolvers map[ResourceType]ResourceResolver
}

// NewAuthorizationMiddleware creates a new authorization middleware
func NewAuthorizationMiddleware(queries *db.Queries, logger *logrus.Logger) *AuthorizationMiddleware {
	m := &AuthorizationMiddleware{
		queries:   queries,
		logger:    logger,
		resolvers: make(map[ResourceType]ResourceResolver),
	}

	// Register default resolvers
	m.resolvers[ResourceTypeProject] = NewProjectResolver(queries)
	m.resolvers[ResourceTypeProjectStatusColumn] = NewProjectStatusColumnResolver(queries)
	m.resolvers[ResourceTypeIssue] = NewIssueResolver(queries)
	m.resolvers[ResourceTypeTeam] = NewTeamResolver()

	return m
}

// RegisterResolver allows registering custom resolvers for resource types
func (m *AuthorizationMiddleware) RegisterResolver(resourceType ResourceType, resolver ResourceResolver) {
	m.resolvers[resourceType] = resolver
}

// wrapAuthorizationError wraps authorization errors with appropriate HTTP status codes
func (m *AuthorizationMiddleware) wrapAuthorizationError(err error, userID int64, resourceType ResourceType, resourceID int64) error {
	if errors.Is(err, ErrResourceNotFound) {
		return httperr.WithStatus(ErrResourceNotFound, http.StatusNotFound)
	}
	if errors.Is(err, ErrUnauthorized) {
		m.logger.WithFields(logrus.Fields{
			"user_id":       userID,
			"resource_type": resourceType,
			"resource_id":   resourceID,
		}).Debug("Unauthorized access attempt")
		return httperr.WithStatus(errors.New("Forbidden: insufficient permissions"), http.StatusForbidden)
	}
	m.logger.WithError(err).Error("Authorization check failed")
	return httperr.WithStatus(errors.New("Internal server error"), http.StatusInternalServerError)
}

// RequireResourceAccess creates a middleware that checks if the user has access to a resource
// resourceType: the type of resource being accessed
// resourceIDParam: the name of the URL parameter containing the resource ID (e.g., "id", "project_id")
func (m *AuthorizationMiddleware) RequireResourceAccess(resourceType ResourceType, resourceIDParam string) func(http.Handler) http.Handler {
	return httperr.WithMiddlewareErrorHandler(func(w http.ResponseWriter, r *http.Request, next http.Handler) error {
		// Get user ID from context (set by authentication middleware)
		userID, ok := GetUserID(r)
		if !ok {
			m.logger.Debug("Missing user ID in context")
			return httperr.WithStatus(ErrMissingAuthentication, http.StatusUnauthorized)
		}

		// Extract resource ID from URL parameter
		resourceIDStr := chi.URLParam(r, resourceIDParam)
		resourceID, err := strconv.ParseInt(resourceIDStr, 10, 64)
		if err != nil {
			m.logger.WithError(err).WithField("param", resourceIDParam).Debug("Invalid resource ID")
			return httperr.WithStatus(ErrInvalidResourceID, http.StatusBadRequest)
		}

		// Check authorization
		if err := m.checkAccess(r.Context(), userID, resourceType, resourceID); err != nil {
			return m.wrapAuthorizationError(err, userID, resourceType, resourceID)
		}

		// Authorization successful, proceed to next handler
		next.ServeHTTP(w, r)
		return nil
	})
}

// checkAccess verifies if a user has access to a specific resource
func (m *AuthorizationMiddleware) checkAccess(ctx context.Context, userID int64, resourceType ResourceType, resourceID int64) error {
	// Get the appropriate resolver
	resolver, ok := m.resolvers[resourceType]
	if !ok {
		return errors.New("unknown resource type")
	}

	// Resolve the team ID for the resource
	teamID, err := resolver.GetTeamID(ctx, resourceID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrResourceNotFound
		}
		return err
	}

	// Check if user is a member of the team
	isMember, err := m.queries.CheckUserTeamMembership(ctx, db.CheckUserTeamMembershipParams{
		TeamID: teamID,
		UserID: userID,
	})

	if err != nil {
		return err
	}

	if !isMember {
		return ErrUnauthorized
	}

	return nil
}

// CheckResourceAccess is a helper function to check access programmatically (not as middleware)
// Useful for controllers that need to check access before performing operations
func (m *AuthorizationMiddleware) CheckResourceAccess(ctx context.Context, userID int64, resourceType ResourceType, resourceID int64) error {
	return m.checkAccess(ctx, userID, resourceType, resourceID)
}

// BodyResourceExtractor is a function type that extracts a resource ID from request body
// It receives the request body bytes and returns the resource ID to check authorization against
type BodyResourceExtractor func(body []byte) (int64, error)

// RequireResourceAccessFromBody creates a middleware that checks resource access based on request body
// resourceType: the type of resource being accessed (used to determine which resolver to use)
// extractor: function that extracts the resource ID from the request body
//
// This is useful for POST/PUT requests where the resource identifier comes from the body
// For example: POST /issues with {"column_id": 5} - we need to check if user has access to column 5
func (m *AuthorizationMiddleware) RequireResourceAccessFromBody(resourceType ResourceType, extractor BodyResourceExtractor) func(http.Handler) http.Handler {
	return httperr.WithMiddlewareErrorHandler(func(w http.ResponseWriter, r *http.Request, next http.Handler) error {
		// Get user ID from context (set by authentication middleware)
		userID, ok := GetUserID(r)
		if !ok {
			m.logger.Debug("Missing user ID in context")
			return httperr.WithStatus(ErrMissingAuthentication, http.StatusUnauthorized)
		}

		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			m.logger.WithError(err).Error("Failed to read request body")
			return httperr.WithStatus(errors.New("Failed to read request body"), http.StatusBadRequest)
		}

		// Restore the body for the next handler
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		// Extract resource ID from body
		resourceID, err := extractor(body)
		if err != nil {
			m.logger.WithError(err).Debug("Failed to extract resource ID from body")
			return httperr.WithStatus(errors.New("Invalid request body"), http.StatusBadRequest)
		}

		// Check authorization
		if err := m.checkAccess(r.Context(), userID, resourceType, resourceID); err != nil {
			return m.wrapAuthorizationError(err, userID, resourceType, resourceID)
		}

		// Authorization successful, proceed to next handler
		next.ServeHTTP(w, r)
		return nil
	})
}
