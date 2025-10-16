package auth

import (
	"acacia/packages/db"
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// Common access checker functions

// CheckTeamMembershipByURLParam checks if user is a member of a team specified in URL parameter
func CheckTeamMembershipByURLParam(paramName string) AccessChecker {
	return func(r *http.Request, queries *db.Queries) error {
		userID, _ := GetUserID(r)

		teamIDStr := chi.URLParam(r, paramName)
		teamID, err := strconv.ParseInt(teamIDStr, 10, 64)
		if err != nil {
			// Invalid ID - let the handler return 400
			return nil
		}

		isMember, err := queries.CheckUserTeamMembership(r.Context(), db.CheckUserTeamMembershipParams{
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
}

// CheckTeamMembershipByBody checks if user is a member of a team specified in request body
func CheckTeamMembershipByBody(fieldName string) AccessChecker {
	return func(r *http.Request, queries *db.Queries) error {
		userID, _ := GetUserID(r)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			return errors.New("failed to read request body")
		}
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		var data map[string]interface{}
		if err := json.Unmarshal(body, &data); err != nil {
			// Invalid JSON - let the handler return 400
			return nil
		}

		teamIDFloat, ok := data[fieldName].(float64)
		if !ok {
			// Missing or invalid team_id - let the handler validate and return appropriate error
			return nil
		}
		teamID := int64(teamIDFloat)

		isMember, err := queries.CheckUserTeamMembership(r.Context(), db.CheckUserTeamMembershipParams{
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
}

// CheckProjectAccessByURLParam checks if user has access to a project via team membership
func CheckProjectAccessByURLParam(paramName string) AccessChecker {
	return func(r *http.Request, queries *db.Queries) error {
		userID, _ := GetUserID(r)

		projectIDStr := chi.URLParam(r, paramName)
		projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
		if err != nil {
			// Invalid ID - let the handler return 400
			return nil
		}

		// Get team ID from project
		teamID, err := queries.GetTeamIDByProject(r.Context(), projectID)
		if err != nil {
			if err == sql.ErrNoRows {
				// Project doesn't exist - let the handler return 404
				return nil
			}
			return err
		}

		// Check team membership
		isMember, err := queries.CheckUserTeamMembership(r.Context(), db.CheckUserTeamMembershipParams{
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
}

// CheckIssueAccessByURLParam checks if user has access to an issue via team membership
func CheckIssueAccessByURLParam(paramName string) AccessChecker {
	return func(r *http.Request, queries *db.Queries) error {
		userID, _ := GetUserID(r)

		issueIDStr := chi.URLParam(r, paramName)
		issueID, err := strconv.ParseInt(issueIDStr, 10, 64)
		if err != nil {
			// Invalid ID - let the handler return 400
			return nil
		}

		// Get team ID from issue
		teamID, err := queries.GetTeamIDByIssue(r.Context(), issueID)
		if err != nil {
			if err == sql.ErrNoRows {
				// Issue doesn't exist - let the handler return 404
				return nil
			}
			return err
		}

		// Check team membership
		isMember, err := queries.CheckUserTeamMembership(r.Context(), db.CheckUserTeamMembershipParams{
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
}

// CheckIssueAccessByBody checks if user has access to an issue from request body
func CheckIssueAccessByBody() AccessChecker {
	return func(r *http.Request, queries *db.Queries) error {
		userID, _ := GetUserID(r)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			return errors.New("failed to read request body")
		}
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		var data struct {
			ID int64 `json:"id"`
		}
		if err := json.Unmarshal(body, &data); err != nil {
			// Invalid JSON - let the handler return 400
			return nil
		}

		if data.ID == 0 {
			// Missing or invalid issue ID - let the handler validate
			return nil
		}

		// Get team ID from issue
		teamID, err := queries.GetTeamIDByIssue(r.Context(), data.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				// Issue doesn't exist - let the handler return 404
				return nil
			}
			return err
		}

		// Check team membership
		isMember, err := queries.CheckUserTeamMembership(r.Context(), db.CheckUserTeamMembershipParams{
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
}

// CheckColumnAccessByBody checks if user has access to a project status column from request body
func CheckColumnAccessByBody() AccessChecker {
	return func(r *http.Request, queries *db.Queries) error {
		userID, _ := GetUserID(r)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			return errors.New("failed to read request body")
		}
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		var data struct {
			ColumnID int64 `json:"column_id"`
		}
		if err := json.Unmarshal(body, &data); err != nil {
			// Invalid JSON - let the handler return 400
			return nil
		}

		if data.ColumnID == 0 {
			// Missing or invalid column_id - let the handler validate
			return nil
		}

		// Get team ID from column
		teamID, err := queries.GetTeamIDByProjectStatusColumn(r.Context(), data.ColumnID)
		if err != nil {
			if err == sql.ErrNoRows {
				// Column doesn't exist - let the handler return 404
				return nil
			}
			return err
		}

		// Check team membership
		isMember, err := queries.CheckUserTeamMembership(r.Context(), db.CheckUserTeamMembershipParams{
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
}

// CheckColumnAccessByURLParam checks if user has access to a project status column via team membership
func CheckColumnAccessByURLParam(paramName string) AccessChecker {
	return func(r *http.Request, queries *db.Queries) error {
		userID, _ := GetUserID(r)

		columnIDStr := chi.URLParam(r, paramName)
		columnID, err := strconv.ParseInt(columnIDStr, 10, 64)
		if err != nil {
			// Invalid ID - let the handler return 400
			return nil
		}

		// Get team ID from column
		teamID, err := queries.GetTeamIDByProjectStatusColumn(r.Context(), columnID)
		if err != nil {
			if err == sql.ErrNoRows {
				// Column doesn't exist - let the handler return 404
				return nil
			}
			return err
		}

		// Check team membership
		isMember, err := queries.CheckUserTeamMembership(r.Context(), db.CheckUserTeamMembershipParams{
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
}

// CheckProjectAccessByBody checks if user has access to a project via team membership from request body
func CheckProjectAccessByBody() AccessChecker {
	return func(r *http.Request, queries *db.Queries) error {
		userID, _ := GetUserID(r)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			return errors.New("failed to read request body")
		}
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		var data struct {
			ProjectID int32 `json:"project_id"`
		}
		if err := json.Unmarshal(body, &data); err != nil {
			// Invalid JSON - let the handler return 400
			return nil
		}

		if data.ProjectID == 0 {
			// Missing or invalid project_id - let the handler validate
			return nil
		}

		// Get team ID from project
		teamID, err := queries.GetTeamIDByProject(r.Context(), int64(data.ProjectID))
		if err != nil {
			if err == sql.ErrNoRows {
				// Project doesn't exist - let the handler return 404
				return nil
			}
			return err
		}

		// Check team membership
		isMember, err := queries.CheckUserTeamMembership(r.Context(), db.CheckUserTeamMembershipParams{
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
}

// CheckConversationOwnershipByBody checks if user owns a conversation from request body
func CheckConversationOwnershipByBody() AccessChecker {
	return func(r *http.Request, queries *db.Queries) error {
		userID, _ := GetUserID(r)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			return errors.New("failed to read request body")
		}
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		var data struct {
			ConversationID int64 `json:"conversation_id"`
		}
		if err := json.Unmarshal(body, &data); err != nil {
			// Invalid JSON - let the handler return 400
			return nil
		}

		if data.ConversationID == 0 {
			// Missing or invalid conversation_id - let the handler validate
			return nil
		}

		// Get conversation
		conversation, err := queries.GetConversationByID(r.Context(), data.ConversationID)
		if err != nil {
			if err == sql.ErrNoRows {
				// Conversation doesn't exist - let the handler return 404
				return nil
			}
			return err
		}

		// Check ownership
		if conversation.UserID != userID {
			return errors.New("user does not own this conversation")
		}

		return nil
	}
}
