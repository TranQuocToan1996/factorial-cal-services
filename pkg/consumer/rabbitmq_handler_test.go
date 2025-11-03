package consumer

import (
	"testing"
)

func TestFindMaxNumber(t *testing.T) {
	// Create a handler instance (we only need to test the findMaxNumber method)
	handler := &FactorialMessageHandler{}

	tests := []struct {
		name     string
		numbers  []string
		expected string
	}{
		{
			name:     "Empty slice",
			numbers:  []string{},
			expected: "0",
		},
		{
			name:     "Single number",
			numbers:  []string{"10"},
			expected: "10",
		},
		{
			name:     "Multiple numbers - ascending",
			numbers:  []string{"5", "10", "15", "20"},
			expected: "20",
		},
		{
			name:     "Multiple numbers - descending",
			numbers:  []string{"20", "15", "10", "5"},
			expected: "20",
		},
		{
			name:     "Multiple numbers - unsorted",
			numbers:  []string{"15", "5", "25", "10", "20"},
			expected: "25",
		},
		{
			name:     "Large numbers",
			numbers:  []string{"100", "1000", "500", "10000"},
			expected: "10000",
		},
		{
			name:     "Numbers with different digit lengths",
			numbers:  []string{"9", "10", "100", "1000"},
			expected: "1000",
		},
		{
			name:     "Same numbers",
			numbers:  []string{"10", "10", "10"},
			expected: "10",
		},
		{
			name:     "Zero included",
			numbers:  []string{"0", "5", "10"},
			expected: "10",
		},
		{
			name:     "Large string numbers",
			numbers:  []string{"10000", "20000", "50000"},
			expected: "50000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.findMaxNumber(tt.numbers)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestFindMaxNumberInvalidStrings(t *testing.T) {
	handler := &FactorialMessageHandler{}

	tests := []struct {
		name     string
		numbers  []string
		expected string // Should use lexicographic comparison for invalid strings
	}{
		{
			name:     "Mix of valid and invalid strings",
			numbers:  []string{"10", "abc", "20"},
			expected: "20", // Invalid strings are skipped, valid numbers compared
		},
		{
			name:     "All invalid strings - fallback",
			numbers:  []string{"abc", "xyz", "def"},
			expected: "abc", // Returns first string when no valid numbers found
		},
		{
			name:     "Invalid string longer",
			numbers:  []string{"10", "verylongstring"},
			expected: "10", // Invalid strings are skipped, valid number returned
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.findMaxNumber(tt.numbers)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestFindMaxNumberEdgeCases(t *testing.T) {
	handler := &FactorialMessageHandler{}

	tests := []struct {
		name     string
		numbers  []string
		expected string
	}{
		{
			name:     "Very large numbers",
			numbers:  []string{"999999999", "1000000000", "999999998"},
			expected: "1000000000",
		},
		{
			name:     "Numbers with leading zeros",
			numbers:  []string{"0010", "010", "10", "00010"},
			expected: "0010", // Returns original string, all parse to same value (10)
		},
		{
			name:     "Mixed format",
			numbers:  []string{"5", "10", "15", "100", "1000"},
			expected: "1000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.findMaxNumber(tt.numbers)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}
