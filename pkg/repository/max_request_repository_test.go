package repository

import (
	"testing"

	"factorial-cal-services/pkg/domain"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDBForMaxRequest(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto migrate
	err = db.AutoMigrate(&domain.FactorialMaxRequestNumber{})
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	return db
}

func TestGetMaxNumber(t *testing.T) {
	db := setupTestDBForMaxRequest(t)
	repo := NewMaxRequestRepository(db)

	// Test empty database
	max, err := repo.GetMaxNumber()
	if err != nil {
		t.Errorf("Failed to get max number: %v", err)
	}
	if max != "0" {
		t.Errorf("Expected '0' for empty database, got %s", max)
	}

	// Create a record
	maxReq := &domain.FactorialMaxRequestNumber{
		MaxNumber: "100",
	}
	err = db.Create(maxReq).Error
	if err != nil {
		t.Fatalf("Failed to create record: %v", err)
	}

	// Get max number
	max, err = repo.GetMaxNumber()
	if err != nil {
		t.Errorf("Failed to get max number: %v", err)
	}
	if max != "100" {
		t.Errorf("Expected '100', got %s", max)
	}
}

func TestUpdateMaxNumber(t *testing.T) {
	db := setupTestDBForMaxRequest(t)
	repo := NewMaxRequestRepository(db)

	// Test creating new record
	err := repo.UpdateMaxNumber("50")
	if err != nil {
		t.Errorf("Failed to update max number: %v", err)
	}

	max, err := repo.GetMaxNumber()
	if err != nil {
		t.Errorf("Failed to get max number: %v", err)
	}
	if max != "50" {
		t.Errorf("Expected '50', got %s", max)
	}

	// Test updating existing record
	err = repo.UpdateMaxNumber("100")
	if err != nil {
		t.Errorf("Failed to update max number: %v", err)
	}

	max, err = repo.GetMaxNumber()
	if err != nil {
		t.Errorf("Failed to get max number: %v", err)
	}
	if max != "100" {
		t.Errorf("Expected '100', got %s", max)
	}
}

func TestSetMaxNumberIfGreater(t *testing.T) {
	db := setupTestDBForMaxRequest(t)
	repo := NewMaxRequestRepository(db)

	// Test with empty database - should create
	err := repo.SetMaxNumberIfGreater("100")
	if err != nil {
		t.Errorf("Failed to set max number: %v", err)
	}

	max, err := repo.GetMaxNumber()
	if err != nil {
		t.Errorf("Failed to get max number: %v", err)
	}
	if max != "100" {
		t.Errorf("Expected '100', got %s", max)
	}

	// Test with greater number - should update
	err = repo.SetMaxNumberIfGreater("200")
	if err != nil {
		t.Errorf("Failed to set max number: %v", err)
	}

	max, err = repo.GetMaxNumber()
	if err != nil {
		t.Errorf("Failed to get max number: %v", err)
	}
	if max != "200" {
		t.Errorf("Expected '200', got %s", max)
	}

	// Test with smaller number - should not update
	err = repo.SetMaxNumberIfGreater("150")
	if err != nil {
		t.Errorf("Failed to set max number: %v", err)
	}

	max, err = repo.GetMaxNumber()
	if err != nil {
		t.Errorf("Failed to get max number: %v", err)
	}
	if max != "200" {
		t.Errorf("Expected '200' to remain, got %s", max)
	}

	// Test with same number - should not update
	err = repo.SetMaxNumberIfGreater("200")
	if err != nil {
		t.Errorf("Failed to set max number: %v", err)
	}

	max, err = repo.GetMaxNumber()
	if err != nil {
		t.Errorf("Failed to get max number: %v", err)
	}
	if max != "200" {
		t.Errorf("Expected '200' to remain, got %s", max)
	}
}

func TestSetMaxNumberIfGreaterLargeNumbers(t *testing.T) {
	db := setupTestDBForMaxRequest(t)
	repo := NewMaxRequestRepository(db)

	// Test with large number (string comparison)
	err := repo.SetMaxNumberIfGreater("10000")
	if err != nil {
		t.Errorf("Failed to set max number: %v", err)
	}

	// Test with number that has more digits
	err = repo.SetMaxNumberIfGreater("100000")
	if err != nil {
		t.Errorf("Failed to set max number: %v", err)
	}

	max, err := repo.GetMaxNumber()
	if err != nil {
		t.Errorf("Failed to get max number: %v", err)
	}
	if max != "100000" {
		t.Errorf("Expected '100000', got %s", max)
	}
}
