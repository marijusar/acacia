package routes

import (
	"acacia/packages/api"
	"acacia/packages/httperr"

	"github.com/go-chi/chi/v5"
)

func UsersRoutes(controller *api.UsersController) chi.Router {
	r := chi.NewRouter()

	r.Post("/register", httperr.WithCustomErrorHandler(controller.Register))
	r.Post("/login", httperr.WithCustomErrorHandler(controller.Login))

	return r
}
