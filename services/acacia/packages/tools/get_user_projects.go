package tools

import (
	"acacia/packages/auth"
	"acacia/packages/db"
	"context"
	"errors"

	"github.com/sirupsen/logrus"
)

// GetUserProjectsTool returns all projects the user has access to
type GetUserProjectsTool struct {
	queries *db.Queries
	logger  *logrus.Logger
}

// NewGetUserProjectsTool creates a new GetUserProjectsTool
func NewGetUserProjectsTool(queries *db.Queries, logger *logrus.Logger) *GetUserProjectsTool {
	return &GetUserProjectsTool{
		queries: queries,
		logger:  logger,
	}
}

func (t *GetUserProjectsTool) Name() string {
	return "get_user_projects"
}

func (t *GetUserProjectsTool) Description() string {
	return "Get all projects that the user has access to. Returns a list of projects with their basic information."
}

func (t *GetUserProjectsTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type":       "object",
		"properties": map[string]interface{}{},
		"required":   []string{},
	}
}

func (t *GetUserProjectsTool) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	t.logger.Info("[GET_USER_PROJECTS] Tool called")

	// Extract user ID from context (set by auth middleware)
	userID, ok := ctx.Value(auth.UserIDKey).(int64)
	if !ok {
		t.logger.Error("[GET_USER_PROJECTS] User not authenticated")
		return nil, errors.New("unauthorized: user not authenticated")
	}

	t.logger.WithField("user_id", userID).Info("[GET_USER_PROJECTS] Fetching projects for user")

	// Get user's projects (via team membership)
	projects, err := t.queries.GetProjects(ctx, userID)
	if err != nil {
		t.logger.WithError(err).Error("[GET_USER_PROJECTS] Failed to fetch projects")
		return nil, err
	}

	t.logger.WithField("project_count", len(projects)).Info("[GET_USER_PROJECTS] Successfully fetched projects")
	return projects, nil
}
