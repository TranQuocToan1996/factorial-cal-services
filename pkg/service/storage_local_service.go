package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

const (
	LocalStoragePrefix = "factorials/"
)

// LocalStorageService handles local filesystem storage operations
type LocalStorageService struct {
	basePath string
}

// NewLocalStorageService creates a new local storage service
func NewLocalStorageService(basePath string) StorageService {
	if basePath == "" {
		basePath = "/tmp/factorial-storage"
	}
	// Ensure directory exists
	_ = os.MkdirAll(basePath, 0o755)
	return &LocalStorageService{
		basePath: basePath,
	}
}

// GenerateKey generates a storage key for a given number (consistent with S3 format)
func (s *LocalStorageService) GenerateKey(number int64) string {
	return fmt.Sprintf("%v%v.txt", LocalStoragePrefix, number)
}

// Upload saves a factorial result to local filesystem
func (s *LocalStorageService) UploadFactorial(ctx context.Context, number int64, result string) (string, error) {
	key := s.GenerateKey(number)
	filePath := filepath.Join(s.basePath, key)

	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Write file
	if err := os.WriteFile(filePath, []byte(result), 0o644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return key, nil
}

// Download reads a factorial result from local filesystem
func (s *LocalStorageService) DownloadFactorial(ctx context.Context, s3Key string) (string, error) {
	filePath := filepath.Join(s.basePath, s3Key)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return string(data), nil
}
