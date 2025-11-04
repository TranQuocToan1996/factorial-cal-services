package service

// func TestCalculateFactorial(t *testing.T) {
// 	service := NewFactorialService()

// 	tests := []struct {
// 		name        string
// 		number      string
// 		expected    string
// 		expectError bool
// 	}{
// 		{
// 			name:        "Factorial of 0",
// 			number:      "0",
// 			expected:    "1",
// 			expectError: false,
// 		},
// 		{
// 			name:        "Factorial of 1",
// 			number:      "1",
// 			expected:    "1",
// 			expectError: false,
// 		},
// 		{
// 			name:        "Factorial of 5",
// 			number:      "5",
// 			expected:    "120",
// 			expectError: false,
// 		},
// 		{
// 			name:        "Factorial of 10",
// 			number:      "10",
// 			expected:    "3628800",
// 			expectError: false,
// 		},
// 		{
// 			name:        "Factorial of 20",
// 			number:      "20",
// 			expected:    "2432902008176640000",
// 			expectError: false,
// 		},
// 		{
// 			name:        "Negative number",
// 			number:      "-1",
// 			expected:    "",
// 			expectError: true,
// 		},
// 		{
// 			name:        "Number exceeds max",
// 			number:      "20001",
// 			expected:    "",
// 			expectError: true,
// 		},
// 		{
// 			name:        "Invalid format",
// 			number:      "abc",
// 			expected:    "",
// 			expectError: true,
// 		},
// 		{
// 			name:        "Empty string",
// 			number:      "",
// 			expected:    "",
// 			expectError: true,
// 		},
// 		{
// 			name:        "Float number",
// 			number:      "5.5",
// 			expected:    "",
// 			expectError: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			result, err := service.CalculateFactorial(tt.number)

// 			if tt.expectError {
// 				if err == nil {
// 					t.Errorf("Expected error but got none")
// 				}
// 				return
// 			}

// 			if err != nil {
// 				t.Errorf("Unexpected error: %v", err)
// 				return
// 			}

// 			if result != tt.expected {
// 				t.Errorf("Expected %s, got %s", tt.expected, result)
// 			}
// 		})
// 	}
// }

// func TestValidateNumber(t *testing.T) {
// 	service := NewFactorialServiceWithLimit(20000)

// 	tests := []struct {
// 		name        string
// 		number      string
// 		expected    int64
// 		expectError bool
// 	}{
// 		{
// 			name:        "Valid number 0",
// 			number:      "0",
// 			expected:    0,
// 			expectError: false,
// 		},
// 		{
// 			name:        "Valid number 100",
// 			number:      "100",
// 			expected:    100,
// 			expectError: false,
// 		},
// 		{
// 			name:        "Valid number 20000",
// 			number:      "20000",
// 			expected:    20000,
// 			expectError: false,
// 		},
// 		{
// 			name:        "Negative number",
// 			number:      "-5",
// 			expected:    0,
// 			expectError: true,
// 		},
// 		{
// 			name:        "Exceeds maximum",
// 			number:      "20001",
// 			expected:    0,
// 			expectError: true,
// 		},
// 		{
// 			name:        "Invalid format",
// 			number:      "not_a_number",
// 			expected:    0,
// 			expectError: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			result, err := service.ValidateNumber(tt.number)

// 			if tt.expectError {
// 				if err == nil {
// 					t.Errorf("Expected error but got none")
// 				}
// 				return
// 			}

// 			if err != nil {
// 				t.Errorf("Unexpected error: %v", err)
// 				return
// 			}

// 			if result != tt.expected {
// 				t.Errorf("Expected %d, got %d", tt.expected, result)
// 			}
// 		})
// 	}
// }
