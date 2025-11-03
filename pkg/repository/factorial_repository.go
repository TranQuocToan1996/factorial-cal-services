package repository

import (
	"errors"
	"factorial-cal-services/pkg/domain"

	"gorm.io/gorm"
)

// factorialRepository implements FactorialRepository interface
type factorialRepository struct {
	db *gorm.DB
}

// NewFactorialRepository creates a new factorial repository
func NewFactorialRepository(db *gorm.DB) domain.FactorialRepository {
	return &factorialRepository{
		db: db,
	}
}

// Create inserts a new factorial calculation record
func (r *factorialRepository) Create(calc *domain.FactorialCalculation) error {
	result := r.db.Create(calc)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// FindByNumber retrieves a factorial calculation by number
func (r *factorialRepository) FindByNumber(number string) (*domain.FactorialCalculation, error) {
	var calc domain.FactorialCalculation
	result := r.db.Where("number = ?", number).First(&calc)
	
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
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
		Updates(map[string]interface{}{
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
func (r *factorialRepository) UpdateS3KeyWithChecksum(number string, s3Key string, checksum string, size int64, status string) error {
	result := r.db.Model(&domain.FactorialCalculation{}).
		Where("number = ?", number).
		Updates(map[string]interface{}{
			"s3_key":  s3Key,
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

