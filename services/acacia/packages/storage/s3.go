package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sirupsen/logrus"
)

type S3Storage struct {
	client *s3.Client
	bucket string
	logger *logrus.Logger
}

type S3Config struct {
	Bucket          string
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	Endpoint        string // S3 endpoint URL (for localstack: http://localstack:4566, for AWS: https://s3.amazonaws.com)
}

// NewS3Storage creates a new S3Storage instance
func NewS3Storage(cfg S3Config, logger *logrus.Logger) (*S3Storage, error) {
	// Always use custom endpoint resolver
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:               cfg.Endpoint,
			HostnameImmutable: true,
			Source:            aws.EndpointSourceCustom,
		}, nil
	})

	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.Region),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKeyID,
			cfg.SecretAccessKey,
			"",
		)),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create S3 client with path-style addressing
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	return &S3Storage{
		client: client,
		bucket: cfg.Bucket,
		logger: logger,
	}, nil
}

// UploadDescription uploads the serialized description to S3
func (s *S3Storage) UploadDescription(ctx context.Context, issueID int64, content string) error {
	key := s.getDescriptionKey(issueID)

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader([]byte(content)),
		ContentType: aws.String("application/json"),
	})

	if err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"issue_id": issueID,
			"key":      key,
		}).Error("Failed to upload description to S3")
		return fmt.Errorf("failed to upload description to S3: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"issue_id": issueID,
		"key":      key,
	}).Info("Successfully uploaded description to S3")

	return nil
}

// GetDescription retrieves the serialized description from S3
func (s *S3Storage) GetDescription(ctx context.Context, issueID int64) (string, error) {
	key := s.getDescriptionKey(issueID)

	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"issue_id": issueID,
			"key":      key,
		}).Error("Failed to get description from S3")
		return "", fmt.Errorf("failed to get description from S3: %w", err)
	}
	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		s.logger.WithError(err).Error("Failed to read S3 object body")
		return "", fmt.Errorf("failed to read S3 object body: %w", err)
	}

	return string(body), nil
}

// getDescriptionKey generates the S3 key for an issue's description
func (s *S3Storage) getDescriptionKey(issueID int64) string {
	return fmt.Sprintf("issues/%d/description.json", issueID)
}
