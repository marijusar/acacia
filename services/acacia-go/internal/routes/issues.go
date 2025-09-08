package routes

import (
	"acacia/internal/controllers"

	"github.com/go-chi/chi/v5"
)

func IssuesRoutes(controller *controllers.IssuesController) chi.Router {
	r := chi.NewRouter()

	r.Get("/", controller.GetAllIssues)
	r.Get("/{id}", controller.GetIssueByID)
	r.Post("/", controller.CreateIssue)
	r.Put("/{id}", controller.UpdateIssue)
	r.Delete("/{id}", controller.DeleteIssue)

	return r
}