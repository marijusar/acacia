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
		teamIDStr := chi.URLParam(r, paramName)
		teamID, err := strconv.ParseInt(teamIDStr, 10, 64)
		if err != nil {
			// Invalid ID - let the handler return 400
			return nil
		}

		return CheckTeamMembership(r.Context(), queries, teamID)
	}
}

// CheckTeamMembershipByBody checks if user is a member of a team specified in request body
func CheckTeamMembershipByBody(fieldName string) AccessChecker {
	return func(r *http.Request, queries *db.Queries) error {
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

		return CheckTeamMembership(r.Context(), queries, teamID)
	}
}

// CheckProjectAccessByURLParam checks if user has access to a project via team membership
func CheckProjectAccessByURLParam(paramName string) AccessChecker {
	return func(r *http.Request, queries *db.Queries) error {
		projectIDStr := chi.URLParam(r, paramName)
		projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
		if err != nil {
			// Invalid ID - let the handler return 400
			return nil
		}

		return CheckProjectAccess(r.Context(), queries, projectID)
	}
}

// CheckIssueAccessByURLParam checks if user has access to an issue via team membership
func CheckIssueAccessByURLParam(paramName string) AccessChecker {
	return func(r *http.Request, queries *db.Queries) error {
		issueIDStr := chi.URLParam(r, paramName)
		issueID, err := strconv.ParseInt(issueIDStr, 10, 64)
		if err != nil {
			// Invalid ID - let the handler return 400
			return nil
		}

		return CheckIssueAccess(r.Context(), queries, issueID)
	}
}

// CheckIssueAccessByBody checks if user has access to an issue from request body
func CheckIssueAccessByBody() AccessChecker {
	return func(r *http.Request, queries *db.Queries) error {
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

		return CheckIssueAccess(r.Context(), queries, data.ID)
	}
}

// CheckColumnAccessByBody checks if user has access to a project status column from request body
func CheckColumnAccessByBody() AccessChecker {
	return func(r *http.Request, queries *db.Queries) error {
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

		return CheckColumnAccess(r.Context(), queries, data.ColumnID)
	}
}

// CheckColumnAccessByURLParam checks if user has access to a project status column via team membership
func CheckColumnAccessByURLParam(paramName string) AccessChecker {
	return func(r *http.Request, queries *db.Queries) error {
		columnIDStr := chi.URLParam(r, paramName)
		columnID, err := strconv.ParseInt(columnIDStr, 10, 64)
		if err != nil {
			// Invalid ID - let the handler return 400
			return nil
		}

		return CheckColumnAccess(r.Context(), queries, columnID)
	}
}

// CheckProjectAccessByBody checks if user has access to a project via team membership from request body
func CheckProjectAccessByBody() AccessChecker {
	return func(r *http.Request, queries *db.Queries) error {
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

		return CheckProjectAccess(r.Context(), queries, int64(data.ProjectID))
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
