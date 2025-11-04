package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	DefaultRedisTTL = 24 * time.Hour
	RedisKeyPrefix  = "factorial:"
)

// RedisService handles Redis caching operations
type RedisService interface {
	Get(ctx context.Context, number string) (string, error)
	Set(ctx context.Context, number string, result string) error
	ShouldCache(number string) bool
}

type redisService struct {
	client    *redis.Client
	ttl       time.Duration
	threshold int64
}

// NewRedisService creates a new Redis service
func NewRedisService(client *redis.Client, ttl time.Duration, threshold int64) RedisService {
	if ttl == 0 {
		ttl = DefaultRedisTTL
	}
	if threshold == 0 {
		threshold = 1000 // Default threshold
	}
	return &redisService{
		client:    client,
		ttl:       ttl,
		threshold: threshold,
	}
}

// factorialKey creates a Redis key with prefix
func (s *redisService) factorialKey(number string) string {
	return fmt.Sprintf("%s%s", RedisKeyPrefix, number)
}

// ShouldCache determines if a number should be cached based on threshold
func (s *redisService) ShouldCache(number string) bool {
	// Parse number directly without validation limits (threshold check only)
	n, err := strconv.ParseInt(number, 10, 64)
	if err != nil {
		return false
	}
	// Only check if negative
	if n < 0 {
		return false
	}
	return n < s.threshold
}

// Get retrieves a factorial result from Redis cache
func (s *redisService) Get(ctx context.Context, number string) (string, error) {
	key := s.factorialKey(number)
	result, err := s.client.Get(ctx, key).Result()

	if err == redis.Nil {
		return "", nil // Cache miss, not an error
	}

	if err != nil {
		return "", fmt.Errorf("redis get error: %w", err)
	}

	return result, nil
}

// Set stores a factorial result in Redis cache with TTL
func (s *redisService) Set(ctx context.Context, number string, result string) error {
	// Only cache if below threshold
	if !s.ShouldCache(number) {
		return nil
	}

	key := s.factorialKey(number)
	err := s.client.Set(ctx, key, result, s.ttl).Err()
	if err != nil {
		return fmt.Errorf("redis set error: %w", err)
	}

	return nil
}
