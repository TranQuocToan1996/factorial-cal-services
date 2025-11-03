package service

import (
	"context"
	"fmt"
	"math/big"
	"strconv"

	"factorial-cal-services/pkg/domain"
)

// IncrementalFactorialService handles incremental factorial calculations
type IncrementalFactorialService interface {
	CalculateIncremental(ctx context.Context, curNumber string, curFactorial string, targetNumber string) (string, error)
	CalculateFromCurrent(ctx context.Context, currentRepo domain.CurrentCalculatedRepository, redisService RedisService, storageService StorageService, targetNumber string) (string, error)
}

type incrementalFactorialService struct {
	factorialService FactorialService
}

// NewIncrementalFactorialService creates a new incremental factorial service
func NewIncrementalFactorialService(factorialService FactorialService) IncrementalFactorialService {
	return &incrementalFactorialService{
		factorialService: factorialService,
	}
}

// CalculateIncremental calculates factorial incrementally from current to target
// Logic: next_fac = (cur_number + 1) * cur_fac
func (s *incrementalFactorialService) CalculateIncremental(ctx context.Context, curNumber string, curFactorial string, targetNumber string) (string, error) {
	curNum, err := strconv.ParseInt(curNumber, 10, 64)
	if err != nil {
		return "", fmt.Errorf("invalid current number: %w", err)
	}

	targetNum, err := strconv.ParseInt(targetNumber, 10, 64)
	if err != nil {
		return "", fmt.Errorf("invalid target number: %w", err)
	}

	if curNum >= targetNum {
		return curFactorial, nil // Already calculated
	}

	// Parse current factorial as big.Int
	result := new(big.Int)
	result, ok := result.SetString(curFactorial, 10)
	if !ok {
		return "", fmt.Errorf("invalid current factorial format")
	}

	// Calculate incrementally: next_fac = (cur_number + 1) * cur_fac
	for i := curNum + 1; i <= targetNum; i++ {
		multiplier := big.NewInt(i)
		result.Mul(result, multiplier)
	}

	return result.String(), nil
}

// CalculateFromCurrent gets current number and factorial from repository/storage and calculates incrementally
func (s *incrementalFactorialService) CalculateFromCurrent(ctx context.Context, currentRepo domain.CurrentCalculatedRepository, redisService RedisService, storageService StorageService, targetNumber string) (string, error) {
	// Get current calculated number from DB
	curNumber, err := currentRepo.GetCurrentNumber()
	if err != nil {
		return "", fmt.Errorf("failed to get current number: %w", err)
	}

	if curNumber == "0" || curNumber == "" {
		// First run - calculate from scratch
		return s.factorialService.CalculateFactorial(targetNumber)
	}

	// Get current factorial from Redis or storage
	var curFactorial string
	if redisService.ShouldCache(curNumber) {
		curFactorial, _ = redisService.Get(ctx, curNumber)
	}

	if curFactorial == "" {
		// Try to get from storage
		s3Key := storageService.GenerateS3Key(curNumber)
		curFactorial, _ = storageService.Download(ctx, s3Key)
	}

	if curFactorial == "" {
		// Fallback: calculate from scratch
		return s.factorialService.CalculateFactorial(targetNumber)
	}

	// Calculate incrementally
	return s.CalculateIncremental(ctx, curNumber, curFactorial, targetNumber)
}

