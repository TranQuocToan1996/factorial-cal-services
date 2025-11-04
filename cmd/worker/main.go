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
		log.Printf("Warning: Redis connection failed: %v", err)
	} else {
		log.Println("Connected to Redis successfully")
	}
	defer redisClient.Close()

	// Initialize services
	factorialService := service.NewFactorialServiceWithLimit(int64(cfg.MAX_FACTORIAL))
	redisService := service.NewRedisService(redisClient, 24*time.Hour, int64(cfg.REDIS_THRESHOLD))
	s3Service := service.NewS3Service(ctx, cfg)

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

	incrementalService := service.NewIncrementalFactorialService(factorialService, currentCalculatedRepo)

	// Create batch handler
	batchHandler := consumer.NewFactorialBatchHandler(
		factorialService,
		redisService,
		s3Service,
		factorialRepo,
		maxRequestRepo,
		currentCalculatedRepo,
		incrementalService,
	)

	batchSize := cfg.WORKER_BATCH_SIZE
	maxBatches := cfg.WORKER_MAX_BATCHES
	if maxBatches <= 0 {
		maxBatches = 16 // Default
	}
	if batchSize <= 0 {
		batchSize = 100 // Default
	}

	log.Printf("Starting %d batch consumers with batch size %d", maxBatches, batchSize)

	// Start multiple batch consumers concurrently
	for i := 0; i < maxBatches; i++ {
		go func(batchID int) {
			if err := mqConsumer.ConsumeBatch(ctx, cfg.FACTORIAL_CAL_SERVICES_QUEUE_NAME, batchSize, batchHandler); err != nil {
				log.Fatalf("Failed to start batch consumer %d: %v", batchID, err)
			}
			log.Printf("Batch consumer %d started", batchID)
		}(i)
	}

	log.Println("Worker started, waiting for messages...")

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down worker...")

	// Give time for current message processing to complete
	time.Sleep(2 * time.Second)

	log.Println("Worker stopped")
}
