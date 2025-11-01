package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"factorial-cal-services/migrations"
	"factorial-cal-services/pkg/config"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag/example/override/docs"
)

// @title           Simple Order Service API
// @version         1.0
// @description     A simple order service with outbox/inbox pattern for reliable message processing
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @schemes http https

func main() {
	cfg := config.LoadConfig()

	// db, err := db.NewGormDB(cfg.DSN())
	// if err != nil {
	// 	log.Fatalf("Failed to connect to database: %v", err)
	// }

	if err := migrations.RunMigrations(cfg.DSN()); err != nil {
		log.Printf("Migration failed: %v", err)
	}

	// Setup routes
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Group("/api/v1").
		Use(gin.Logger()).
		Use(gin.Recovery())

	if cfg.SWAGGER_HOST != "" {
		docs.SwaggerInfo.Host = cfg.SWAGGER_HOST
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check endpoints
	r.GET("/health", healthCheck)

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

	srv.Shutdown(context.Background())

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
