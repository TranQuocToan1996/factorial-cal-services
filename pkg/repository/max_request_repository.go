package repository

import (
	"errors"

	"factorial-cal-services/pkg/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// maxRequestRepository implements MaxRequestRepository interface
type maxRequestRepository struct {
	db *gorm.DB
}

// MaxRequestRepository defines the interface for max request number operations
type MaxRequestRepository interface {
	GetMaxNumber() (int64, error)
	UpdateMaxNumber(maxNumber int64) error
	SetMaxNumberIfGreater(maxNumber int64) (int64, error)
}

// NewMaxRequestRepository creates a new max request repository
func NewMaxRequestRepository(db *gorm.DB) MaxRequestRepository {
	return &maxRequestRepository{
		db: db,
	}
}

// GetMaxNumber retrieves the current maximum requested number
func (r *maxRequestRepository) GetMaxNumber() (int64, error) {
	var maxReq domain.FactorialMaxRequestNumber
	result := r.db.Order("max_number DESC").First(&maxReq)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, nil // Return "0" if no record exists
		}
		return 0, result.Error
	}

	return maxReq.MaxNumber, nil
}

// UpdateMaxNumber updates the maximum requested number
func (r *maxRequestRepository) UpdateMaxNumber(maxNumber int64) error {
	// Check if record exists
	var existing domain.FactorialMaxRequestNumber
	db := r.db.Session(&gorm.Session{
		Logger: logger.Discard, // Disable print error when not found
	})
	result := db.First(&existing)

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
func (r *maxRequestRepository) SetMaxNumberIfGreater(maxNumber int64) (int64, error) {
	query := r.db.Model(&domain.FactorialMaxRequestNumber{}).
		Where("max_number < ?", maxNumber).
		Update("max_number", maxNumber)
	rowAffected, err := query.RowsAffected, query.Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return rowAffected, nil
	}
	return rowAffected, err
}
