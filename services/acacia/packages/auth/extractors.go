package auth

import (
	"encoding/json"
	"errors"
)

// Common body extractors for authorization middleware

// ExtractProjectIDFromBody extracts project_id from request body
func ExtractProjectIDFromBody(body []byte) (int64, error) {
	var data struct {
		ProjectID int32 `json:"project_id"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, err
	}
	if data.ProjectID == 0 {
		return 0, errors.New("project_id is required")
	}
	return int64(data.ProjectID), nil
}

// ExtractColumnIDFromBody extracts column_id from request body
func ExtractColumnIDFromBody(body []byte) (int64, error) {
	var data struct {
		ColumnID int64 `json:"column_id"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, err
	}
	if data.ColumnID == 0 {
		return 0, errors.New("column_id is required")
	}
	return data.ColumnID, nil
}

// ExtractIssueIDFromBody extracts id from request body (for update operations)
func ExtractIssueIDFromBody(body []byte) (int64, error) {
	var data struct {
		ID int64 `json:"id"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, err
	}
	if data.ID == 0 {
		return 0, errors.New("id is required")
	}
	return data.ID, nil
}

// ExtractTeamIDFromBody extracts team_id from request body
func ExtractTeamIDFromBody(body []byte) (int64, error) {
	var data struct {
		TeamID int64 `json:"team_id"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, err
	}
	if data.TeamID == 0 {
		return 0, errors.New("team_id is required")
	}
	return data.TeamID, nil
}
