package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"acacia/packages/db"
	"acacia/packages/httperr"
	"acacia/packages/schemas"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type IssuesController struct {
	queries *db.Queries
	logger  *logrus.Logger
}

func NewIssuesController(queries *db.Queries, logger *logrus.Logger) *IssuesController {
	return &IssuesController{
		queries: queries,
		logger:  logger,
	}
}

func (c *IssuesController) GetAllIssues(w http.ResponseWriter, r *http.Request) error {
	issues, err := c.queries.GetAllIssues(r.Context())
	if err != nil {
		c.logger.WithError(err).Error("Failed to get all issues")
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

	json.NewEncoder(w).Encode(issue)

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
		Description: stringPtrToNullString(req.Description),
	}

	issue, err := c.queries.CreateIssue(r.Context(), params)
	if err != nil {
		c.logger.WithError(err).Error("Failed to create issue")
		return httperr.WithStatus(errors.New("packages server error"), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(issue)
	return nil
}

func (c *IssuesController) UpdateIssue(w http.ResponseWriter, r *http.Request) error {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return httperr.WithStatus(errors.New("Invalid issue ID"), http.StatusBadRequest)
	}

	var req schemas.UpdateIssueInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return httperr.WithStatus(errors.New("Invalid JSON"), http.StatusBadRequest)
	}

	if req.Name == "" {
		return httperr.WithStatus(errors.New("Name is required"), http.StatusBadRequest)
	}

	params := db.UpdateIssueParams{
		ID:          id,
		Name:        req.Name,
		Description: stringPtrToNullString(req.Description),
	}

	issue, err := c.queries.UpdateIssue(r.Context(), params)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperr.WithStatus(errors.New("Issue not found"), http.StatusNotFound)
		}
		c.logger.WithError(err).Error("Failed to update issue")
		return httperr.WithStatus(errors.New("Internal server error"), http.StatusInternalServerError)
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

func stringPtrToNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}
