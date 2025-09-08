package main

import (
	"database/sql"
	"net/http"

	"acacia/internal/config"
	"acacia/internal/controllers"
	"acacia/internal/db"
	"acacia/internal/routes"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	env := config.LoadEnvironment()

	database, err := sql.Open("postgres", env.DatabaseURL)
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to database")
	}
	defer database.Close()

	if err := database.Ping(); err != nil {
		logger.WithError(err).Fatal("Failed to ping database")
	}

	queries := db.New(database)
	issuesController := controllers.NewIssuesController(queries, logger)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/issues", routes.IssuesRoutes(issuesController))

	logger.WithField("port", env.Port).Info("Starting server")
	if err := http.ListenAndServe(":"+env.Port, r); err != nil {
		logger.WithError(err).Fatal("Server failed to start")
	}
}

