package service

import (
	"crypto/sha256"
	"fmt"
)

// ChecksumService handles checksum calculation
type ChecksumService interface {
	Calculate(data string) string
}

type checksumService struct{}

// NewChecksumService creates a new checksum service
func NewChecksumService() ChecksumService {
	return &checksumService{}
}

// Calculate calculates SHA256 checksum for the given data
func (s *checksumService) Calculate(data string) string {
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}

