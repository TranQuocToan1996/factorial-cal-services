package repository

import (
	"testing"

	"factorial-cal-services/pkg/domain"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto migrate
	err = db.AutoMigrate(&domain.FactorialCalculation{})
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	return db
}

func TestCreate(t *testing.T) {
	db := setupTestDB(t)
	repo := NewFactorialRepository(db)

	calc := &domain.FactorialCalculation{
		Number: "10",
		Status: domain.StatusCalculating,
		S3Key:  "",
	}

	err := repo.Create(calc)
	if err != nil {
		t.Errorf("Failed to create record: %v", err)
	}

	if calc.ID == 0 {
		t.Error("Expected ID to be set after creation")
	}
}

func TestFindByNumber(t *testing.T) {
	db := setupTestDB(t)
	repo := NewFactorialRepository(db)

	// Create a record
	calc := &domain.FactorialCalculation{
		Number: "10",
		Status: domain.StatusDone,
		S3Key:  "factorials/10.txt",
	}

	err := repo.Create(calc)
	if err != nil {
		t.Fatalf("Failed to create record: %v", err)
	}

	// Find by number
	found, err := repo.FindByNumber("10")
	if err != nil {
		t.Errorf("Failed to find record: %v", err)
	}

	if found == nil {
		t.Error("Expected to find record, got nil")
		return
	}

	if found.Number != "10" {
		t.Errorf("Expected number 10, got %s", found.Number)
	}

	if found.Status != domain.StatusDone {
		t.Errorf("Expected status done, got %s", found.Status)
	}
}

func TestFindByNumberNotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewFactorialRepository(db)

	found, err := repo.FindByNumber("999")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if found != nil {
		t.Error("Expected nil for non-existent record")
	}
}

func TestUpdateStatus(t *testing.T) {
	db := setupTestDB(t)
	repo := NewFactorialRepository(db)

	// Create a record
	calc := &domain.FactorialCalculation{
		Number: "10",
		Status: domain.StatusCalculating,
		S3Key:  "",
	}

	err := repo.Create(calc)
	if err != nil {
		t.Fatalf("Failed to create record: %v", err)
	}

	// Update status
	err = repo.UpdateStatus("10", domain.StatusDone)
	if err != nil {
		t.Errorf("Failed to update status: %v", err)
	}

	// Verify update
	found, err := repo.FindByNumber("10")
	if err != nil {
		t.Errorf("Failed to find record: %v", err)
	}

	if found.Status != domain.StatusDone {
		t.Errorf("Expected status done, got %s", found.Status)
	}
}

func TestUpdateS3Key(t *testing.T) {
	db := setupTestDB(t)
	repo := NewFactorialRepository(db)

	// Create a record
	calc := &domain.FactorialCalculation{
		Number: "10",
		Status: domain.StatusUploading,
		S3Key:  "",
	}

	err := repo.Create(calc)
	if err != nil {
		t.Fatalf("Failed to create record: %v", err)
	}

	// Update S3 key and status
	s3Key := "factorials/10.txt"
	err = repo.UpdateS3Key("10", s3Key, domain.StatusDone)
	if err != nil {
		t.Errorf("Failed to update S3 key: %v", err)
	}

	// Verify update
	found, err := repo.FindByNumber("10")
	if err != nil {
		t.Errorf("Failed to find record: %v", err)
	}

	if found.S3Key != s3Key {
		t.Errorf("Expected S3 key %s, got %s", s3Key, found.S3Key)
	}

	if found.Status != domain.StatusDone {
		t.Errorf("Expected status done, got %s", found.Status)
	}
}

