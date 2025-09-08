package config

import (
	"acacia/internal/api"
	"acacia/internal/db"
	"acacia/internal/routes"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

type Server struct {
	router *chi.Mux
	db     *db.Queries
	logger *logrus.Logger
}

func NewServer(q *db.Queries, l *logrus.Logger) *Server {
	issuesController := api.NewIssuesController(q, l)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/issues", routes.IssuesRoutes(issuesController))

	server := &Server{
		router: r,
		db:     q,
		logger: l,
	}

	return server
}

func (s *Server) ListenAndServe(port string) {
	s.logger.WithField("port", port).Info("Starting server")
	if err := http.ListenAndServe(":"+port, s.router); err != nil {
		s.logger.WithError(err).Fatal("Server failed to start")
	}
}
