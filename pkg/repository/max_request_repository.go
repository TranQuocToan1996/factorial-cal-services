package repository

import (
	"errors"

	"factorial-cal-services/pkg/domain"

	"gorm.io/gorm"
)

// maxRequestRepository implements MaxRequestRepository interface
type maxRequestRepository struct {
	db *gorm.DB
}

// MaxRequestRepository defines the interface for max request number operations
type MaxRequestRepository interface {
	GetMaxNumber() (string, error)
	UpdateMaxNumber(maxNumber string) error
	SetMaxNumberIfGreater(maxNumber string) error
}

// NewMaxRequestRepository creates a new max request repository
func NewMaxRequestRepository(db *gorm.DB) MaxRequestRepository {
	return &maxRequestRepository{
		db: db,
	}
}

// GetMaxNumber retrieves the current maximum requested number
func (r *maxRequestRepository) GetMaxNumber() (string, error) {
	var maxReq domain.FactorialMaxRequestNumber
	result := r.db.Order("max_number DESC").First(&maxReq)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "0", nil // Return "0" if no record exists
		}
		return "", result.Error
	}

	return maxReq.MaxNumber, nil
}

// UpdateMaxNumber updates the maximum requested number
func (r *maxRequestRepository) UpdateMaxNumber(maxNumber string) error {
	// Check if record exists
	var existing domain.FactorialMaxRequestNumber
	result := r.db.First(&existing)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Create new record
		maxReq := &domain.FactorialMaxRequestNumber{
			MaxNumber: maxNumber,
		}
		return r.db.Create(maxReq).Error
	}

	if result.Error != nil {
		return result.Error
	}

	// Update existing record
	return r.db.Model(&existing).Update("max_number", maxNumber).Error
}

// SetMaxNumberIfGreater updates the max number only if the new value is greater
func (r *maxRequestRepository) SetMaxNumberIfGreater(maxNumber string) error {
	currentMax, err := r.GetMaxNumber()
	if err != nil {
		return err
	}

	// Compare as strings (numbers as strings)
	// Simple comparison: if new number is longer, it's greater
	// If same length, compare lexicographically
	if len(maxNumber) > len(currentMax) ||
		(len(maxNumber) == len(currentMax) && maxNumber > currentMax) {
		return r.UpdateMaxNumber(maxNumber)
	}

	return nil
}
