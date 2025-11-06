package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"time"

	"factorial-cal-services/pkg/domain"
	"factorial-cal-services/pkg/repository"

	"gorm.io/gorm"
)

// FactorialService handles factorial calculations
type FactorialService interface {
	ValidateNumber(number string) (int64, error)
	StartContinuelyCalculateFactorial()
}

type factorialService struct {
	maxFactorial                int64
	repository                  repository.FactorialRepository
	currentCalculatedRepository repository.CurrentCalculatedRepository
	maxRequestRepository        repository.MaxRequestRepository
	storage                     StorageService
}

// NewFactorialService creates a new factorial service
func NewFactorialService(
	repository repository.FactorialRepository,
	currentCalculatedRepository repository.CurrentCalculatedRepository,
	maxRequestRepository repository.MaxRequestRepository,
	storage StorageService,
) FactorialService {
	return &factorialService{
		repository:                  repository,
		currentCalculatedRepository: currentCalculatedRepository,
		maxRequestRepository:        maxRequestRepository,
		storage:                     storage,
		maxFactorial:                100000,
	}
}

// ValidateNumber validates and parses the input number string
func (s *factorialService) ValidateNumber(number string) (int64, error) {
	// Parse string to int64
	n, err := strconv.ParseInt(number, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number format: %w", err)
	}

	// Check bounds
	if n < 0 {
		return 0, errors.New("number must be non-negative")
	}

	if n > s.maxFactorial {
		return 0, fmt.Errorf("number exceeds maximum allowed value of %d", s.maxFactorial)
	}

	return n, nil
}

// TODO: Add ctx, graceful shutdown
func (s *factorialService) StartContinuelyCalculateFactorial() {
	go func() {
		for {
			time.Sleep(1 * time.Second)
			current, err := s.currentCalculatedRepository.GetCurrentNumber()
			if err != nil {
				log.Printf("failed to get current number: %v", err)
				continue
			}
			max, err := s.maxRequestRepository.GetMaxNumber()
			if err != nil {
				log.Printf("failed to get max number: %v", err)
				continue
			}
			if current > s.maxFactorial {
				log.Printf("current number exceeds maximum allowed value of %d", s.maxFactorial)
				continue
			}
			err = s.continuelyCalculateFactorial(current, max, nil)
			if err != nil {
				log.Printf("failed to calculate factorial: %v", err)
			}
		}
	}()
}

func (s *factorialService) continuelyCalculateFactorial(current, max int64, factorialBigInt *big.Int) error {
	for ; current <= max; current++ {
		// check status process
		factorial, err := s.repository.FindByNumber(current)
		if factorial != nil && factorial.Status == domain.StatusDone {
			continue // Skip already done factorial, continue to next
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("failed to query factorial %d: %w", current, err)
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err := s.repository.Create(&domain.FactorialCalculation{
				Number: current,
				Status: domain.StatusCalculating,
			})
			if err != nil {
				return fmt.Errorf("failed to create factorial record: %w", err)
			}
		}

		if factorialBigInt == nil {
			factorialBigInt, err = s.getPreviousFactorial(current - 1)
			if err != nil {
				return fmt.Errorf("failed to get previous factorial: %w", err)
			}
		}

		// Calculate and save
		factorialBigInt = new(big.Int).Mul(factorialBigInt, big.NewInt(current))
		factorialStr := factorialBigInt.String()
		s3Key, err := s.storage.UploadFactorial(context.Background(), current, factorialStr)
		if err != nil {
			return fmt.Errorf("failed to upload factorial to S3: %w", err)
		}

		// Size should be the string length (bytes), not bit length
		err = s.repository.UpdateWithCurrentNumber(current, s3Key, checksum(factorialStr), int64(len(factorialStr)), domain.StatusDone)
		if err != nil {
			return fmt.Errorf("failed to update factorial record: %w", err)
		}
	}
	return nil
}

func (s *factorialService) getPreviousFactorial(number int64) (*big.Int, error) {
	// number is already (current - 1), so we use it directly
	// For example: if current=4, we call getPreviousFactorial(3), so we want factorial(3)
	if number < 0 {
		// Factorial of 0 is 1, no previous needed
		return big.NewInt(1), nil
	}
	key := s.storage.GenerateKey(number)
	result, err := s.storage.DownloadFactorial(context.Background(), key)
	if err != nil {
		return nil, fmt.Errorf("failed to download factorial from S3: %w", err)
	}
	currentFactorial, ok := new(big.Int).SetString(result, 10)
	if !ok || currentFactorial == nil {
		return nil, fmt.Errorf("failed to parse factorial result: invalid format")
	}
	return currentFactorial, nil
}

func checksum(data string) string {
	h := sha256.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
