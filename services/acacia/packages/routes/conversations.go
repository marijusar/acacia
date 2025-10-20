package routes

import (
	"acacia/packages/api"
	"acacia/packages/auth"
	"acacia/packages/httperr"

	"github.com/go-chi/chi/v5"
)

func ConversationsRoutes(
	controller *api.ConversationsController,
	authMiddlewares chi.Middlewares,
	authzMiddleware *auth.AuthorizationMiddleware,
) chi.Router {
	r := chi.NewRouter()

	// Apply authentication middleware to all routes
	r.Use(authMiddlewares...)

	// POST /conversations - create new conversation (need access to the project to use team's API key)
	r.Group(func(r chi.Router) {
		r.Use(authzMiddleware.RequireAccess(auth.CheckProjectAccessByBody()))
		r.Post("/", httperr.WithCustomErrorHandler(controller.CreateConversation))
	})

	// POST /conversations/messages - send message to conversation (check conversation ownership)
	r.Group(func(r chi.Router) {
		r.Use(authzMiddleware.RequireAccess(auth.CheckConversationOwnershipByBody()))
		r.Post("/messages", httperr.WithCustomErrorHandler(controller.SendMessage))
	})

	r.Get("/latest", httperr.WithCustomErrorHandler(controller.GetLatestConversation))

	return r
}
