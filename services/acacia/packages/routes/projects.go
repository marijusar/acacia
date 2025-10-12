package routes

import (
	"acacia/packages/api"
	"acacia/packages/auth"
	"acacia/packages/httperr"

	"github.com/go-chi/chi/v5"
)

func ProjectsRoutes(controller *api.ProjectsController, authMiddlewares chi.Middlewares, authzMiddleware *auth.AuthorizationMiddleware) chi.Router {
	r := chi.NewRouter()

	// Apply authentication middleware to all routes
	r.Use(authMiddlewares...)

	// Routes that don't require resource-level authorization
	r.Get("/", httperr.WithCustomErrorHandler(controller.GetProjects))

	// POST /projects - check team membership from team_id in request body
	r.Group(func(r chi.Router) {
		r.Use(authzMiddleware.RequireResourceAccessFromBody(auth.ResourceTypeTeam, auth.ExtractTeamIDFromBody))
		r.Post("/", httperr.WithCustomErrorHandler(controller.CreateProject))
	})

	// Routes that require project-level authorization
	r.Group(func(r chi.Router) {
		r.Use(authzMiddleware.RequireResourceAccess(auth.ResourceTypeProject, "id"))
		r.Get("/{id}", httperr.WithCustomErrorHandler(controller.GetProjectByID))
		r.Get("/{id}/details", httperr.WithCustomErrorHandler(controller.GetProjectDetailsByID))
		r.Put("/{id}", httperr.WithCustomErrorHandler(controller.UpdateProject))
		r.Delete("/{id}", httperr.WithCustomErrorHandler(controller.DeleteProject))
	})

	return r
}
