package config

import (
	"fmt"
	"os"
	"strconv"
)

func LoadConfig() *Config {
	maxFactorial, _ := strconv.Atoi(getEnvOrDefault("MAX_FACTORIAL", "100"))
	redisThreshold, _ := strconv.Atoi(getEnvOrDefault("REDIS_THRESHOLD", "50"))
	workerBatchSize, _ := strconv.Atoi(getEnvOrDefault("WORKER_BATCH_SIZE", "100"))
	workerMaxBatches, _ := strconv.Atoi(getEnvOrDefault("WORKER_MAX_BATCHES", "1"))

	return &Config{
		SERVER_PORT:                       getEnvOrDefault("SERVER_PORT", ":8080"),
		DB_USER:                           getEnvOrDefault("DB_USER", "postgres"),
		DB_PASSWORD:                       getEnvOrDefault("DB_PASSWORD", "password"),
		DB_HOST:                           getEnvOrDefault("DB_HOST", "localhost"),
		DB_PORT:                           getEnvOrDefault("DB_PORT", "5432"),
		DB_NAME:                           getEnvOrDefault("DB_NAME", "factorial-cal-services"),
		DB_SSLMODE:                        getEnvOrDefault("DB_SSLMODE", "disable"),
		DB_TYPE:                           getEnvOrDefault("DB_TYPE", "postgres"),
		RABBITMQ_USER:                     getEnvOrDefault("RABBITMQ_USER", "guest"),
		RABBITMQ_PASSWORD:                 getEnvOrDefault("RABBITMQ_PASSWORD", "guest"),
		RABBITMQ_HOST:                     getEnvOrDefault("RABBITMQ_HOST", "localhost"),
		RABBITMQ_PORT:                     getEnvOrDefault("RABBITMQ_PORT", "5672"),
		FACTORIAL_CAL_SERVICES_QUEUE_NAME: getEnvOrDefault("FACTORIAL_CAL_SERVICES_QUEUE_NAME", "factorial-cal-queue"),
		SWAGGER_HOST:                      getEnvOrDefault("SWAGGER_HOST", "localhost:8080"),
		RABBITMQ_CA:                       getEnvOrDefault("RABBITMQ_CA", ""),
		REDIS_HOST:                        getEnvOrDefault("REDIS_HOST", "localhost"),
		REDIS_PORT:                        getEnvOrDefault("REDIS_PORT", "6379"),
		REDIS_PASSWORD:                    getEnvOrDefault("REDIS_PASSWORD", ""),
		AWS_REGION:                        getEnvOrDefault("AWS_REGION", "us-east-1"),
		S3_BUCKET_NAME:                    getEnvOrDefault("S3_BUCKET_NAME", "factorial-calculator-service"),
		STORAGE_TYPE:                      getEnvOrDefault("STORAGE_TYPE", "local"),
		QUEUE_TYPE:                        getEnvOrDefault("QUEUE_TYPE", "rabbitmq"),
		MAX_FACTORIAL:                     maxFactorial,
		REDIS_THRESHOLD:                   redisThreshold,
		WORKER_BATCH_SIZE:                 workerBatchSize,
		WORKER_MAX_BATCHES:                workerMaxBatches,
	}
}

type Config struct {
	SERVER_PORT                       string `mapstructure:"SERVER_PORT"`
	DB_USER                           string `mapstructure:"DB_USER"`
	DB_PASSWORD                       string `mapstructure:"DB_PASSWORD"`
	DB_HOST                           string `mapstructure:"DB_HOST"`
	DB_PORT                           string `mapstructure:"DB_PORT"`
	DB_NAME                           string `mapstructure:"DB_NAME"`
	DB_SSLMODE                        string `mapstructure:"DB_SSLMODE"`
	DB_TYPE                           string `mapstructure:"DB_TYPE"`
	RABBITMQ_USER                     string `mapstructure:"RABBITMQ_USER"`
	RABBITMQ_PASSWORD                 string `mapstructure:"RABBITMQ_PASSWORD"`
	RABBITMQ_HOST                     string `mapstructure:"RABBITMQ_HOST"`
	RABBITMQ_PORT                     string `mapstructure:"RABBITMQ_PORT"`
	FACTORIAL_CAL_SERVICES_QUEUE_NAME string `mapstructure:"FACTORIAL_CAL_SERVICES_QUEUE_NAME"`
	SWAGGER_HOST                      string `mapstructure:"SWAGGER_HOST"`
	RABBITMQ_CA                       string `mapstructure:"RABBITMQ_CA"`
	REDIS_HOST                        string `mapstructure:"REDIS_HOST"`
	REDIS_PORT                        string `mapstructure:"REDIS_PORT"`
	REDIS_PASSWORD                    string `mapstructure:"REDIS_PASSWORD"`
	AWS_REGION                        string `mapstructure:"AWS_REGION"`
	S3_BUCKET_NAME                    string `mapstructure:"S3_BUCKET_NAME"`
	STORAGE_TYPE                      string `mapstructure:"STORAGE_TYPE"`
	QUEUE_TYPE                        string `mapstructure:"QUEUE_TYPE"`
	MAX_FACTORIAL                     int    `mapstructure:"MAX_FACTORIAL"`
	REDIS_THRESHOLD                   int    `mapstructure:"REDIS_THRESHOLD"`
	WORKER_BATCH_SIZE                 int    `mapstructure:"WORKER_BATCH_SIZE"`
	WORKER_MAX_BATCHES                int    `mapstructure:"WORKER_MAX_BATCHES"`
}

