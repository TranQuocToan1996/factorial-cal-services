package service

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
)

// FactorialService handles factorial calculations
type FactorialService interface {
	CalculateFactorial(number string) (string, error)
	ValidateNumber(number string) (int64, error)
}

type factorialService struct {
	maxFactorial int64
}

// NewFactorialService creates a new factorial service
func NewFactorialService() FactorialService {
	return &factorialService{
		maxFactorial: 10000, // Default, should be overridden via NewFactorialServiceWithLimit
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

// CalculateFactorial calculates the factorial of a number given as string
func (s *factorialService) CalculateFactorial(number string) (string, error) {
	// Validate input
	n, err := s.ValidateNumber(number)
	if err != nil {
		return "", err
	}

	// Calculate factorial using big.Int
	result := big.NewInt(1)

	for i := int64(2); i <= n; i++ {
		result.Mul(result, big.NewInt(i))
	}

	// Return result as string
	return result.String(), nil
}
