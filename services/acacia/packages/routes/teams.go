package routes

import (
	"acacia/packages/api"
	"acacia/packages/httperr"

	"github.com/go-chi/chi/v5"
)

func TeamsRoutes(controller *api.TeamsController, authMiddlewares chi.Middlewares) chi.Router {
	r := chi.NewRouter()

	// Apply authentication middleware to all routes
	r.Use(authMiddlewares...)

	r.Post("/", httperr.WithCustomErrorHandler(controller.CreateTeam))
	r.Get("/", httperr.WithCustomErrorHandler(controller.GetUserTeams))

	return r
}