func (c *Config) DSN() string {
	// postgres://postgres:secret@localhost:5432/mydb?sslmode=disable
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DB_USER, c.DB_PASSWORD, c.DB_HOST, c.DB_PORT, c.DB_NAME, c.DB_SSLMODE)
}

func (c *Config) RabbitMQURL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s",
		c.RABBITMQ_USER, c.RABBITMQ_PASSWORD, c.RABBITMQ_HOST, c.RABBITMQ_PORT)
}

func (c *Config) RabbitMQURLWithSecure() string {
	return fmt.Sprintf("amqps://%s:%s@%s:%s",
		c.RABBITMQ_USER, c.RABBITMQ_PASSWORD, c.RABBITMQ_HOST, c.RABBITMQ_PORT)
}

func (c *Config) RedisAddr() string {
	return fmt.Sprintf("%s:%s", c.REDIS_HOST, c.REDIS_PORT)
}

// Validate validates the configuration and returns an error if required fields are missing or invalid
func (c *Config) Validate() error {
	// Required database fields
	if c.DB_HOST == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if c.DB_PORT == "" {
		return fmt.Errorf("DB_PORT is required")
	}
	if c.DB_NAME == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	if c.DB_USER == "" {
		return fmt.Errorf("DB_USER is required")
	}

	// Validate DB_PORT is numeric
	if _, err := strconv.Atoi(c.DB_PORT); err != nil {
		return fmt.Errorf("DB_PORT must be numeric: %w", err)
	}

	// Required RabbitMQ fields
	if c.RABBITMQ_HOST == "" {
		return fmt.Errorf("RABBITMQ_HOST is required")
	}
	if c.RABBITMQ_PORT == "" {
		return fmt.Errorf("RABBITMQ_PORT is required")
	}

	// Validate RABBITMQ_PORT is numeric
	if _, err := strconv.Atoi(c.RABBITMQ_PORT); err != nil {
		return fmt.Errorf("RABBITMQ_PORT must be numeric: %w", err)
	}

	// Validate batch sizes are positive
	if c.WORKER_BATCH_SIZE <= 0 {
		return fmt.Errorf("WORKER_BATCH_SIZE must be positive (got %d)", c.WORKER_BATCH_SIZE)
	}
	if c.WORKER_MAX_BATCHES <= 0 {
		return fmt.Errorf("WORKER_MAX_BATCHES must be positive (got %d)", c.WORKER_MAX_BATCHES)
	}

	// Validate MAX_FACTORIAL is positive
	if c.MAX_FACTORIAL <= 0 {
		return fmt.Errorf("MAX_FACTORIAL must be positive (got %d)", c.MAX_FACTORIAL)
	}

	// Validate REDIS_THRESHOLD is positive
	if c.REDIS_THRESHOLD <= 0 {
		return fmt.Errorf("REDIS_THRESHOLD must be positive (got %d)", c.REDIS_THRESHOLD)
	}

	// Warn about optional but recommended fields
	if c.AWS_REGION == "" || c.S3_BUCKET_NAME == "" {
		// Log warning but don't fail (S3 might not be configured for local development)
		fmt.Println("Warning: AWS_REGION and S3_BUCKET_NAME are not set. S3 operations will fail if attempted.")
	}

	return nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
