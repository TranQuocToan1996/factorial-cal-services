package producer

import "context"

// Producer represents a queue producer (SQS, Kafka, RabbitMQ, etc.)
type Producer interface {
	Publish(ctx context.Context, queueName string, payload []byte) error
	Close() error
}
