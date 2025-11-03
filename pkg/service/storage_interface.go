package service

import (
	"context"
)

// StorageService handles storage operations (S3, local filesystem, etc.)
type StorageService interface {
	Upload(ctx context.Context, number string, result string) (string, error)
	Download(ctx context.Context, key string) (string, error)
	GenerateS3Key(number string) string
}

