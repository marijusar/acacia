package tools

import (
	"acacia/packages/auth"
	"acacia/packages/db"
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

// SearchIssuesTool searches for issues across all projects the user has access to
type SearchIssuesTool struct {
	queries *db.Queries
	logger  *logrus.Logger
}

// NewSearchIssuesTool creates a new SearchIssuesTool
func NewSearchIssuesTool(queries *db.Queries, logger *logrus.Logger) *SearchIssuesTool {
	return &SearchIssuesTool{
		queries: queries,
		logger:  logger,
	}
}

func (t *SearchIssuesTool) Name() string {
	return "search_issues"
}

func (t *SearchIssuesTool) Description() string {
	return "Search for issues across all projects the user has access to. Returns issues matching the search query."
}

func (t *SearchIssuesTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"query": map[string]interface{}{
				"type":        "string",
				"description": "The search query to find matching issues",
			},
		},
		"required": []string{"query"},
	}
}

func (t *SearchIssuesTool) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	t.logger.WithField("args", args).Info("[SEARCH_ISSUES] Tool called")

	// Extract user ID from context
	userID, ok := ctx.Value(auth.UserIDKey).(int64)
	if !ok {
		t.logger.Error("[SEARCH_ISSUES] User not authenticated")
		return nil, fmt.Errorf("unauthorized: user not authenticated")
	}

	// Extract search query from arguments
	query, ok := args["query"].(string)
	if !ok {
		t.logger.Error("[SEARCH_ISSUES] Invalid query argument")
		return nil, fmt.Errorf("invalid query: expected string")
	}

	t.logger.WithFields(logrus.Fields{
		"user_id": userID,
		"query":   query,
	}).Info("[SEARCH_ISSUES] Fetching user projects")

	queryLower := strings.ToLower(query)

	// Get all projects the user has access to
	projects, err := t.queries.GetProjects(ctx, userID)
	if err != nil {
		t.logger.WithError(err).Error("[SEARCH_ISSUES] Failed to fetch projects")
		return nil, err
	}

	t.logger.WithField("project_count", len(projects)).Info("[SEARCH_ISSUES] Searching issues across projects")

	// Search issues in each project
	var matchingIssues []db.Issue
	for _, project := range projects {
		issues, err := t.queries.GetProjectIssues(ctx, int32(project.ID))
		if err != nil {
			t.logger.WithError(err).WithField("project_id", project.ID).Warn("[SEARCH_ISSUES] Failed to fetch issues for project, skipping")
			// Skip projects with errors, continue searching
			continue
		}

		// Filter issues by query (simple text matching)
		for _, issue := range issues {
			titleMatch := strings.Contains(strings.ToLower(issue.Name), queryLower)
			descMatch := false
			if issue.Description.Valid {
				descMatch = strings.Contains(strings.ToLower(issue.Description.String), queryLower)
			}

			if titleMatch || descMatch {
				matchingIssues = append(matchingIssues, issue)
			}
		}
	}

	t.logger.WithFields(logrus.Fields{
		"query":        query,
		"match_count":  len(matchingIssues),
		"searched_projects": len(projects),
	}).Info("[SEARCH_ISSUES] Search completed successfully")

	return matchingIssues, nil
}
