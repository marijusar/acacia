package main

import (
	"database/sql"

	"acacia/packages/config"
	"acacia/packages/db"

	_ "github.com/lib/pq"
)

func main() {
	env := config.LoadEnvironment()

	logger := config.NewLogger(env.Env)

	database, err := sql.Open("postgres", env.DatabaseURL)
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to database")
	}
	defer database.Close()

	if err := database.Ping(); err != nil {
		logger.WithError(err).Fatal("Failed to ping database")
	}

	queries := db.New(database)

	databaseConn := &config.Database{Queries: queries, Conn: database}

	s := config.NewServer(databaseConn, logger)

	s.ListenAndServe(env.Port)
}
