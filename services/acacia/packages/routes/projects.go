package routes

import (
	"acacia/packages/api"
	"acacia/packages/httperr"

	"github.com/go-chi/chi/v5"
)

func ProjectsRoutes(controller *api.ProjectsController) chi.Router {
	r := chi.NewRouter()

	r.Get("/", httperr.WithCustomErrorHandler(controller.GetProjects))
	r.Get("/{id}", httperr.WithCustomErrorHandler(controller.GetProjectByID))
	r.Get("/{id}/details", httperr.WithCustomErrorHandler(controller.GetProjectDetailsByID))
	r.Post("/", httperr.WithCustomErrorHandler(controller.CreateProject))
	r.Put("/{id}", httperr.WithCustomErrorHandler(controller.UpdateProject))
	r.Delete("/{id}", httperr.WithCustomErrorHandler(controller.DeleteProject))

	return r
}
