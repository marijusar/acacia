package config

import (
	"acacia/packages/api"
	"acacia/packages/auth"
	"acacia/packages/db"
	"acacia/packages/routes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

type Database struct {
	Queries *db.Queries
	Conn    *sql.DB
}

type Server struct {
	httpServer *http.Server
	db         *Database
	logger     *logrus.Logger
}

func NewServer(d *Database, l *logrus.Logger, env *Environment) *Server {
	// Initialize JWT manager
	jwtManager := auth.NewJWTManager(
		env.JWTSecret,
		15*time.Minute, // Access token: 15 minutes
		30*24*time.Hour, // Refresh token: 30 days
	)

	issuesController := api.NewIssuesController(d.Queries, l)
	projectsController := api.NewProjectsController(d.Queries, l)
	projectColumnsController := api.NewProjectStatusColumnsController(d.Queries, l, d.Conn)
	usersController := api.NewUsersController(d.Queries, l, jwtManager)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/issues", routes.IssuesRoutes(issuesController))
	r.Mount("/projects", routes.ProjectsRoutes(projectsController))
	r.Mount("/project-columns", routes.ProjectStatusColumnsRoutes(projectColumnsController))
	r.Mount("/users", routes.UsersRoutes(usersController))

	httpServer := &http.Server{
		Handler: r,
	}

	server := &Server{
		httpServer: httpServer,
		db:         d,
		logger:     l,
	}

	return server
}

func (s *Server) ListenAndServe(port string) {
	s.httpServer.Addr = ":" + port
	s.logger.WithField("port", port).Info("Starting server")
	fmt.Printf("Starting server at port %s\n", port)
	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.WithError(err).Error("Server failed to start")
	}
}

func (s *Server) Close() {
	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*5000)
	s.httpServer.Shutdown(ctx)
}
