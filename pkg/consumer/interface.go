package consumer

import "context"

// Consumer represents a queue consumer (SQS, Kafka, RabbitMQ, etc.)
type Consumer interface {
	Consume(ctx context.Context, queueName string, handler MessageHandler) error
	Close() error
}

// MessageHandler is a function type for processing messages
type MessageHandler func(ctx context.Context, payload []byte) error
