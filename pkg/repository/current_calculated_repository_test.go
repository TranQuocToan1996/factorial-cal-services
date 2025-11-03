package repository

import (
	"testing"

	"factorial-cal-services/pkg/domain"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDBForCurrentCalculated(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto migrate
	err = db.AutoMigrate(&domain.FactorialCurrentCalculatedNumber{})
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	return db
}

func TestGetCurrentNumber(t *testing.T) {
	db := setupTestDBForCurrentCalculated(t)
	repo := NewCurrentCalculatedRepository(db)

	// Test empty database
	current, err := repo.GetCurrentNumber()
	if err != nil {
		t.Errorf("Failed to get current number: %v", err)
	}
	if current != "0" {
		t.Errorf("Expected '0' for empty database, got %s", current)
	}

	// Create a record
	curCalc := &domain.FactorialCurrentCalculatedNumber{
		CurNumber: "50",
	}
	err = db.Create(curCalc).Error
	if err != nil {
		t.Fatalf("Failed to create record: %v", err)
	}

	// Get current number
	current, err = repo.GetCurrentNumber()
	if err != nil {
		t.Errorf("Failed to get current number: %v", err)
	}
	if current != "50" {
		t.Errorf("Expected '50', got %s", current)
	}
}

func TestUpdateCurrentNumber(t *testing.T) {
	db := setupTestDBForCurrentCalculated(t)
	repo := NewCurrentCalculatedRepository(db)

	// Test creating new record
	err := repo.UpdateCurrentNumber("10")
	if err != nil {
		t.Errorf("Failed to update current number: %v", err)
	}

	current, err := repo.GetCurrentNumber()
	if err != nil {
		t.Errorf("Failed to get current number: %v", err)
	}
	if current != "10" {
		t.Errorf("Expected '10', got %s", current)
	}

	// Test updating existing record
	err = repo.UpdateCurrentNumber("20")
	if err != nil {
		t.Errorf("Failed to update current number: %v", err)
	}

	current, err = repo.GetCurrentNumber()
	if err != nil {
		t.Errorf("Failed to get current number: %v", err)
	}
	if current != "20" {
		t.Errorf("Expected '20', got %s", current)
	}

	// Test updating to higher number
	err = repo.UpdateCurrentNumber("100")
	if err != nil {
		t.Errorf("Failed to update current number: %v", err)
	}

	current, err = repo.GetCurrentNumber()
	if err != nil {
		t.Errorf("Failed to get current number: %v", err)
	}
	if current != "100" {
		t.Errorf("Expected '100', got %s", current)
	}
}
