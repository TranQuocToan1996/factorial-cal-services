package service

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/alicebob/miniredis/v2"
)

func TestShouldCache(t *testing.T) {
	// Create a mock Redis server
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to create miniredis: %v", err)
	}
	defer mr.Close()

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer client.Close()

	service := NewRedisService(client, time.Hour)

	tests := []struct {
		name     string
		number   string
		expected bool
	}{
		{
			name:     "Small number 0",
			number:   "0",
			expected: true,
		},
		{
			name:     "Small number 9999",
			number:   "9999",
			expected: true,
		},
		{
			name:     "Boundary number 10000",
			number:   "10000",
			expected: false,
		},
		{
			name:     "Large number 15000",
			number:   "15000",
			expected: false,
		},
		{
			name:     "Invalid number",
			number:   "invalid",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.ShouldCache(tt.number)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestRedisSetAndGet(t *testing.T) {
	// Create a mock Redis server
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to create miniredis: %v", err)
	}
	defer mr.Close()

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer client.Close()

	service := NewRedisService(client, time.Hour)
	ctx := context.Background()

	tests := []struct {
		name   string
		number string
		value  string
	}{
		{
			name:   "Set and get factorial of 5",
			number: "5",
			value:  "120",
		},
		{
			name:   "Set and get factorial of 10",
			number: "10",
			value:  "3628800",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set value
			err := service.Set(ctx, tt.number, tt.value)
			if err != nil {
				t.Errorf("Failed to set value: %v", err)
				return
			}

			// Get value
			result, err := service.Get(ctx, tt.number)
			if err != nil {
				t.Errorf("Failed to get value: %v", err)
				return
			}

			if result != tt.value {
				t.Errorf("Expected %s, got %s", tt.value, result)
			}
		})
	}
}

func TestRedisGetNonExistent(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to create miniredis: %v", err)
	}
	defer mr.Close()

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer client.Close()

	service := NewRedisService(client, time.Hour)
	ctx := context.Background()

	result, err := service.Get(ctx, "999")
	if err != nil {
		t.Errorf("Expected no error for cache miss, got: %v", err)
	}

	if result != "" {
		t.Errorf("Expected empty string for cache miss, got: %s", result)
	}
}

