package config

import (
	"fmt"
	"os"
)

func LoadConfig() *Config {
	return &Config{
		SERVER_PORT:                       os.Getenv("SERVER_PORT"),
		DB_USER:                           os.Getenv("DB_USER"),
		DB_PASSWORD:                       os.Getenv("DB_PASSWORD"),
		DB_HOST:                           os.Getenv("DB_HOST"),
		DB_PORT:                           os.Getenv("DB_PORT"),
		DB_NAME:                           os.Getenv("DB_NAME"),
		DB_SSLMODE:                        os.Getenv("DB_SSLMODE"),
		RABBITMQ_USER:                     os.Getenv("RABBITMQ_USER"),
		RABBITMQ_PASSWORD:                 os.Getenv("RABBITMQ_PASSWORD"),
		RABBITMQ_HOST:                     os.Getenv("RABBITMQ_HOST"),
		RABBITMQ_PORT:                     os.Getenv("RABBITMQ_PORT"),
		FACTORIAL_CAL_SERVICES_QUEUE_NAME: os.Getenv("FACTORIAL_CAL_SERVICES_QUEUE_NAME"),
		SWAGGER_HOST:                      os.Getenv("SWAGGER_HOST"),
		RABBITMQ_CA:                       os.Getenv("RABBITMQ_CA"),
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
	RABBITMQ_USER                     string `mapstructure:"RABBITMQ_USER"`
	RABBITMQ_PASSWORD                 string `mapstructure:"RABBITMQ_PASSWORD"`
	RABBITMQ_HOST                     string `mapstructure:"RABBITMQ_HOST"`
	RABBITMQ_PORT                     string `mapstructure:"RABBITMQ_PORT"`
	FACTORIAL_CAL_SERVICES_QUEUE_NAME string `mapstructure:"FACTORIAL_CAL_SERVICES_QUEUE_NAME"`
	SWAGGER_HOST                      string `mapstructure:"SWAGGER_HOST"`
	RABBITMQ_CA                       string `mapstructure:"RABBITMQ_CA"`
}

func (c *Config) DSN() string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		c.DB_USER, c.DB_PASSWORD, c.DB_HOST, c.DB_PORT, c.DB_NAME, c.DB_SSLMODE)
}

func (c *Config) RabbitMQURL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s",
		c.RABBITMQ_USER, c.RABBITMQ_PASSWORD, c.RABBITMQ_HOST, c.RABBITMQ_PORT)
}
