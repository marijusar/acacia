package routes

import (
	"acacia/packages/api"
	"acacia/packages/auth"
	"acacia/packages/httperr"

	"github.com/go-chi/chi/v5"
)

func TeamsRoutes(
	controller *api.TeamsController,
	teamLLMAPIKeysController *api.TeamLLMAPIKeysController,
	authMiddlewares chi.Middlewares,
	authzMiddleware *auth.AuthorizationMiddleware,
) chi.Router {
	r := chi.NewRouter()

	// Apply authentication middleware to all routes
	r.Use(authMiddlewares...)

	r.Post("/", httperr.WithCustomErrorHandler(controller.CreateTeam))
	r.Get("/", httperr.WithCustomErrorHandler(controller.GetUserTeams))

	// Team LLM API Keys routes - require team-level authorization
	r.Group(func(r chi.Router) {
		r.Use(authzMiddleware.RequireResourceAccess(auth.ResourceTypeTeam, "id"))
		r.Post("/{id}/llm-api-keys", httperr.WithCustomErrorHandler(teamLLMAPIKeysController.CreateOrUpdateAPIKey))
		r.Get("/{id}/llm-api-keys", httperr.WithCustomErrorHandler(teamLLMAPIKeysController.GetAPIKeys))
		r.Delete("/{id}/llm-api-keys/{keyId}", httperr.WithCustomErrorHandler(teamLLMAPIKeysController.DeleteAPIKey))
	})

	return r
}
