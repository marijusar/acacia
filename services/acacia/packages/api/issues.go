package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"acacia/packages/db"
	"acacia/packages/httperr"
	"acacia/packages/schemas"

	"github.com/go-chi/chi/v5"
	"github.com/guregu/null"
	"github.com/sirupsen/logrus"
)

type IssuesController struct {
	queries *db.Queries
	logger  *logrus.Logger
	storage S3Storage
}

type S3Storage interface {
	UploadDescription(ctx context.Context, issueID int64, content string) error
	GetDescription(ctx context.Context, issueID int64) (string, error)
}

func NewIssuesController(queries *db.Queries, logger *logrus.Logger, storage S3Storage) *IssuesController {
	return &IssuesController{
		queries: queries,
		logger:  logger,
		storage: storage,
	}
}

func (c *IssuesController) GetIssuesByColumnId(w http.ResponseWriter, r *http.Request) error {
	columnIdStr := chi.URLParam(r, "columnId")
	columnId, err := strconv.ParseInt(columnIdStr, 10, 64)
	if err != nil {
		return httperr.WithStatus(errors.New("Invalid column ID"), http.StatusBadRequest)
	}

	issues, err := c.queries.GetIssuesByColumnId(r.Context(), columnId)
	if err != nil {
		c.logger.WithError(err).Error("Failed to get issues by column ID")
		return httperr.WithStatus(errors.New("packages server error"), http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(issues)
	return nil
}

func (c *IssuesController) GetIssueByID(w http.ResponseWriter, r *http.Request) error {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return httperr.WithStatus(errors.New("Invalid issue ID"), http.StatusBadRequest)
	}

	issue, err := c.queries.GetIssueByID(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperr.WithStatus(errors.New("Issue not found"), http.StatusNotFound)
		}
		c.logger.WithError(err).Error("Failed to get issue by ID")
		return httperr.WithStatus(errors.New("packages server error"), http.StatusInternalServerError)
	}

	// Try to fetch serialized description from S3
	var descriptionSerialized *string
	serialized, err := c.storage.GetDescription(r.Context(), id)
	if err != nil {
		// Log error but don't fail the request if S3 fetch fails
		c.logger.WithError(err).Warn("Failed to fetch serialized description from S3")
	} else if serialized != "" {
		descriptionSerialized = &serialized
	}

	// Create response with serialized description
	response := map[string]interface{}{
		"id":                      issue.ID,
		"name":                    issue.Name,
		"description":             issue.Description,
		"column_id":               issue.ColumnID,
		"created_at":              issue.CreatedAt,
		"updated_at":              issue.UpdatedAt,
		"description_serialized":  descriptionSerialized,
	}

	json.NewEncoder(w).Encode(response)

	return nil
}

func (c *IssuesController) CreateIssue(w http.ResponseWriter, r *http.Request) error {
	var req schemas.CreateIssueInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return httperr.WithStatus(errors.New("Invalid JSON"), http.StatusBadRequest)
	}

	if req.Name == "" {
		return httperr.WithStatus(errors.New("Name is required"), http.StatusBadRequest)
	}

	params := db.CreateIssueParams{
		Name:        req.Name,
		ColumnID:    req.ColumnId,
		Description: null.NewString(*req.Description, true),
	}

	issue, err := c.queries.CreateIssue(r.Context(), params)
	if err != nil {
		c.logger.WithError(err).Error("Failed to create issue")
		return httperr.WithStatus(errors.New("packages server error"), http.StatusInternalServerError)
	}

	// Upload serialized description to S3 if provided
	if req.DescriptionSerialized != nil && *req.DescriptionSerialized != "" {
		if err := c.storage.UploadDescription(r.Context(), issue.ID, *req.DescriptionSerialized); err != nil {
			// S3 upload failed - should we rollback? For now, we'll delete the issue
			c.queries.DeleteIssue(r.Context(), issue.ID)
			c.logger.WithError(err).Error("Failed to upload description to S3, rolled back issue creation")
			return httperr.WithStatus(errors.New("Failed to save issue description"), http.StatusInternalServerError)
		}
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(issue)
	return nil
}

func (c *IssuesController) UpdateIssue(w http.ResponseWriter, r *http.Request) error {
	var req schemas.UpdateIssueInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return httperr.WithStatus(errors.New("Invalid JSON"), http.StatusBadRequest)
	}
	fmt.Println(req)

	params := db.UpdateIssueParams{
		ID:          req.ID,
		Name:        req.Name,
		Description: null.NewString(req.Description, req.Description != ""),
		ColumnID:    req.ColumnId,
	}
	fmt.Println(params)

	issue, err := c.queries.UpdateIssue(r.Context(), params)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperr.WithStatus(errors.New("Issue not found"), http.StatusNotFound)
		}
		c.logger.WithError(err).Error("Failed to update issue")
		return httperr.WithStatus(errors.New("Internal server error"), http.StatusInternalServerError)
	}

	// Upload serialized description to S3 if provided
	if req.DescriptionSerialized != nil && *req.DescriptionSerialized != "" {
		if err := c.storage.UploadDescription(r.Context(), issue.ID, *req.DescriptionSerialized); err != nil {
			c.logger.WithError(err).Error("Failed to upload description to S3 during update")
			return httperr.WithStatus(errors.New("Failed to save issue description"), http.StatusInternalServerError)
		}
	}

	json.NewEncoder(w).Encode(issue)
	return nil
}

func (c *IssuesController) DeleteIssue(w http.ResponseWriter, r *http.Request) error {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return httperr.WithStatus(errors.New("Invalid issue ID"), http.StatusBadRequest)
	}

	err = c.queries.DeleteIssue(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperr.WithStatus(errors.New("Issue not found"), http.StatusNotFound)
		}
		c.logger.WithError(err).Error("Failed to delete issue")
		return httperr.WithStatus(errors.New("Internal server error"), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
