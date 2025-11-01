package consumer

import (
	"fmt"
)

// setupQueues declares the main queue
func (c *RabbitMQConsumer) setupQueues(queueName string) error {
	_, err := c.channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // no arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}
	return nil
}
