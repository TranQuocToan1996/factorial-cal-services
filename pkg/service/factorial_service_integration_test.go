package service

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"factorial-cal-services/pkg/domain"
	"factorial-cal-services/pkg/repository"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// mockStorageService implements StorageService for testing
type mockStorageService struct {
	storage map[string]string // key -> value
}

func newMockStorageService() *mockStorageService {
	return &mockStorageService{
		storage: make(map[string]string),
	}
}

func (m *mockStorageService) UploadFactorial(ctx context.Context, number int64, result string) (string, error) {
	key := m.GenerateKey(number)
	m.storage[key] = result
	return key, nil
}

func (m *mockStorageService) DownloadFactorial(ctx context.Context, s3Key string) (string, error) {
	value, ok := m.storage[s3Key]
	if !ok {
		return "", fmt.Errorf("key not found: %s", s3Key)
	}
	return value, nil
}

func (m *mockStorageService) GenerateKey(number int64) string {
	return fmt.Sprintf("%d.txt", number)
}

// setupTestDB creates an in-memory SQLite database with schema
func setupTestDB(t *testing.T) *gorm.DB {
	// Use in-memory SQLite database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	// Auto-migrate all tables
	err = db.AutoMigrate(
		&domain.FactorialCalculation{},
		&domain.FactorialMaxRequestNumber{},
		&domain.FactorialCurrentCalculatedNumber{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

// setupBaseCases inserts base cases (0, 1, 2, 3) into database and storage
func setupBaseCases(t *testing.T, db *gorm.DB, storage StorageService, ctx context.Context) {
	baseCases := []struct {
		number int64
		value  string
	}{
		{0, "1"},
		{1, "1"},
		{2, "2"},
		{3, "6"},
	}

	for _, bc := range baseCases {
		// Insert into database
		calc := &domain.FactorialCalculation{
			Number: bc.number,
			Status: domain.StatusDone,
			S3Key:  fmt.Sprintf("%d.txt", bc.number),
			Size:   int64(len(bc.value)),
		}
		if err := db.Create(calc).Error; err != nil {
			t.Fatalf("Failed to create base case for %d: %v", bc.number, err)
		}

		// Store in mock storage
		_, err := storage.UploadFactorial(ctx, bc.number, bc.value)
		if err != nil {
			t.Fatalf("Failed to upload base case for %d: %v", bc.number, err)
		}
	}

	// Set current calculated number to 4
	curCalc := &domain.FactorialCurrentCalculatedNumber{
		CurNumber: 4,
	}
	if err := db.Create(curCalc).Error; err != nil {
		t.Fatalf("Failed to create current calculated number: %v", err)
	}
}

// calculateFactorial calculates factorial using big.Int for verification
func calculateFactorial(n int64) *big.Int {
	result := big.NewInt(1)
	for i := int64(2); i <= n; i++ {
		result.Mul(result, big.NewInt(i))
	}
	return result
}

func TestFactorialService_Integration_UpTo10(t *testing.T) {
	// Setup test database
	db := setupTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}()

	// Setup mock storage (use mockStorageService for consistent key format)
	mockStorage := newMockStorageService()
	ctx := context.Background()

	// Setup base cases (0, 1, 2, 3)
	setupBaseCases(t, db, mockStorage, ctx)

	// Initialize repositories
	factorialRepo := repository.NewFactorialRepository(db)
	currentCalculatedRepo := repository.NewCurrentCalculatedRepository(db)
	maxRequestRepo := repository.NewMaxRequestRepository(db)

	// Initialize service
	factorialService := NewFactorialService(
		factorialRepo,
		currentCalculatedRepo,
		maxRequestRepo,
		mockStorage,
	).(*factorialService)

	// Set max number to 10
	err := maxRequestRepo.UpdateMaxNumber(10)
	if err != nil {
		t.Fatalf("Failed to set max number: %v", err)
	}

	// Get current number (should be 4)
	current, err := currentCalculatedRepo.GetCurrentNumber()
	if err != nil {
		t.Fatalf("Failed to get current number: %v", err)
	}
	if current != 4 {
		t.Fatalf("Expected current number to be 4, got %d", current)
	}

	// Get max number (should be 10)
	max, err := maxRequestRepo.GetMaxNumber()
	if err != nil {
		t.Fatalf("Failed to get max number: %v", err)
	}
	if max != 10 {
		t.Fatalf("Expected max number to be 10, got %d", max)
	}

	// Calculate factorials from 4 to 10
	err = factorialService.continuelyCalculateFactorial(current, max, nil)
	if err != nil {
		t.Fatalf("Failed to calculate factorials: %v", err)
	}

	// Verify all factorials are calculated correctly (0-10)
	expectedResults := make(map[int64]*big.Int)
	for i := int64(0); i <= 10; i++ {
		expectedResults[i] = calculateFactorial(i)
	}

	// Verify database records
	for number := int64(0); number <= 10; number++ {
		t.Run(fmt.Sprintf("Verify_factorial_%d", number), func(t *testing.T) {
			// Check database record
			calc, err := factorialRepo.FindByNumber(number)
			if err != nil {
				t.Fatalf("Failed to find factorial %d: %v", number, err)
			}

			if calc.Status != domain.StatusDone {
				t.Errorf("Expected status 'done' for %d, got '%s'", number, calc.Status)
			}

			// Verify S3 key format (mockStorageService uses "number.txt" format)
			expectedS3Key := fmt.Sprintf("%d.txt", number)
			if calc.S3Key != expectedS3Key {
				t.Errorf("Expected S3 key '%s' for %d, got '%s'", expectedS3Key, number, calc.S3Key)
			}

			// Verify storage content
			result, err := mockStorage.DownloadFactorial(ctx, calc.S3Key)
			if err != nil {
				t.Fatalf("Failed to download factorial %d from storage: %v", number, err)
			}

			// Parse and verify value
			actualValue, ok := new(big.Int).SetString(result, 10)
			if !ok {
				t.Fatalf("Failed to parse factorial result for %d: %s", number, result)
			}

			expectedValue := expectedResults[number]
			if actualValue.Cmp(expectedValue) != 0 {
				t.Errorf("Factorial %d: expected %s, got %s", number, expectedValue.String(), actualValue.String())
			}

			// Verify size
			expectedSize := int64(len(result))
			if calc.Size != expectedSize {
				t.Errorf("Expected size %d for %d, got %d", expectedSize, number, calc.Size)
			}
		})
	}

	// Verify current calculated number is updated to 11
	updatedCurrent, err := currentCalculatedRepo.GetCurrentNumber()
	if err != nil {
		t.Fatalf("Failed to get updated current number: %v", err)
	}
	if updatedCurrent != 11 {
		t.Errorf("Expected current calculated number to be 11, got %d", updatedCurrent)
	}

	// Verify all factorials are in storage (using mockStorageService)
	for number := int64(0); number <= 10; number++ {
		key := fmt.Sprintf("%d.txt", number)
		value, exists := mockStorage.storage[key]
		if !exists {
			t.Errorf("Factorial %d not found in storage (key: %s)", number, key)
			continue
		}

		expectedValue := expectedResults[number]
		actualValue, ok := new(big.Int).SetString(value, 10)
		if !ok {
			t.Errorf("Failed to parse stored value for %d: %s", number, value)
			continue
		}

		if actualValue.Cmp(expectedValue) != 0 {
			t.Errorf("Storage mismatch for %d: expected %s, got %s", number, expectedValue.String(), actualValue.String())
		}
	}
}

func TestFactorialService_Integration_SequentialCalculation(t *testing.T) {
	// Setup test database
	db := setupTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}()

	// Setup mock storage
	mockStorage := newMockStorageService()
	ctx := context.Background()

	// Setup base cases (0, 1, 2, 3)
	setupBaseCases(t, db, mockStorage, ctx)

	// Initialize repositories
	factorialRepo := repository.NewFactorialRepository(db)
	currentCalculatedRepo := repository.NewCurrentCalculatedRepository(db)
	maxRequestRepo := repository.NewMaxRequestRepository(db)

	// Initialize service
	factorialService := NewFactorialService(
		factorialRepo,
		currentCalculatedRepo,
		maxRequestRepo,
		mockStorage,
	).(*factorialService)

	// Set max number to 10
	err := maxRequestRepo.UpdateMaxNumber(10)
	if err != nil {
		t.Fatalf("Failed to set max number: %v", err)
	}

	// Calculate factorials sequentially: 4, 5, 6, 7, 8, 9, 10
	current, _ := currentCalculatedRepo.GetCurrentNumber()
	max, _ := maxRequestRepo.GetMaxNumber()
	err = factorialService.continuelyCalculateFactorial(current, max, nil)
	if err != nil {
		t.Fatalf("Failed to calculate factorials: %v", err)
	}

	// Verify sequential calculation correctness: each factorial = previous * current
	for number := int64(4); number <= 10; number++ {
		t.Run(fmt.Sprintf("Verify_sequential_factorial_%d", number), func(t *testing.T) {
			// Get current factorial
			calc, err := factorialRepo.FindByNumber(number)
			if err != nil {
				t.Fatalf("Failed to find factorial %d: %v", number, err)
			}

			result, err := mockStorage.DownloadFactorial(ctx, calc.S3Key)
			if err != nil {
				t.Fatalf("Failed to download factorial %d: %v", number, err)
			}

			currentFactorial, _ := new(big.Int).SetString(result, 10)

			// Get previous factorial
			prevCalc, err := factorialRepo.FindByNumber(number - 1)
			if err != nil {
				t.Fatalf("Failed to find previous factorial %d: %v", number-1, err)
			}

			prevResult, err := mockStorage.DownloadFactorial(ctx, prevCalc.S3Key)
			if err != nil {
				t.Fatalf("Failed to download previous factorial %d: %v", number-1, err)
			}

			previousFactorial, _ := new(big.Int).SetString(prevResult, 10)

			// Verify: factorial(n) = factorial(n-1) * n
			expectedFactorial := new(big.Int).Mul(previousFactorial, big.NewInt(number))
			if currentFactorial.Cmp(expectedFactorial) != 0 {
				t.Errorf("Factorial %d: expected %s (from %s * %d), got %s",
					number, expectedFactorial.String(), previousFactorial.String(), number, currentFactorial.String())
			}

			// Also verify against calculated value
			expectedValue := calculateFactorial(number)
			if currentFactorial.Cmp(expectedValue) != 0 {
				t.Errorf("Factorial %d: expected %s, got %s", number, expectedValue.String(), currentFactorial.String())
			}
		})
	}
}

func TestFactorialService_Integration_WithRedis(t *testing.T) {
	// Setup test database
	db := setupTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}()

	// Setup Redis mock
	mr := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	// Setup mock storage
	storage := newMockStorageService()
	ctx := context.Background()

	// Setup Redis service
	redisService := NewRedisService(redisClient, 24*time.Hour, 1000)

	// Setup base cases
	setupBaseCases(t, db, storage, ctx)

	// Initialize repositories
	factorialRepo := repository.NewFactorialRepository(db)
	currentCalculatedRepo := repository.NewCurrentCalculatedRepository(db)
	maxRequestRepo := repository.NewMaxRequestRepository(db)

	// Initialize service
	factorialService := NewFactorialService(
		factorialRepo,
		currentCalculatedRepo,
		maxRequestRepo,
		storage,
	).(*factorialService)

	// Set max number to 10
	err := maxRequestRepo.UpdateMaxNumber(10)
	if err != nil {
		t.Fatalf("Failed to set max number: %v", err)
	}

	// Calculate factorials
	current, _ := currentCalculatedRepo.GetCurrentNumber()
	max, _ := maxRequestRepo.GetMaxNumber()
	err = factorialService.continuelyCalculateFactorial(current, max, nil)
	if err != nil {
		t.Fatalf("Failed to calculate factorials: %v", err)
	}

	// Test Redis caching for numbers < 1000 (threshold)
	for number := int64(0); number <= 10; number++ {
		numberStr := fmt.Sprintf("%d", number)
		if redisService.ShouldCache(numberStr) {
			// Get result from storage
			calc, _ := factorialRepo.FindByNumber(number)
			result, _ := storage.DownloadFactorial(ctx, calc.S3Key)

			// Set in Redis
			err := redisService.Set(ctx, numberStr, result)
			if err != nil {
				t.Errorf("Failed to set Redis cache for %d: %v", number, err)
				continue
			}

			// Get from Redis
			cached, err := redisService.Get(ctx, numberStr)
			if err != nil {
				t.Errorf("Failed to get from Redis for %d: %v", number, err)
				continue
			}

			if cached != result {
				t.Errorf("Redis cache mismatch for %d: expected %s, got %s", number, result, cached)
			}
		}
	}
}

