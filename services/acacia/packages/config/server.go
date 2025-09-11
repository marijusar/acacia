package config

import (
	"acacia/packages/api"
	"acacia/packages/db"
	"acacia/packages/routes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

type Server struct {
	httpServer *http.Server
	db         *db.Queries
	logger     *logrus.Logger
}

func NewServer(q *db.Queries, l *logrus.Logger) *Server {
	issuesController := api.NewIssuesController(q, l)
	projectsController := api.NewProjectsController(q, l)
	projectColumnsController := api.NewProjectStatusColumnsController(q, l)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/issues", routes.IssuesRoutes(issuesController))
	r.Mount("/projects", routes.ProjectsRoutes(projectsController))
	r.Mount("/project-columns", routes.ProjectStatusColumnsRoutes(projectColumnsController))

	httpServer := &http.Server{
		Handler: r,
	}

	server := &Server{
		httpServer: httpServer,
		db:         q,
		logger:     l,
	}

	return server
}

func (s *Server) ListenAndServe(port string) {
	s.httpServer.Addr = ":" + port
	s.logger.WithField("port", port).Info("Starting server")
	fmt.Printf("Starting server at port %s\n", port)
	if err := s.httpServer.ListenAndServe(); err != nil {
		s.logger.WithError(err).Error("Server failed to start")
	}
}

func (s *Server) Close() {
	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*5000)
	s.httpServer.Shutdown(ctx)
}
