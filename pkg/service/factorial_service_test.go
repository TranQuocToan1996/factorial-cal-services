package service

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"testing"

	"factorial-cal-services/pkg/domain"

	"gorm.io/gorm"
)

// mockFactorialRepository is a mock implementation of FactorialRepository
type mockFactorialRepository struct {
	calculations map[int64]*domain.FactorialCalculation
	createError  error
	findError    error
	updateError  error
}

func newMockFactorialRepository() *mockFactorialRepository {
	return &mockFactorialRepository{
		calculations: make(map[int64]*domain.FactorialCalculation),
	}
}

func (m *mockFactorialRepository) Create(calc *domain.FactorialCalculation) error {
	if m.createError != nil {
		return m.createError
	}
	m.calculations[calc.Number] = calc
	return nil
}

func (m *mockFactorialRepository) FindByNumber(number int64) (*domain.FactorialCalculation, error) {
	if m.findError != nil {
		return nil, m.findError
	}
	calc, exists := m.calculations[number]
	if !exists {
		return nil, gorm.ErrRecordNotFound
	}
	return calc, nil
}

func (m *mockFactorialRepository) UpdateStatus(number string, status string) error {
	return m.updateError
}

func (m *mockFactorialRepository) UpdateWithCurrentNumber(number int64, s3Key string, checksum string, size int64, status string, bucket string) error {
	if m.updateError != nil {
		return m.updateError
	}
	if calc, exists := m.calculations[number]; exists {
		calc.S3Key = s3Key
		calc.Checksum = checksum
		calc.Size = size
		calc.Status = status
		calc.Bucket = bucket
	}
	return nil
}

// mockCurrentCalculatedRepository is a mock implementation of CurrentCalculatedRepository
type mockCurrentCalculatedRepository struct {
	currentNumber int64
	getError      error
	updateError   error
}

func newMockCurrentCalculatedRepository() *mockCurrentCalculatedRepository {
	return &mockCurrentCalculatedRepository{
		currentNumber: 0,
	}
}

func (m *mockCurrentCalculatedRepository) GetCurrentNumber() (int64, error) {
	if m.getError != nil {
		return 0, m.getError
	}
	return m.currentNumber, nil
}

func (m *mockCurrentCalculatedRepository) UpdateCurrentNumber(curNumber int64) error {
	if m.updateError != nil {
		return m.updateError
	}
	m.currentNumber = curNumber
	return nil
}

// mockMaxRequestRepository is a mock implementation of MaxRequestRepository
type mockMaxRequestRepository struct {
	maxNumber   int64
	getError    error
	updateError error
	setError    error
}

func newMockMaxRequestRepository() *mockMaxRequestRepository {
	return &mockMaxRequestRepository{
		maxNumber: 0,
	}
}

func (m *mockMaxRequestRepository) GetMaxNumber() (int64, error) {
	if m.getError != nil {
		return 0, m.getError
	}
	return m.maxNumber, nil
}

func (m *mockMaxRequestRepository) UpdateMaxNumber(maxNumber int64) error {
	if m.updateError != nil {
		return m.updateError
	}
	m.maxNumber = maxNumber
	return nil
}

func (m *mockMaxRequestRepository) SetMaxNumberIfGreater(maxNumber int64) (int64, error) {
	if m.setError != nil {
		return 0, m.setError
	}
	if maxNumber > m.maxNumber {
		m.maxNumber = maxNumber
		return 1, nil
	}
	return 0, nil
}

// unitTestMockStorageService is a mock implementation of StorageService for unit tests
type unitTestMockStorageService struct {
	storage       map[string]string
	uploadError   error
	downloadError error
}

func newUnitTestMockStorageService() *unitTestMockStorageService {
	return &unitTestMockStorageService{
		storage: make(map[string]string),
	}
}

func (m *unitTestMockStorageService) UploadFactorial(ctx context.Context, number int64, result string) (string, error) {
	if m.uploadError != nil {
		return "", m.uploadError
	}
	key := m.GenerateKey(number)
	m.storage[key] = result
	return key, nil
}

func (m *unitTestMockStorageService) DownloadFactorial(ctx context.Context, s3Key string) (string, error) {
	if m.downloadError != nil {
		return "", m.downloadError
	}
	result, exists := m.storage[s3Key]
	if !exists {
		return "", errors.New("key not found")
	}
	return result, nil
}

func (m *unitTestMockStorageService) GenerateKey(number int64) string {
	return fmt.Sprintf("%d.txt", number)
}

func (m *unitTestMockStorageService) GetBucket() string {
	return "test-bucket"
}

