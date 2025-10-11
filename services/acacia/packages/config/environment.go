package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Environment struct {
	Port        string
	DatabaseURL string
	Env         string
	JWTSecret   string
}

const (
	EnvProduction = "production"
	EnvDev        = "development"
)

func LoadEnvironment() *Environment {
	godotenv.Load()

	port := os.Getenv("PORT")
	env := os.Getenv("ENV")
	if port == "" {
		port = "8080"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		logrus.Fatal("DATABASE_URL environment variable required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		logrus.Fatal("JWT_SECRET environment variable required")
	}

	return &Environment{
		Env:         env,
		Port:        port,
		DatabaseURL: databaseURL,
		JWTSecret:   jwtSecret,
	}
}

