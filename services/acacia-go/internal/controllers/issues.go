package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"acacia/internal/db"
	"acacia/internal/schemas"

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

func (c *IssuesController) GetAllIssues(w http.ResponseWriter, r *http.Request) {
	issues, err := c.queries.GetAllIssues(r.Context())
	if err != nil {
		c.logger.WithError(err).Error("Failed to get all issues")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(issues)
}

func (c *IssuesController) GetIssueByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid issue ID", http.StatusBadRequest)
		return
	}

	issue, err := c.queries.GetIssueByID(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Issue not found", http.StatusNotFound)
			return
		}
		c.logger.WithError(err).Error("Failed to get issue by ID")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(issue)
}

func (c *IssuesController) CreateIssue(w http.ResponseWriter, r *http.Request) {
	var req schemas.CreateIssueInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	params := db.CreateIssueParams{
		Name:        req.Name,
		Description: stringPtrToNullString(req.Description),
	}

	issue, err := c.queries.CreateIssue(r.Context(), params)
	if err != nil {
		c.logger.WithError(err).Error("Failed to create issue")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(issue)
}

func (c *IssuesController) UpdateIssue(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid issue ID", http.StatusBadRequest)
		return
	}

	var req schemas.UpdateIssueInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	params := db.UpdateIssueParams{
		ID:          id,
		Name:        req.Name,
		Description: stringPtrToNullString(req.Description),
	}

	issue, err := c.queries.UpdateIssue(r.Context(), params)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Issue not found", http.StatusNotFound)
			return
		}
		c.logger.WithError(err).Error("Failed to update issue")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(issue)
}

func (c *IssuesController) DeleteIssue(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid issue ID", http.StatusBadRequest)
		return
	}

	_, err = c.queries.DeleteIssue(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Issue not found", http.StatusNotFound)
			return
		}
		c.logger.WithError(err).Error("Failed to delete issue")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func stringPtrToNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}