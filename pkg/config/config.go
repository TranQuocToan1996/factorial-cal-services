package config

import (
	"fmt"
	"os"
	"strconv"
)

func LoadConfig() *Config {
	maxFactorial, _ := strconv.Atoi(getEnvOrDefault("MAX_FACTORIAL", "10000"))
	redisThreshold, _ := strconv.Atoi(getEnvOrDefault("REDIS_THRESHOLD", "1000"))
	workerBatchSize, _ := strconv.Atoi(getEnvOrDefault("WORKER_BATCH_SIZE", "100"))
	workerMaxBatches, _ := strconv.Atoi(getEnvOrDefault("WORKER_MAX_BATCHES", "16"))

	return &Config{
		SERVER_PORT:                       os.Getenv("SERVER_PORT"),
		DB_USER:                           os.Getenv("DB_USER"),
		DB_PASSWORD:                       os.Getenv("DB_PASSWORD"),
		DB_HOST:                           os.Getenv("DB_HOST"),
		DB_PORT:                           os.Getenv("DB_PORT"),
		DB_NAME:                           os.Getenv("DB_NAME"),
		DB_SSLMODE:                        os.Getenv("DB_SSLMODE"),
		DB_TYPE:                           getEnvOrDefault("DB_TYPE", "postgres"),
		RABBITMQ_USER:                     os.Getenv("RABBITMQ_USER"),
		RABBITMQ_PASSWORD:                 os.Getenv("RABBITMQ_PASSWORD"),
		RABBITMQ_HOST:                     os.Getenv("RABBITMQ_HOST"),
		RABBITMQ_PORT:                     os.Getenv("RABBITMQ_PORT"),
		FACTORIAL_CAL_SERVICES_QUEUE_NAME: os.Getenv("FACTORIAL_CAL_SERVICES_QUEUE_NAME"),
		SWAGGER_HOST:                      os.Getenv("SWAGGER_HOST"),
		RABBITMQ_CA:                       os.Getenv("RABBITMQ_CA"),
		REDIS_HOST:                        os.Getenv("REDIS_HOST"),
		REDIS_PORT:                        os.Getenv("REDIS_PORT"),
		REDIS_PASSWORD:                    os.Getenv("REDIS_PASSWORD"),
		AWS_REGION:                        os.Getenv("AWS_REGION"),
		S3_BUCKET_NAME:                    os.Getenv("S3_BUCKET_NAME"),
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
	return fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		c.DB_USER, c.DB_PASSWORD, c.DB_HOST, c.DB_PORT, c.DB_NAME, c.DB_SSLMODE)
}

func (c *Config) RabbitMQURL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s",
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
