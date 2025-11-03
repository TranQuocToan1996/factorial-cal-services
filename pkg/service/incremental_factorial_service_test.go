package service

// func TestCalculateIncremental(t *testing.T) {
// 	factorialService := NewFactorialService()
// 	service := NewIncrementalFactorialService(factorialService)
// 	ctx := context.Background()

// 	tests := []struct {
// 		name          string
// 		curNumber     string
// 		curFactorial  string
// 		targetNumber  string
// 		expected      string
// 		expectError   bool
// 		expectedError string
// 	}{
// 		{
// 			name:         "Calculate from 5 to 10",
// 			curNumber:    "5",
// 			curFactorial: "120", // 5! = 120
// 			targetNumber: "10",
// 			expected:     "3628800", // 10! = 3628800
// 			expectError:  false,
// 		},
// 		{
// 			name:         "Calculate from 0 to 5",
// 			curNumber:    "0",
// 			curFactorial: "1", // 0! = 1
// 			targetNumber: "5",
// 			expected:     "120", // 5! = 120
// 			expectError:  false,
// 		},
// 		{
// 			name:         "Calculate from 10 to 15",
// 			curNumber:    "10",
// 			curFactorial: "3628800", // 10! = 3628800
// 			targetNumber: "15",
// 			expected:     "1307674368000", // 15! = 1307674368000
// 			expectError:  false,
// 		},
// 		{
// 			name:         "Already calculated (cur >= target)",
// 			curNumber:    "10",
// 			curFactorial: "3628800",
// 			targetNumber: "10",
// 			expected:     "3628800",
// 			expectError:  false,
// 		},
// 		{
// 			name:         "Current greater than target",
// 			curNumber:    "15",
// 			curFactorial: "1307674368000",
// 			targetNumber: "10",
// 			expected:     "1307674368000",
// 			expectError:  false,
// 		},
// 		{
// 			name:          "Invalid current number",
// 			curNumber:     "invalid",
// 			curFactorial:  "120",
// 			targetNumber:  "10",
// 			expectError:   true,
// 			expectedError: "invalid current number",
// 		},
// 		{
// 			name:          "Invalid target number",
// 			curNumber:     "5",
// 			curFactorial:  "120",
// 			targetNumber:  "invalid",
// 			expectError:   true,
// 			expectedError: "invalid target number",
// 		},
// 		{
// 			name:          "Invalid factorial format",
// 			curNumber:     "5",
// 			curFactorial:  "not-a-number",
// 			targetNumber:  "10",
// 			expectError:   true,
// 			expectedError: "invalid current factorial format",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			result, err := service.CalculateIncremental(ctx, tt.curNumber, tt.curFactorial, tt.targetNumber)

// 			if tt.expectError {
// 				if err == nil {
// 					t.Errorf("Expected error but got none")
// 					return
// 				}
// 				if tt.expectedError != "" && err.Error()[:len(tt.expectedError)] != tt.expectedError {
// 					t.Errorf("Expected error containing '%s', got '%s'", tt.expectedError, err.Error())
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

// func TestCalculateIncrementalLargeNumbers(t *testing.T) {
// 	factorialService := NewFactorialServiceWithLimit(10000)
// 	db := setupTestDBForCurrentCalculated(t)
// 	currentCalRepo := repository.NewCurrentCalculatedRepository(db)
// 	service := NewIncrementalFactorialService(factorialService, currentCalRepo)
// 	ctx := context.Background()

// 	// Test with large numbers using big.Int
// 	result, err := service.CalculateIncremental(ctx, "100", "93326215443944152681699238856266700490715968264381621468592963895217599993229915608941463976156518286253697920827223758251185210916864000000000000000000000000", "105")
// 	if err != nil {
// 		t.Fatalf("Unexpected error: %v", err)
// 	}

// 	if len(result) == 0 {
// 		t.Error("Expected non-empty result for large number calculation")
// 	}

// 	// Verify it's a valid number (contains only digits)
// 	for _, char := range result {
// 		if char < '0' || char > '9' {
// 			t.Errorf("Result contains non-digit character: %c", char)
// 		}
// 	}
// }