func TestFactorialService_ValidateNumber(t *testing.T) {
	service := NewFactorialService(
		newMockFactorialRepository(),
		newMockCurrentCalculatedRepository(),
		newMockMaxRequestRepository(),
		newUnitTestMockStorageService(),
	).(*factorialService)

	tests := []struct {
		name      string
		number    string
		want      int64
		wantError bool
		errorMsg  string
	}{
		{
			name:      "Valid number zero",
			number:    "0",
			want:      0,
			wantError: false,
		},
		{
			name:      "Valid number one",
			number:    "1",
			want:      1,
			wantError: false,
		},
		{
			name:      "Valid number at max boundary",
			number:    "10000",
			want:      10000,
			wantError: false,
		},
		{
			name:      "Valid number below max",
			number:    "5000",
			want:      5000,
			wantError: false,
		},
		{
			name:      "Negative number",
			number:    "-1",
			want:      0,
			wantError: true,
			errorMsg:  "number must be non-negative",
		},
		{
			name:      "Number exceeds maximum",
			number:    "10001",
			want:      0,
			wantError: true,
			errorMsg:  "number exceeds maximum allowed value of 10000",
		},
		{
			name:      "Invalid format - non-numeric",
			number:    "abc",
			want:      0,
			wantError: true,
			errorMsg:  "invalid number format",
		},
		{
			name:      "Invalid format - empty string",
			number:    "",
			want:      0,
			wantError: true,
			errorMsg:  "invalid number format",
		},
		{
			name:      "Invalid format - decimal",
			number:    "10.5",
			want:      0,
			wantError: true,
			errorMsg:  "invalid number format",
		},
		{
			name:      "Invalid format - with spaces",
			number:    " 10 ",
			want:      0,
			wantError: true,
			errorMsg:  "invalid number format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.ValidateNumber(tt.number)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateNumber() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError && got != tt.want {
				t.Errorf("ValidateNumber() = %v, want %v", got, tt.want)
			}
			if tt.wantError && tt.errorMsg != "" {
				if err == nil || err.Error() == "" {
					t.Errorf("ValidateNumber() expected error message containing '%s', got %v", tt.errorMsg, err)
				}
			}
		})
	}
}

func TestChecksum(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		wantLen  int // SHA256 hex string length is 64
		wantSame bool
	}{
		{
			name:     "Empty string",
			data:     "",
			wantLen:  64,
			wantSame: true,
		},
		{
			name:     "Simple string",
			data:     "hello",
			wantLen:  64,
			wantSame: true,
		},
		{
			name:     "Large number string",
			data:     "3628800",
			wantLen:  64,
			wantSame: true,
		},
		{
			name:     "Very large factorial result",
			data:     "93326215443944152681699238856266700490715968264381621468592963895217599993229915608941463976156518286253697920827223758251185210916864000000000000000000000000",
			wantLen:  64,
			wantSame: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checksum(tt.data)
			if len(got) != tt.wantLen {
				t.Errorf("checksum() length = %v, want %v", len(got), tt.wantLen)
			}
			if tt.wantSame {
				// Same input should produce same checksum
				got2 := checksum(tt.data)
				if got != got2 {
					t.Errorf("checksum() not deterministic: first = %v, second = %v", got, got2)
				}
			}
		})
	}

	// Test that different inputs produce different checksums
	checksum1 := checksum("123")
	checksum2 := checksum("456")
	if checksum1 == checksum2 {
		t.Errorf("checksum() should produce different values for different inputs: both = %v", checksum1)
	}
}

