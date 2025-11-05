package repository

import (
	"errors"

	"factorial-cal-services/pkg/domain"

	"gorm.io/gorm"
)

// currentCalculatedRepository implements CurrentCalculatedRepository interface
type currentCalculatedRepository struct {
	db *gorm.DB
}

// CurrentCalculatedRepository defines the interface for current calculated number operations
type CurrentCalculatedRepository interface {
	GetCurrentNumber() (int64, error)
	UpdateCurrentNumber(curNumber int64) error
}

// NewCurrentCalculatedRepository creates a new current calculated repository
func NewCurrentCalculatedRepository(db *gorm.DB) CurrentCalculatedRepository {
	return &currentCalculatedRepository{
		db: db,
	}
}

// GetCurrentNumber retrieves the current calculated number
func (r *currentCalculatedRepository) GetCurrentNumber() (int64, error) {
	var curCalc domain.FactorialCurrentCalculatedNumber
	result := r.db.Order("next_number DESC").First(&curCalc)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, nil // Return "0" if no record exists
		}
		return 0, result.Error
	}
	return curCalc.NextNumber, nil
}

// UpdateCurrentNumber updates the current calculated number
func (r *currentCalculatedRepository) UpdateCurrentNumber(curNumber int64) error {
	// Check if record exists
	var existing domain.FactorialCurrentCalculatedNumber
	result := r.db.Order("next_number DESC").First(&existing).Limit(1)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Create new record
		curCalc := &domain.FactorialCurrentCalculatedNumber{
			NextNumber: curNumber,
		}
		return r.db.Create(curCalc).Error
	}

	if result.Error != nil {
		return result.Error
	}

	// Update existing record
	return r.db.Model(&existing).Update("next_number", curNumber).Error
}
