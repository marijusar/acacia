package routes

import (
	"acacia/packages/api"
	"acacia/packages/httperr"

	"github.com/go-chi/chi/v5"
)

func IssuesRoutes(controller *api.IssuesController) chi.Router {
	r := chi.NewRouter()

	r.Get("/{id}", httperr.WithCustomErrorHandler(controller.GetIssueByID))
	r.Post("/", httperr.WithCustomErrorHandler(controller.CreateIssue))
	r.Put("/", httperr.WithCustomErrorHandler(controller.UpdateIssue))
	r.Delete("/{id}", httperr.WithCustomErrorHandler(controller.DeleteIssue))

	return r
}
