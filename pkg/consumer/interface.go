package consumer

import "context"

// Consumer represents a queue consumer (SQS, Kafka, RabbitMQ, etc.)
type Consumer interface {
	Consume(ctx context.Context, queueName string, handler MessageHandler) error
	ConsumeBatch(ctx context.Context, queueName string, batchSize int, handler BatchMessageHandler) error
	Close() error
}

// MessageHandler is a function type for processing messages
type MessageHandler func(ctx context.Context, payload []byte) error

// BatchMessageHandler is a function type for processing batches of messages
type BatchMessageHandler func(ctx context.Context, payloads [][]byte) error
