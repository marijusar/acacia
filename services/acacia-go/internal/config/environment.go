package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Environment struct {
	Port        string
	DatabaseURL string
}

func LoadEnvironment() *Environment {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		logrus.Fatal("DATABASE_URL environment variable required")
	}

	return &Environment{
		Port:        port,
		DatabaseURL: databaseURL,
	}
}