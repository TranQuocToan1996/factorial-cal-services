package domain

import "time"

// FactorialMaxRequestNumber represents the maximum requested factorial number
type FactorialMaxRequestNumber struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	MaxNumber int64     `gorm:"type:bigint;not null;index" json:"max_number"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for GORM
func (FactorialMaxRequestNumber) TableName() string {
	return "factorial_max_request_numbers"
}
