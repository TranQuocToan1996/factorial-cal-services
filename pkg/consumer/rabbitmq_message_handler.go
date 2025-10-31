package consumer

import (
	"context"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// handleMessage handles a single message with panic recovery and retry logic
func (c *RabbitMQConsumer) handleMessage(ctx context.Context, msg amqp.Delivery, queueName string, handler MessageHandler) {
	defer func() {
		if r := recover(); r != nil {
			c.handlePanic(r, msg, queueName)
		}
	}()

	log.Printf("Received message: message_id: %s, body: %s", msg.MessageId, string(msg.Body))

	retryCount := c.getRetryCount(msg.Headers)

	if err := handler(ctx, msg.Body); err != nil {
		c.handleError(err, msg, queueName, retryCount)
	} else {
		msg.Ack(false)
		log.Printf("Message processed successfully: message_id: %s, body: %s", msg.MessageId, string(msg.Body))
	}
}

func (c *RabbitMQConsumer) handlePanic(r any, msg amqp.Delivery, queueName string) {
	log.Printf("PANIC recovered while processing message: %v", r)
	log.Printf("Message body: %s", string(msg.Body))

	retryCount := c.getRetryCount(msg.Headers)

	if retryCount >= 3 {
		log.Printf("Max retries reached after panic, sending to DLQ")
		if err := c.sendToDLQ(queueName+".dlq", msg, fmt.Errorf("panic: %v", r)); err != nil {
			log.Printf("Failed to send to DLQ: %v", err)
			msg.Nack(false, false)
		} else {
			msg.Ack(false)
		}
	} else {
		log.Printf("Nacking message after panic for retry (attempt %d)", retryCount+1)
		msg.Nack(false, false)
	}
}

func (c *RabbitMQConsumer) handleError(err error, msg amqp.Delivery, queueName string, retryCount int) {
	log.Printf("Error processing message (retry %d/3): %v", retryCount, err)

	if retryCount >= 3 {
		log.Printf("Max retries reached, sending to DLQ")
		if err := c.sendToDLQ(queueName+".dlq", msg, err); err != nil {
			log.Printf("Failed to send to DLQ: %v", err)
			msg.Nack(false, false)
		} else {
			msg.Ack(false)
		}
	} else {
		log.Printf("Nacking message for retry (attempt %d)", retryCount+1)
		msg.Nack(false, false)
	}
}
