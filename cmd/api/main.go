package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"factorial-cal-services/migrations"
	"factorial-cal-services/pkg/config"
	"factorial-cal-services/pkg/db"
	"factorial-cal-services/pkg/handler"
	"factorial-cal-services/pkg/producer"
	"factorial-cal-services/pkg/repository"
	"factorial-cal-services/pkg/service"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag/example/override/docs"
)

// @title           Factorial Calculation Service API
// @version         1.0
// @description     Async factorial calculation service with RabbitMQ, Redis, and S3
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @schemes http https

func main() {
	cfg := config.LoadConfig()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	// Run migrations
	if err := migrations.RunMigrations(cfg.DSN()); err != nil {
		log.Printf("Migration failed: %v", err)
	}

	// Initialize database
	database, err := db.NewGormDB(cfg.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr(),
		Password: cfg.REDIS_PASSWORD,
		DB:       cfg.REDIS_DB,
	})
	defer redisClient.Close()

	// Test Redis connection
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Printf("Warning: Redis connection failed: %v", err)
	} else {
		log.Println("Connected to Redis successfully")
	}

	// Initialize RabbitMQ producer
	mqProducer, err := producer.NewRabbitMQProducer(cfg.RabbitMQURL())
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ producer: %v", err)
	}
	defer mqProducer.Close()

	// Initialize services
	factorialService := service.NewFactorialServiceWithLimit(int64(cfg.MAX_FACTORIAL))
	redisService := service.NewRedisService(redisClient, 24*time.Hour, int64(cfg.REDIS_THRESHOLD))
	s3Service := service.NewS3Service(ctx, cfg)

	// Initialize repository
	factorialRepo := repository.NewFactorialRepository(database)

	// Initialize handler
	factorialHandler := handler.NewFactorialHandler(
		factorialService,
		redisService,
		s3Service,
		factorialRepo,
		mqProducer,
		cfg.FACTORIAL_CAL_SERVICES_QUEUE_NAME,
	)

	// Setup routes
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Swagger
	if cfg.SWAGGER_HOST != "" {
		docs.SwaggerInfo.Host = cfg.SWAGGER_HOST
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check
	r.GET("/health", healthCheck)

	// API routes
	v1 := r.Group("/api/v1")
	{
		v1.POST("/factorial", factorialHandler.SubmitCalculation)
		v1.GET("/factorial/:number", factorialHandler.GetResult)
		v1.GET("/factorial/metadata/:number", factorialHandler.GetMetadata)
	}

	srv := &http.Server{
		Addr:    cfg.SERVER_PORT,
		Handler: r,
	}

	// Run server in goroutine
	go func() {
		log.Printf("API server listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down API server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("API server stopped")
}

// healthCheck godoc
// @Summary      Health check
// @Description  Check if the service is running
// @Tags         health
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /health [get]
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
