package consumer

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// handleMessageAndAck handles a single message with panic recovery
func (c *RabbitMQConsumer) handleMessageAndAck(ctx context.Context, msg amqp.Delivery, handler MessageHandler) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("PANIC recovered while processing message: %v", r)
			log.Printf("Message body: %s", string(msg.Body))
			msg.Nack(false, false) // Reject and don't requeue
		}
	}()

	log.Printf("Received message: message_id: %s, body: %s", msg.MessageId, string(msg.Body))

	if err := handler(ctx, msg.Body); err != nil {
		log.Printf("Error processing message: %v", err)
		msg.Nack(false, false) // Reject and don't requeue
	} else {
		msg.Ack(false)
		log.Printf("Message processed successfully: message_id: %s, body: %s", msg.MessageId, string(msg.Body))
	}
}
