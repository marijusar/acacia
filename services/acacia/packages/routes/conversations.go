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

	// POST /conversations - create new conversation (no specific resource auth needed, just authenticated)
	r.Post("/", httperr.WithCustomErrorHandler(controller.CreateConversation))

	// POST /conversations/messages - send message to conversation (check conversation ownership)
	r.Group(func(r chi.Router) {
		r.Use(authzMiddleware.RequireAccess(auth.CheckConversationOwnershipByBody()))
		r.Post("/messages", httperr.WithCustomErrorHandler(controller.SendMessage))
	})

	// GET /conversations/latest - get latest conversation (no specific resource auth, just authenticated)
	r.Get("/latest", httperr.WithCustomErrorHandler(controller.GetLatestConversation))

	return r
}
