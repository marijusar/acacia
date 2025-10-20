package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Environment struct {
	Port            string
	DatabaseURL     string
	Env             string
	JWTSecret       string
	EncryptionKey   []byte
	AWSS3Bucket     string
	AWSRegion       string
	AWSAccessKeyID  string
	AWSSecretKey    string
	AWSEndpoint     string // For localstack
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

	encryptionKey := os.Getenv("ENCRYPTION_KEY")
	if encryptionKey == "" {
		logrus.Fatal("ENCRYPTION_KEY environment variable required (must be 32 bytes)")
	}

	if len(encryptionKey) != 32 {
		logrus.Fatal("ENCRYPTION_KEY must be exactly 32 bytes for AES-256")
	}

	awsS3Bucket := os.Getenv("AWS_S3_BUCKET")
	if awsS3Bucket == "" {
		logrus.Fatal("AWS_S3_BUCKET environment variable required")
	}

	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		logrus.Fatal("AWS_REGION environment variable required")
	}

	awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	if awsAccessKeyID == "" {
		logrus.Fatal("AWS_ACCESS_KEY_ID environment variable required")
	}

	awsSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if awsSecretKey == "" {
		logrus.Fatal("AWS_SECRET_ACCESS_KEY environment variable required")
	}

	// Optional: for localstack development
	awsEndpoint := os.Getenv("AWS_ENDPOINT")

	return &Environment{
		Env:            env,
		Port:           port,
		DatabaseURL:    databaseURL,
		JWTSecret:      jwtSecret,
		EncryptionKey:  []byte(encryptionKey),
		AWSS3Bucket:    awsS3Bucket,
		AWSRegion:      awsRegion,
		AWSAccessKeyID: awsAccessKeyID,
		AWSSecretKey:   awsSecretKey,
		AWSEndpoint:    awsEndpoint,
	}
}

