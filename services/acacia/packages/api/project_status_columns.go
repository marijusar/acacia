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
	"acacia/packages/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type ProjectStatusColumnsController struct {
	queries                    *db.Queries
	logger                     *logrus.Logger
	validator                  *validator.Validate
	projectStatusColumnService *services.ProjectStatusColumnService
}

func NewProjectStatusColumnsController(queries *db.Queries, logger *logrus.Logger, database *sql.DB) *ProjectStatusColumnsController {
	return &ProjectStatusColumnsController{
		queries:                    queries,
		logger:                     logger,
		validator:                  validator.New(),
		projectStatusColumnService: services.NewProjectStatusColumnService(queries, database),
	}
}

func (c *ProjectStatusColumnsController) GetProjectStatusColumnsByProjectID(w http.ResponseWriter, r *http.Request) error {
	projectIDStr := chi.URLParam(r, "project_id")
	projectID, err := strconv.ParseInt(projectIDStr, 10, 32)
	if err != nil {
		return httperr.WithStatus(errors.New("Invalid project ID"), http.StatusBadRequest)
	}

	columns, err := c.queries.GetProjectStatusColumnsByProjectID(r.Context(), int32(projectID))
	if err != nil {
		c.logger.WithError(err).Error("Failed to get project status columns by project ID")
		return httperr.WithStatus(errors.New("Internal server error"), http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(columns)
	return nil
}

func (c *ProjectStatusColumnsController) CreateProjectStatusColumn(w http.ResponseWriter, r *http.Request) error {
	var req schemas.CreateProjectStatusColumnInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return httperr.WithStatus(errors.New("Invalid JSON"), http.StatusBadRequest)
	}

	if err := c.validator.Struct(&req); err != nil {
		return httperr.WithStatus(errors.New("Validation failed: "+err.Error()), http.StatusBadRequest)
	}

	params := db.CreateProjectStatusColumnParams{
		ProjectID: req.ProjectID,
		Name:      req.Name,
	}

	column, err := c.queries.CreateProjectStatusColumn(r.Context(), params)
	if err != nil {
		c.logger.WithError(err).Error("Failed to create project status column")
		return httperr.WithStatus(errors.New("Internal server error"), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(column)
	return nil
}

func (c *ProjectStatusColumnsController) UpdateProjectStatusColumn(w http.ResponseWriter, r *http.Request) error {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return httperr.WithStatus(errors.New("Invalid column ID"), http.StatusBadRequest)
	}

	var req schemas.UpdateProjectStatusColumnInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return httperr.WithStatus(errors.New("Invalid JSON"), http.StatusBadRequest)
	}

	if err := c.validator.Struct(&req); err != nil {
		return httperr.WithStatus(errors.New("Validation failed: "+err.Error()), http.StatusBadRequest)
	}

	params := db.UpdateProjectStatusColumnParams{
		ID:            id,
		Name:          req.Name,
		PositionIndex: req.PositionIndex,
	}

	column, err := c.queries.UpdateProjectStatusColumn(r.Context(), params)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperr.WithStatus(errors.New("Project status column not found"), http.StatusNotFound)
		}
		c.logger.WithError(err).Error("Failed to update project status column")
		return httperr.WithStatus(errors.New("Internal server error"), http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(column)
	return nil
}

func (c *ProjectStatusColumnsController) DeleteProjectStatusColumn(w http.ResponseWriter, r *http.Request) error {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return httperr.WithStatus(errors.New("Invalid column ID"), http.StatusBadRequest)
	}

	_, err = c.projectStatusColumnService.DeleteProjectStatusColumnWithReorder(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperr.WithStatus(errors.New("Project status column not found"), http.StatusNotFound)
		}
		if errors.Is(err, services.ErrCannotDeleteLastColumn) {
			return httperr.WithStatus(err, http.StatusBadRequest)
		}
		c.logger.WithError(err).Error("Failed to delete project status column")
		return httperr.WithStatus(errors.New("Internal server error"), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