func TestFactorialService_GetPreviousFactorial(t *testing.T) {
	tests := []struct {
		name         string
		number       int64
		setupStorage func(*unitTestMockStorageService)
		want         *big.Int
		wantError    bool
		errorMsg     string
	}{
		{
			name:   "Base case - negative number (factorial of 0)",
			number: -1,
			setupStorage: func(m *unitTestMockStorageService) {
				// No storage needed for base case
			},
			want:      big.NewInt(1),
			wantError: false,
		},
		{
			name:   "Zero (factorial of 0)",
			number: 0,
			setupStorage: func(m *unitTestMockStorageService) {
				key := m.GenerateKey(0)
				m.storage[key] = "1"
			},
			want:      big.NewInt(1),
			wantError: false,
		},
		{
			name:   "Valid previous factorial",
			number: 5,
			setupStorage: func(m *unitTestMockStorageService) {
				key := m.GenerateKey(5)
				m.storage[key] = "120" // factorial of 5
			},
			want:      big.NewInt(120),
			wantError: false,
		},
		{
			name:   "Large factorial",
			number: 10,
			setupStorage: func(m *unitTestMockStorageService) {
				key := m.GenerateKey(10)
				m.storage[key] = "3628800" // factorial of 10
			},
			want:      big.NewInt(3628800),
			wantError: false,
		},
		{
			name:   "Storage download error",
			number: 5,
			setupStorage: func(m *unitTestMockStorageService) {
				m.downloadError = errors.New("storage error")
			},
			want:      nil,
			wantError: true,
			errorMsg:  "failed to download factorial from S3",
		},
		{
			name:   "Key not found in storage",
			number: 5,
			setupStorage: func(m *unitTestMockStorageService) {
				// Don't set the key
			},
			want:      nil,
			wantError: true,
			errorMsg:  "failed to download factorial from S3",
		},
		{
			name:   "Invalid format in storage",
			number: 5,
			setupStorage: func(m *unitTestMockStorageService) {
				key := m.GenerateKey(5)
				m.storage[key] = "not-a-number"
			},
			want:      nil,
			wantError: true,
			errorMsg:  "failed to parse factorial result: invalid format",
		},
		{
			name:   "Empty string in storage",
			number: 5,
			setupStorage: func(m *unitTestMockStorageService) {
				key := m.GenerateKey(5)
				m.storage[key] = ""
			},
			want:      nil,
			wantError: true,
			errorMsg:  "failed to parse factorial result: invalid format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := newUnitTestMockStorageService()
			tt.setupStorage(mockStorage)

			service := NewFactorialService(
				newMockFactorialRepository(),
				newMockCurrentCalculatedRepository(),
				newMockMaxRequestRepository(),
				mockStorage,
			).(*factorialService)

			got, err := service.getPreviousFactorial(tt.number)
			if (err != nil) != tt.wantError {
				t.Errorf("getPreviousFactorial() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError {
				if got == nil {
					t.Errorf("getPreviousFactorial() = nil, want %v", tt.want)
					return
				}
				if got.Cmp(tt.want) != 0 {
					t.Errorf("getPreviousFactorial() = %v, want %v", got, tt.want)
				}
			}
			if tt.wantError && tt.errorMsg != "" {
				if err == nil || err.Error() == "" {
					t.Errorf("getPreviousFactorial() expected error message containing '%s', got %v", tt.errorMsg, err)
				}
			}
		})
	}
}

func TestFactorialService_ContinuelyCalculateFactorial_EdgeCases(t *testing.T) {
	tests := []struct {
		name       string
		current    int64
		max        int64
		setupMocks func(*mockFactorialRepository, *unitTestMockStorageService)
		wantError  bool
		errorMsg   string
	}{
		{
			name:    "Current greater than max - should do nothing",
			current: 10,
			max:     5,
			setupMocks: func(repo *mockFactorialRepository, storage *unitTestMockStorageService) {
				// No setup needed, loop should not execute
			},
			wantError: false,
		},
		{
			name:    "Current equals max - should calculate one number",
			current: 5,
			max:     5,
			setupMocks: func(repo *mockFactorialRepository, storage *unitTestMockStorageService) {
				// Setup previous factorial
				key := storage.GenerateKey(4)
				storage.storage[key] = "24" // factorial of 4
			},
			wantError: false,
		},
		{
			name:    "Repository find error",
			current: 5,
			max:     5,
			setupMocks: func(repo *mockFactorialRepository, storage *unitTestMockStorageService) {
				repo.findError = errors.New("database error")
				key := storage.GenerateKey(4)
				storage.storage[key] = "24"
			},
			wantError: true,
			errorMsg:  "failed to query factorial",
		},
		{
			name:    "Storage upload error",
			current: 5,
			max:     5,
			setupMocks: func(repo *mockFactorialRepository, storage *unitTestMockStorageService) {
				storage.uploadError = errors.New("upload error")
				key := storage.GenerateKey(4)
				storage.storage[key] = "24"
			},
			wantError: true,
			errorMsg:  "failed to upload factorial to S3",
		},
		{
			name:    "Repository update error",
			current: 5,
			max:     5,
			setupMocks: func(repo *mockFactorialRepository, storage *unitTestMockStorageService) {
				repo.updateError = errors.New("update error")
				key := storage.GenerateKey(4)
				storage.storage[key] = "24"
			},
			wantError: true,
			errorMsg:  "failed to update factorial record",
		},
		{
			name:    "Skip already done factorial",
			current: 5,
			max:     5,
			setupMocks: func(repo *mockFactorialRepository, storage *unitTestMockStorageService) {
				// Mark as already done
				repo.calculations[5] = &domain.FactorialCalculation{
					Number: 5,
					Status: domain.StatusDone,
				}
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockFactorialRepository()
			mockStorage := newUnitTestMockStorageService()
			tt.setupMocks(mockRepo, mockStorage)

			service := NewFactorialService(
				mockRepo,
				newMockCurrentCalculatedRepository(),
				newMockMaxRequestRepository(),
				mockStorage,
			).(*factorialService)

			err := service.continuelyCalculateFactorial(tt.current, tt.max, nil)
			if (err != nil) != tt.wantError {
				t.Errorf("continuelyCalculateFactorial() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if tt.wantError && tt.errorMsg != "" {
				if err == nil || err.Error() == "" {
					t.Errorf("continuelyCalculateFactorial() expected error message containing '%s', got %v", tt.errorMsg, err)
				}
			}
		})
	}
}
