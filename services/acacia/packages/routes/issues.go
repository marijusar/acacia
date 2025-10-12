package routes

import (
	"acacia/packages/api"
	"acacia/packages/auth"
	"acacia/packages/httperr"

	"github.com/go-chi/chi/v5"
)

func IssuesRoutes(controller *api.IssuesController, authMiddlewares chi.Middlewares, authzMiddleware *auth.AuthorizationMiddleware) chi.Router {
	r := chi.NewRouter()

	// Apply authentication middleware to all routes
	r.Use(authMiddlewares...)

	// POST /issues - check access to the column_id from request body
	r.Group(func(r chi.Router) {
		r.Use(authzMiddleware.RequireResourceAccessFromBody(auth.ResourceTypeProjectStatusColumn, auth.ExtractColumnIDFromBody))
		r.Post("/", httperr.WithCustomErrorHandler(controller.CreateIssue))
	})

	// PUT /issues - check access to the issue id from request body
	r.Group(func(r chi.Router) {
		r.Use(authzMiddleware.RequireResourceAccessFromBody(auth.ResourceTypeIssue, auth.ExtractIssueIDFromBody))
		r.Put("/", httperr.WithCustomErrorHandler(controller.UpdateIssue))
	})

	// Routes that require issue-level authorization via URL parameter
	r.Group(func(r chi.Router) {
		r.Use(authzMiddleware.RequireResourceAccess(auth.ResourceTypeIssue, "id"))
		r.Get("/{id}", httperr.WithCustomErrorHandler(controller.GetIssueByID))
		r.Delete("/{id}", httperr.WithCustomErrorHandler(controller.DeleteIssue))
	})

	return r
}
