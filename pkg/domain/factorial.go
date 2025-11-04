package domain

import (
	"time"
)

// Status constants for factorial calculations
const (
	StatusCalculating = "calculating"
	StatusUploading   = "uploading"
	StatusDone        = "done"
	StatusFailed      = "failed"
)

// FactorialCalculation represents a factorial calculation record
type FactorialCalculation struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Number    string    `gorm:"type:varchar(100);not null;uniqueIndex" json:"number"`
	Status    string    `gorm:"type:varchar(20);not null;index" json:"status"`
	S3Key     string    `gorm:"type:varchar(512);not null" json:"s3_key"`
	Checksum  string    `gorm:"type:varchar(64)" json:"checksum,omitempty"`
	Size      int64     `gorm:"type:bigint;default:0" json:"size,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for GORM
func (FactorialCalculation) TableName() string {
	return "factorial_calculations"
}
