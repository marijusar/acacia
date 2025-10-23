package auth

import (
	"acacia/packages/db"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// Resource-based authorization functions that can be used by both HTTP handlers and LLM tools.
// These functions extract user_id from context and verify access to resources via team membership.

// CheckTeamMembership verifies that the user in the context is a member of the specified team
func CheckTeamMembership(ctx context.Context, queries *db.Queries, teamID int64) error {
	userID, ok := ctx.Value(UserIDKey).(int64)
	if !ok {
		return errors.New("unauthorized: user not authenticated")
	}

	isMember, err := queries.CheckUserTeamMembership(ctx, db.CheckUserTeamMembershipParams{
		TeamID: teamID,
		UserID: userID,
	})
	if err != nil {
		return err
	}

	if !isMember {
		return errors.New("user is not a member of this team")
	}

	return nil
}

// CheckProjectAccess verifies that the user has access to a project via team membership
func CheckProjectAccess(ctx context.Context, queries *db.Queries, projectID int64) error {
	userID, ok := ctx.Value(UserIDKey).(int64)
	fmt.Printf("%d has asked for project access", userID)
	if !ok {
		return errors.New("unauthorized: user not authenticated")
	}

	// Get team ID from project
	teamID, err := queries.GetTeamIDByProject(ctx, projectID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("project not found")
		}
		return err
	}

	// Check team membership
	isMember, err := queries.CheckUserTeamMembership(ctx, db.CheckUserTeamMembershipParams{
		TeamID: teamID,
		UserID: userID,
	})
	if err != nil {
		return err
	}

	if !isMember {
		return errors.New("user does not have access to this project")
	}

	return nil
}

// CheckIssueAccess verifies that the user has access to an issue via team membership
func CheckIssueAccess(ctx context.Context, queries *db.Queries, issueID int64) error {
	userID, ok := ctx.Value(UserIDKey).(int64)
	if !ok {
		return errors.New("unauthorized: user not authenticated")
	}

	// Get team ID from issue
	teamID, err := queries.GetTeamIDByIssue(ctx, issueID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("issue not found")
		}
		return err
	}

	// Check team membership
	isMember, err := queries.CheckUserTeamMembership(ctx, db.CheckUserTeamMembershipParams{
		TeamID: teamID,
		UserID: userID,
	})
	if err != nil {
		return err
	}

	if !isMember {
		return errors.New("user does not have access to this issue")
	}

	return nil
}

// CheckColumnAccess verifies that the user has access to a project status column via team membership
func CheckColumnAccess(ctx context.Context, queries *db.Queries, columnID int64) error {
	userID, ok := ctx.Value(UserIDKey).(int64)
	if !ok {
		return errors.New("unauthorized: user not authenticated")
	}

	// Get team ID from column
	teamID, err := queries.GetTeamIDByProjectStatusColumn(ctx, columnID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("column not found")
		}
		return err
	}

	// Check team membership
	isMember, err := queries.CheckUserTeamMembership(ctx, db.CheckUserTeamMembershipParams{
		TeamID: teamID,
		UserID: userID,
	})
	if err != nil {
		return err
	}

	if !isMember {
		return errors.New("user does not have access to this column")
	}

	return nil
}
