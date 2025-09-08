package config

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func NewLogger(env string) *logrus.Logger {
	logger := logrus.New()
	fmt.Println(env)
	if env == EnvProduction {
		logger.SetFormatter(&logrus.JSONFormatter{})
	}

	return logger
}
