package routes

import (
	"acacia/packages/api"
	"acacia/packages/httperr"

	"github.com/go-chi/chi/v5"
)

func UsersRoutes(controller *api.UsersController, authMiddlewares chi.Middlewares) chi.Router {
	r := chi.NewRouter()

	// Public routes
	r.Post("/register", httperr.WithCustomErrorHandler(controller.Register))
	r.Post("/login", httperr.WithCustomErrorHandler(controller.Login))

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(authMiddlewares...)
		r.Get("/auth/me", httperr.WithCustomErrorHandler(controller.GetAuthStatus))
	})

	return r
}
