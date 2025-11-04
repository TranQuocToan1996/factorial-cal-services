package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"factorial-cal-services/pkg/config"
	"factorial-cal-services/pkg/db"
	"factorial-cal-services/pkg/repository"
	"factorial-cal-services/pkg/service"

	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.LoadConfig()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	// Initialize Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr(),
		Password: cfg.REDIS_PASSWORD,
	})

	// Test Redis connection
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Printf("Warning: Redis connection failed: %v", err)
	} else {
		log.Println("Connected to Redis successfully")
	}
	defer redisClient.Close()

	redisService := service.NewRedisService(redisClient, 24*time.Hour, int64(cfg.REDIS_THRESHOLD))

	database, err := db.NewGormDB(cfg.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Connected to database")

	factorialService := service.NewFactorialService(
		repository.NewFactorialRepository(database),
		repository.NewCurrentCalculatedRepository(database),
		repository.NewMaxRequestRepository(database),
		redisService,
		service.NewS3Service(ctx, cfg),
	)

	factorialService.StartContinuelyCalculateFactorial()

	log.Println("Calculate started, waiting for messages...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// Wait for interrupt signal
	<-quit
	log.Println("Received shutdown signal, starting graceful shutdown...")

	// Wait for graceful shutdown or timeout
	<-time.After(5 * time.Second)

	log.Println("Worker stopped")
}
