package auth

import (
	"acacia/packages/db"
	"context"
)

// ResourceResolver defines the interface for resolving team IDs from resources
type ResourceResolver interface {
	GetTeamID(ctx context.Context, resourceID int64) (int64, error)
}

// ProjectResolver resolves team IDs for projects
type ProjectResolver struct {
	queries *db.Queries
}

func NewProjectResolver(queries *db.Queries) *ProjectResolver {
	return &ProjectResolver{queries: queries}
}

func (r *ProjectResolver) GetTeamID(ctx context.Context, resourceID int64) (int64, error) {
	return r.queries.GetTeamIDByProject(ctx, resourceID)
}

// ProjectStatusColumnResolver resolves team IDs for project status columns
type ProjectStatusColumnResolver struct {
	queries *db.Queries
}

func NewProjectStatusColumnResolver(queries *db.Queries) *ProjectStatusColumnResolver {
	return &ProjectStatusColumnResolver{queries: queries}
}

func (r *ProjectStatusColumnResolver) GetTeamID(ctx context.Context, resourceID int64) (int64, error) {
	return r.queries.GetTeamIDByProjectStatusColumn(ctx, resourceID)
}

// IssueResolver resolves team IDs for issues
type IssueResolver struct {
	queries *db.Queries
}

func NewIssueResolver(queries *db.Queries) *IssueResolver {
	return &IssueResolver{queries: queries}
}

func (r *IssueResolver) GetTeamID(ctx context.Context, resourceID int64) (int64, error) {
	return r.queries.GetTeamIDByIssue(ctx, resourceID)
}

// TeamResolver is a simple resolver that returns the team ID itself
// Used when checking if a user is a member of a specific team
type TeamResolver struct{}

func NewTeamResolver() *TeamResolver {
	return &TeamResolver{}
}

func (r *TeamResolver) GetTeamID(ctx context.Context, resourceID int64) (int64, error) {
	// For teams, the resource ID is the team ID itself
	return resourceID, nil
}
