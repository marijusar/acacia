package routes

import (
	"acacia/packages/api"
	"acacia/packages/auth"
	"acacia/packages/httperr"

	"github.com/go-chi/chi/v5"
)

func ProjectStatusColumnsRoutes(controller *api.ProjectStatusColumnsController, authMiddlewares chi.Middlewares, authzMiddleware *auth.AuthorizationMiddleware) chi.Router {
	r := chi.NewRouter()

	// Apply authentication middleware to all routes
	r.Use(authMiddlewares...)

	// POST /project-columns - check access to project_id from request body
	r.Group(func(r chi.Router) {
		r.Use(authzMiddleware.RequireAccess(auth.CheckProjectAccessByBody()))
		r.Post("/", httperr.WithCustomErrorHandler(controller.CreateProjectStatusColumn))
	})

	// Routes that require project status column-level authorization
	r.Group(func(r chi.Router) {
		r.Use(authzMiddleware.RequireAccess(auth.CheckColumnAccessByURLParam("id")))
		r.Put("/{id}", httperr.WithCustomErrorHandler(controller.UpdateProjectStatusColumn))
		r.Delete("/{id}", httperr.WithCustomErrorHandler(controller.DeleteProjectStatusColumn))
	})

	// This route uses project_id parameter, so needs project-level authorization
	r.Group(func(r chi.Router) {
		r.Use(authzMiddleware.RequireAccess(auth.CheckProjectAccessByURLParam("project_id")))
		r.Get("/project/{project_id}", httperr.WithCustomErrorHandler(controller.GetProjectStatusColumnsByProjectID))
	})

	return r
}
