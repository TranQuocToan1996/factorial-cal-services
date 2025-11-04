package service

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"sync"

	"factorial-cal-services/pkg/domain"
	"factorial-cal-services/pkg/repository"

	"gorm.io/gorm"
)

// FactorialService handles factorial calculations
type FactorialService interface {
	ValidateNumber(number string) (int64, error)
}

type factorialService struct {
	start bool
	lock  sync.Mutex

	maxFactorial                int64
	repository                  repository.FactorialRepository
	currentCalculatedRepository repository.CurrentCalculatedRepository
	maxRequestRepository        repository.MaxRequestRepository
	redisService                RedisService
	s3Service                   S3Service
}

// NewFactorialService creates a new factorial service
func NewFactorialService(
	repository repository.FactorialRepository,
	currentCalculatedRepository repository.CurrentCalculatedRepository,
	maxRequestRepository repository.MaxRequestRepository,
	redisService RedisService,
	s3Service S3Service,
) FactorialService {
	return &factorialService{
		maxFactorial:                10000, // Default, should be overridden via NewFactorialServiceWithLimit
		repository:                  repository,
		currentCalculatedRepository: currentCalculatedRepository,
		maxRequestRepository:        maxRequestRepository,
		redisService:                redisService,
		s3Service:                   s3Service,
	}
}

// NewFactorialServiceWithLimit creates a new factorial service with custom max limit
func NewFactorialServiceWithLimit(maxFactorial int64) FactorialService {
	if maxFactorial <= 0 {
		maxFactorial = 10000 // Default
	}
	return &factorialService{
		maxFactorial: maxFactorial,
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

func (s *factorialService) backgroundContinuelyCalculateFactorial(privFactorial *big.Int) error {
	max, err := s.maxRequestRepository.GetMaxNumber()
	if err != nil {
		return fmt.Errorf("failed to get max number: %w", err)
	}
	// get current
	current, err := s.currentCalculatedRepository.GetCurrentNumber()
	if err != nil {
		return fmt.Errorf("failed to get current number: %w", err)
	}

	if current > max {
		return nil
	}

	factorial, err := s.repository.FindByNumber(current)
	if factorial != nil && factorial.Status == domain.StatusDone {
		return nil
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

	// TODO;
	// if privFactorial == nil {
	// 	privFactorial, err = s.getPreviousFactorial(current)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to get previous factorial: %w", err)
	// 	}
	// }

	// for i := current; i <= max; i++ {
	// 	nextFactorial := new(big.Int).Mul(privFactorial, big.NewInt(i))
	// }

	// Save status calculating
	// calculate current factorial
	// save to S3
	// Update state calculated and next number

	return nil
}

func (s *factorialService) previousFactorialKey(number string) string {
	n, _ := strconv.ParseInt(number, 10, 64)
	return fmt.Sprintf("%d.txt", n-1)
}

func (s *factorialService) getPreviousFactorial(number string) (*big.Int, error) {
	key := s.previousFactorialKey(number)
	result, err := s.s3Service.DownloadFactorial(context.Background(), key)
	if err != nil {
		return nil, fmt.Errorf("failed to download factorial from S3: %w", err)
	}
	currentFactorial, _ := new(big.Int).SetString(result, 10)
	return currentFactorial, nil
}
