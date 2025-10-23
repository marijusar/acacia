package tools

import (
	"acacia/packages/auth"
	"acacia/packages/db"
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

// GetIssueDetailsTool returns detailed information about a specific issue
type GetIssueDetailsTool struct {
	queries *db.Queries
	logger  *logrus.Logger
}

// NewGetIssueDetailsTool creates a new GetIssueDetailsTool
func NewGetIssueDetailsTool(queries *db.Queries, logger *logrus.Logger) *GetIssueDetailsTool {
	return &GetIssueDetailsTool{
		queries: queries,
		logger:  logger,
	}
}

func (t *GetIssueDetailsTool) Name() string {
	return "get_issue_details"
}

func (t *GetIssueDetailsTool) Description() string {
	return "Get detailed information about a specific issue. Requires the issue ID."
}

func (t *GetIssueDetailsTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"issue_id": map[string]interface{}{
				"type":        "number",
				"description": "The ID of the issue to retrieve",
			},
		},
		"required": []string{"issue_id"},
	}
}

func (t *GetIssueDetailsTool) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	t.logger.WithField("args", args).Info("[GET_ISSUE_DETAILS] Tool called")

	// Extract issue ID from arguments
	issueIDFloat, ok := args["issue_id"].(float64)
	if !ok {
		t.logger.Error("[GET_ISSUE_DETAILS] Invalid issue_id argument")
		return nil, fmt.Errorf("invalid issue_id: expected number")
	}
	issueID := int64(issueIDFloat)

	t.logger.WithField("issue_id", issueID).Info("[GET_ISSUE_DETAILS] Checking issue access")

	// Check authorization using shared resource checker
	// This uses the SAME authorization logic as GET /issues/{id}
	if err := auth.CheckIssueAccess(ctx, t.queries, issueID); err != nil {
		t.logger.WithError(err).WithField("issue_id", issueID).Error("[GET_ISSUE_DETAILS] Authorization failed")
		return nil, err
	}

	t.logger.WithField("issue_id", issueID).Info("[GET_ISSUE_DETAILS] Authorization passed, fetching issue")

	// User is authorized - fetch issue details
	issue, err := t.queries.GetIssueByID(ctx, issueID)
	if err != nil {
		t.logger.WithError(err).WithField("issue_id", issueID).Error("[GET_ISSUE_DETAILS] Failed to fetch issue")
		return nil, err
	}

	t.logger.WithField("issue_id", issueID).Info("[GET_ISSUE_DETAILS] Successfully fetched issue details")
	return issue, nil
}