func TestFactorialService_Integration_EdgeCases(t *testing.T) {
	// Setup test database
	db := setupTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}()

	storage := newMockStorageService()
	ctx := context.Background()

	// Setup base cases
	setupBaseCases(t, db, storage, ctx)

	// Initialize repositories
	factorialRepo := repository.NewFactorialRepository(db)
	currentCalculatedRepo := repository.NewCurrentCalculatedRepository(db)
	maxRequestRepo := repository.NewMaxRequestRepository(db)

	// Initialize service
	factorialService := NewFactorialService(
		factorialRepo,
		currentCalculatedRepo,
		maxRequestRepo,
		storage,
	).(*factorialService)

	// Test: Calculate when current > max (should do nothing)
	current, _ := currentCalculatedRepo.GetCurrentNumber() // Should be 4
	err := maxRequestRepo.UpdateMaxNumber(2)               // Set max to 2 (less than current)
	if err != nil {
		t.Fatalf("Failed to set max number: %v", err)
	}

	// This should not calculate anything
	err = factorialService.continuelyCalculateFactorial(current, 2, nil)
	if err != nil {
		t.Fatalf("Should not error when current > max: %v", err)
	}

	// Test: Calculate already done factorial (should skip)
	// Set max back to 10
	err = maxRequestRepo.UpdateMaxNumber(10)
	if err != nil {
		t.Fatalf("Failed to set max number: %v", err)
	}

	// Calculate 4 again (already done from previous test)
	err = factorialService.continuelyCalculateFactorial(4, 4, nil)
	if err != nil {
		t.Fatalf("Should skip already done factorial: %v", err)
	}
}

// Benchmark factorial calculation
func BenchmarkFactorialService_CalculateUpTo10(b *testing.B) {
	db := setupTestDB(&testing.T{})
	defer func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}()

	storage := newMockStorageService()
	ctx := context.Background()

	setupBaseCases(&testing.T{}, db, storage, ctx)

	factorialRepo := repository.NewFactorialRepository(db)
	currentCalculatedRepo := repository.NewCurrentCalculatedRepository(db)
	maxRequestRepo := repository.NewMaxRequestRepository(db)

	factorialService := NewFactorialService(
		factorialRepo,
		currentCalculatedRepo,
		maxRequestRepo,
		storage,
	).(*factorialService)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Reset database state
		db.Exec("DELETE FROM factorial_calculations WHERE number > 3")
		db.Exec("UPDATE factorial_current_calculated_numbers SET cur_number = 4")

		maxRequestRepo.UpdateMaxNumber(10)
		current, _ := currentCalculatedRepo.GetCurrentNumber()
		max, _ := maxRequestRepo.GetMaxNumber()
		factorialService.continuelyCalculateFactorial(current, max, nil)
	}
}
