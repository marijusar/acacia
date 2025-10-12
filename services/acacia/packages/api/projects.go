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
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type ProjectsController struct {
	queries   *db.Queries
	logger    *logrus.Logger
	validator *validator.Validate
}

func NewProjectsController(queries *db.Queries, logger *logrus.Logger) *ProjectsController {
	return &ProjectsController{
		queries:   queries,
		logger:    logger,
		validator: validator.New(),
	}
}

func (c *ProjectsController) GetProjectByID(w http.ResponseWriter, r *http.Request) error {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return httperr.WithStatus(errors.New("Invalid project ID"), http.StatusBadRequest)
	}

	project, err := c.queries.GetProjectByID(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperr.WithStatus(errors.New("Project not found"), http.StatusNotFound)
		}
		c.logger.WithError(err).Error("Failed to get project by ID")
		return httperr.WithStatus(errors.New("Internal server error"), http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(project)
	return nil
}

func (c *ProjectsController) CreateProject(w http.ResponseWriter, r *http.Request) error {
	var req schemas.CreateProjectInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return httperr.WithStatus(errors.New("Invalid JSON"), http.StatusBadRequest)
	}

	if err := c.validator.Struct(&req); err != nil {
		return httperr.WithStatus(errors.New("Validation failed: "+err.Error()), http.StatusBadRequest)
	}

	params := db.CreateProjectParams{
		Name:   req.Name,
		TeamID: req.TeamID,
	}

	project, err := c.queries.CreateProject(r.Context(), params)
	if err != nil {
		c.logger.WithError(err).Error("Failed to create project")
		return httperr.WithStatus(errors.New("Internal server error"), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(project)
	return nil
}

func (c *ProjectsController) UpdateProject(w http.ResponseWriter, r *http.Request) error {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return httperr.WithStatus(errors.New("Invalid project ID"), http.StatusBadRequest)
	}

	var req schemas.UpdateProjectInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return httperr.WithStatus(errors.New("Invalid JSON"), http.StatusBadRequest)
	}

	if err := c.validator.Struct(&req); err != nil {
		return httperr.WithStatus(errors.New("Validation failed: "+err.Error()), http.StatusBadRequest)
	}

	params := db.UpdateProjectParams{
		ID:   id,
		Name: req.Name,
	}

	project, err := c.queries.UpdateProject(r.Context(), params)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperr.WithStatus(errors.New("Project not found"), http.StatusNotFound)
		}
		c.logger.WithError(err).Error("Failed to update project")
		return httperr.WithStatus(errors.New("Internal server error"), http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(project)
	return nil
}

func (c *ProjectsController) GetProjectDetailsByID(w http.ResponseWriter, r *http.Request) error {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return httperr.WithStatus(errors.New("Invalid project ID"), http.StatusBadRequest)
	}

	project, err := c.queries.GetProjectByID(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperr.WithStatus(errors.New("Project not found"), http.StatusNotFound)
		}
		c.logger.WithError(err).Error("Failed to get project details by ID")
		return httperr.WithStatus(errors.New("Internal server error"), http.StatusInternalServerError)
	}

	projectColumns, err := c.queries.GetProjectStatusColumnsByProjectID(r.Context(), int32(id))

	projectIssues, err := c.queries.GetProjectIssues(r.Context(), int32(id))

	resp := schemas.GetProjectDetailsResponse{
		Project: project,
		Columns: projectColumns,
		Issues:  projectIssues,
	}

	json.NewEncoder(w).Encode(resp)
	return nil
}

func (c *ProjectsController) GetProjects(w http.ResponseWriter, r *http.Request) error {
	projects, err := c.queries.GetProjects(r.Context())
	if err != nil {
		c.logger.WithError(err).Error("Failed to get projects")
		return httperr.WithStatus(errors.New("Internal server error"), http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(projects)
	return nil
}

func (c *ProjectsController) DeleteProject(w http.ResponseWriter, r *http.Request) error {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return httperr.WithStatus(errors.New("Invalid project ID"), http.StatusBadRequest)
	}

	_, err = c.queries.DeleteProject(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperr.WithStatus(errors.New("Project not found"), http.StatusNotFound)
		}
		c.logger.WithError(err).Error("Failed to delete project")
		return httperr.WithStatus(errors.New("Internal server error"), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
