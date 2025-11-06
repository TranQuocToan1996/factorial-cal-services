package service

import (
	"context"
)

// StorageService handles storage operations (S3, local filesystem, etc.)
type StorageService interface {
	UploadFactorial(ctx context.Context, number int64, result string) (string, error)
	DownloadFactorial(ctx context.Context, s3Key string) (string, error)
	GenerateKey(number int64) string
	GetBucket() string
}
