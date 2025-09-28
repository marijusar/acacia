package config

import (
	"github.com/sirupsen/logrus"
)

func NewLogger(env string) *logrus.Logger {
	logger := logrus.New()
	if env == EnvProduction {
		logger.SetFormatter(&logrus.JSONFormatter{})
	}

	return logger
}
