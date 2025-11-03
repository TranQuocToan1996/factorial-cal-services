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
	GetCurrentNumber() (string, error)
	UpdateCurrentNumber(curNumber string) error
}

// NewCurrentCalculatedRepository creates a new current calculated repository
func NewCurrentCalculatedRepository(db *gorm.DB) CurrentCalculatedRepository {
	return &currentCalculatedRepository{
		db: db,
	}
}

// GetCurrentNumber retrieves the current calculated number
func (r *currentCalculatedRepository) GetCurrentNumber() (string, error) {
	var curCalc domain.FactorialCurrentCalculatedNumber
	result := r.db.Order("cur_number DESC").First(&curCalc)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "0", nil // Return "0" if no record exists
		}
		return "", result.Error
	}

	return curCalc.CurNumber, nil
}

// UpdateCurrentNumber updates the current calculated number
func (r *currentCalculatedRepository) UpdateCurrentNumber(curNumber string) error {
	// Check if record exists
	var existing domain.FactorialCurrentCalculatedNumber
	result := r.db.First(&existing)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Create new record
		curCalc := &domain.FactorialCurrentCalculatedNumber{
			CurNumber: curNumber,
		}
		return r.db.Create(curCalc).Error
	}

	if result.Error != nil {
		return result.Error
	}

	// Update existing record
	return r.db.Model(&existing).Update("cur_number", curNumber).Error
}
