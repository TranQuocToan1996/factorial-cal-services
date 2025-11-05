package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"factorial-cal-services/pkg/config"
	"factorial-cal-services/pkg/consumer"
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

	// Initialize database
	database, err := db.NewGormDB(cfg.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Connected to database")

	// Initialize Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr(),
		Password: cfg.REDIS_PASSWORD,
	})

	// Test Redis connection
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Warning: Redis connection failed: %v", err)
	} else {
		log.Println("Connected to Redis successfully")
	}
	defer redisClient.Close()

	redisService := service.NewRedisService(redisClient, 24*time.Hour, int64(cfg.REDIS_THRESHOLD))
	// Initialize services
	s3Service := service.NewS3Service(ctx, cfg)
	factorialService := service.NewFactorialService(
		repository.NewFactorialRepository(database),
		repository.NewCurrentCalculatedRepository(database),
		repository.NewMaxRequestRepository(database),
		s3Service,
	)

	// Initialize repositories
	factorialRepo := repository.NewFactorialRepository(database)
	maxRequestRepo := repository.NewMaxRequestRepository(database)
	currentCalculatedRepo := repository.NewCurrentCalculatedRepository(database)

	// Initialize RabbitMQ consumer
	mqConsumer, err := consumer.NewRabbitMQConsumer(cfg.RabbitMQURL())
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ consumer: %v", err)
	}
	defer mqConsumer.Close()

	// Create batch handler
	factorialMessageHandler := consumer.NewFactorialMessageHandler(
		factorialService,
		redisService,
		s3Service,
		factorialRepo,
		maxRequestRepo,
		currentCalculatedRepo,
	)

	batchSize := cfg.WORKER_BATCH_SIZE
	maxBatches := cfg.WORKER_MAX_BATCHES
	if maxBatches <= 0 {
		maxBatches = 16 // Default
	}
	if batchSize <= 0 {
		batchSize = 100 // Default
	}

	// Setup signal handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Worker pool with WaitGroup
	workerCount := maxBatches
	if workerCount <= 0 {
		workerCount = 1
	}

	log.Printf("Starting %d batch consumers with batch size %d", workerCount, batchSize)

	// Start consumer in a goroutine
	err = mqConsumer.Consume(ctx, cfg.FACTORIAL_CAL_SERVICES_QUEUE_NAME, factorialMessageHandler)
	if err != nil {
		log.Fatalf("Consumer error: %v", err)
	}

	log.Println("Worker started, waiting for messages...")

	// Wait for interrupt signal
	<-quit
	log.Println("Received shutdown signal, starting graceful shutdown...")

	// Wait for graceful shutdown or timeout
	<-time.After(5 * time.Second)

	log.Println("Worker stopped")
}
