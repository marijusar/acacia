package tools

import (
	"acacia/packages/auth"
	"acacia/packages/db"
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

// GetProjectDetailsTool returns detailed information about a specific project
type GetProjectDetailsTool struct {
	queries *db.Queries
	logger  *logrus.Logger
}

// NewGetProjectDetailsTool creates a new GetProjectDetailsTool
func NewGetProjectDetailsTool(queries *db.Queries, logger *logrus.Logger) *GetProjectDetailsTool {
	return &GetProjectDetailsTool{
		queries: queries,
		logger:  logger,
	}
}

func (t *GetProjectDetailsTool) Name() string {
	return "get_project_details"
}

func (t *GetProjectDetailsTool) Description() string {
	return "Get detailed information about a specific project, including its columns and issues. Requires the project ID."
}

func (t *GetProjectDetailsTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"project_id": map[string]interface{}{
				"type":        "number",
				"description": "The ID of the project to retrieve",
			},
		},
		"required": []string{"project_id"},
	}
}

func (t *GetProjectDetailsTool) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	t.logger.WithField("args", args).Info("[GET_PROJECT_DETAILS] Tool called")

	// Extract project ID from arguments
	projectIDFloat, ok := args["project_id"].(float64)
	if !ok {
		t.logger.Error("[GET_PROJECT_DETAILS] Invalid project_id argument")
		return nil, fmt.Errorf("invalid project_id: expected number")
	}
	projectID := int64(projectIDFloat)

	t.logger.WithField("project_id", projectID).Info("[GET_PROJECT_DETAILS] Checking project access")

	// Check authorization using shared resource checker
	// This uses the SAME authorization logic as GET /projects/{id}/details
	if err := auth.CheckProjectAccess(ctx, t.queries, projectID); err != nil {
		t.logger.WithError(err).WithField("project_id", projectID).Error("[GET_PROJECT_DETAILS] Authorization failed")
		return nil, err
	}

	t.logger.WithField("project_id", projectID).Info("[GET_PROJECT_DETAILS] Authorization passed, fetching project details")

	// User is authorized - fetch project details
	project, err := t.queries.GetProjectByID(ctx, projectID)
	if err != nil {
		t.logger.WithError(err).WithField("project_id", projectID).Error("[GET_PROJECT_DETAILS] Failed to fetch project")
		return nil, err
	}

	columns, err := t.queries.GetProjectStatusColumnsByProjectID(ctx, int32(projectID))
	if err != nil {
		t.logger.WithError(err).WithField("project_id", projectID).Error("[GET_PROJECT_DETAILS] Failed to fetch columns")
		return nil, err
	}

	issues, err := t.queries.GetProjectIssues(ctx, int32(projectID))
	if err != nil {
		t.logger.WithError(err).WithField("project_id", projectID).Error("[GET_PROJECT_DETAILS] Failed to fetch issues")
		return nil, err
	}

	t.logger.WithFields(logrus.Fields{
		"project_id":   projectID,
		"column_count": len(columns),
		"issue_count":  len(issues),
	}).Info("[GET_PROJECT_DETAILS] Successfully fetched project details")

	// Return structured response
	return map[string]interface{}{
		"project": project,
		"columns": columns,
		"issues":  issues,
	}, nil
}
