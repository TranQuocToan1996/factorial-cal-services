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
	_ = os.MkdirAll(basePath, 0755)
	return &LocalStorageService{
		basePath: basePath,
	}
}

// GenerateS3Key generates a storage key for a given number (consistent with S3 format)
func (s *LocalStorageService) GenerateS3Key(number string) string {
	return fmt.Sprintf("%s%s.txt", LocalStoragePrefix, number)
}

// Upload saves a factorial result to local filesystem
func (s *LocalStorageService) Upload(ctx context.Context, number string, result string) (string, error) {
	key := s.GenerateS3Key(number)
	filePath := filepath.Join(s.basePath, key)
	
	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}
	
	// Write file
	if err := os.WriteFile(filePath, []byte(result), 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}
	
	return key, nil
}

// Download reads a factorial result from local filesystem
func (s *LocalStorageService) Download(ctx context.Context, key string) (string, error) {
	filePath := filepath.Join(s.basePath, key)
	
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	
	return string(data), nil
}

