package domain

import "time"

// FactorialCurrentCalculatedNumber represents the current calculated factorial number
type FactorialCurrentCalculatedNumber struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	CurNumber int64     `gorm:"type:bigint;not null;index" json:"cur_number"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for GORM
func (FactorialCurrentCalculatedNumber) TableName() string {
	return "factorial_current_calculated_numbers"
}
