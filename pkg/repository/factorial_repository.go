package repository

import (
	"errors"
	"fmt"

	"factorial-cal-services/pkg/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// factorialRepository implements FactorialRepository interface
type factorialRepository struct {
	db *gorm.DB
}

// FactorialRepository defines the interface for factorial data operations
type FactorialRepository interface {
	Create(calc *domain.FactorialCalculation) error
	FindByNumber(number int64) (*domain.FactorialCalculation, error)
	UpdateStatus(number string, status string) error
	UpdateS3Key(number string, s3Key string, status string) error
	UpdateS3KeyWithChecksum(number int64, s3Key string, checksum string, size int64, status string) error
	UpdateWithCurrentNumber(number int64, s3Key string, checksum string, size int64, status string) error
}

// NewFactorialRepository creates a new factorial repository
func NewFactorialRepository(db *gorm.DB) FactorialRepository {
	return &factorialRepository{
		db: db,
	}
}

// Create inserts a new factorial calculation record
func (r *factorialRepository) Create(calc *domain.FactorialCalculation) error {
	return r.db.Create(calc).Error
}

// FindByNumber retrieves a factorial calculation by number
func (r *factorialRepository) FindByNumber(number int64) (*domain.FactorialCalculation, error) {
	var calc domain.FactorialCalculation

	db := r.db.Session(&gorm.Session{
		Logger: logger.Discard, // Disable print error when not found
	})
	result := db.Where("number = ?", number).First(&calc)

	if result.Error != nil {
		return nil, result.Error
	}

	return &calc, nil
}

// UpdateStatus updates the status of a factorial calculation
func (r *factorialRepository) UpdateStatus(number string, status string) error {
	result := r.db.Model(&domain.FactorialCalculation{}).
		Where("number = ?", number).
		Update("status", status)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// UpdateS3Key updates the S3 key and status of a factorial calculation
func (r *factorialRepository) UpdateS3Key(number string, s3Key string, status string) error {
	result := r.db.Model(&domain.FactorialCalculation{}).
		Where("number = ?", number).
		Updates(map[string]any{
			"s3_key": s3Key,
			"status": status,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// UpdateS3KeyWithChecksum updates the S3 key, checksum, size, and status of a factorial calculation
func (r *factorialRepository) UpdateS3KeyWithChecksum(number int64, s3Key string, checksum string, size int64, status string) error {
	result := r.db.Model(&domain.FactorialCalculation{}).
		Where("number = ?", number).
		Updates(map[string]any{
			"s3_key":   s3Key,
			"checksum": checksum,
			"size":     size,
			"status":   status,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// UpdateWithCurrentNumber atomically updates factorial metadata and current calculated number
func (r *factorialRepository) UpdateWithCurrentNumber(
	number int64,
	s3Key string,
	checksum string,
	size int64,
	status string,
) error {
	// Use transaction to ensure atomicity
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Update factorial calculation metadata
		result := tx.Model(&domain.FactorialCalculation{}).
			Where("number = ?", number).
			Updates(map[string]any{
				"s3_key":   s3Key,
				"checksum": checksum,
				"size":     size,
				"status":   status,
			})

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		// Update current calculated number within the same transaction
		var existing domain.FactorialCurrentCalculatedNumber
		if err := tx.First(&existing).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Create new record
				curCalc := &domain.FactorialCurrentCalculatedNumber{
					NextNumber: number + 1,
				}
				if err := tx.Create(curCalc).Error; err != nil {
					return fmt.Errorf("failed to create current number: %w", err)
				}
			} else {
				return fmt.Errorf("failed to query current number: %w", err)
			}
		} else {
			// Update existing record
			if err := tx.Model(&existing).Update("next_number", number+1).Error; err != nil {
				return fmt.Errorf("failed to update current number: %w", err)
			}
		}

		return nil
	})
}
