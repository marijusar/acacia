package routes

import (
	"acacia/packages/api"
	"acacia/packages/httperr"

	"github.com/go-chi/chi/v5"
)

func ProjectStatusColumnsRoutes(controller *api.ProjectStatusColumnsController) chi.Router {
	r := chi.NewRouter()

	r.Get("/{id}", httperr.WithCustomErrorHandler(controller.GetProjectStatusColumnByID))
	r.Post("/", httperr.WithCustomErrorHandler(controller.CreateProjectStatusColumn))
	r.Put("/{id}", httperr.WithCustomErrorHandler(controller.UpdateProjectStatusColumn))
	r.Delete("/{id}", httperr.WithCustomErrorHandler(controller.DeleteProjectStatusColumn))

	return r
}

func ProjectStatusColumnsByProjectRoutes(controller *api.ProjectStatusColumnsController) chi.Router {
	r := chi.NewRouter()

	r.Get("/{project_id}/columns", httperr.WithCustomErrorHandler(controller.GetProjectStatusColumnsByProjectID))

	return r
}