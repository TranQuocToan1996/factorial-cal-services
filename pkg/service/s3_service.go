package service

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	S3KeyPrefix = "factorials/"
)

// S3Service handles S3 storage operations
type S3Service interface {
	UploadFactorial(ctx context.Context, number string, result string) (string, error)
	DownloadFactorial(ctx context.Context, s3Key string) (string, error)
	GenerateS3Key(number string) string
	// StorageService interface methods
	Upload(ctx context.Context, number string, result string) (string, error)
	Download(ctx context.Context, key string) (string, error)
}

type s3Service struct {
	client     *s3.Client
	bucketName string
}

// NewS3Service creates a new S3 service
func NewS3Service(client *s3.Client, bucketName string) S3Service {
	return &s3Service{
		client:     client,
		bucketName: bucketName,
	}
}

// GenerateS3Key generates an S3 key for a given number
func (s *s3Service) GenerateS3Key(number string) string {
	return fmt.Sprintf("%s%s.txt", S3KeyPrefix, number)
}

// UploadFactorial uploads a factorial result to S3
func (s *s3Service) UploadFactorial(ctx context.Context, number string, result string) (string, error) {
	key := s.GenerateS3Key(number)
	
	// Convert result to byte buffer
	body := bytes.NewReader([]byte(result))
	
	// Upload to S3
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String("text/plain"),
	})
	
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}
	
	return key, nil
}

// DownloadFactorial downloads a factorial result from S3
func (s *s3Service) DownloadFactorial(ctx context.Context, s3Key string) (string, error) {
	// Get object from S3
	output, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(s3Key),
	})
	
	if err != nil {
		return "", fmt.Errorf("failed to download from S3: %w", err)
	}
	defer output.Body.Close()
	
	// Read body
	body, err := io.ReadAll(output.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read S3 object body: %w", err)
	}
	
	return string(body), nil
}

// Upload implements StorageService interface (alias for UploadFactorial)
func (s *s3Service) Upload(ctx context.Context, number string, result string) (string, error) {
	return s.UploadFactorial(ctx, number, result)
}

// Download implements StorageService interface (alias for DownloadFactorial)
func (s *s3Service) Download(ctx context.Context, key string) (string, error) {
	return s.DownloadFactorial(ctx, key)
}

