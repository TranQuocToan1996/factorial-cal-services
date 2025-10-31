package consumer

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// setupQueues declares all required queues, exchanges, and bindings for retry pattern
func (c *RabbitMQConsumer) setupQueues(queueName string) error {
	// Step 1: Declare retry exchange
	if err := c.declareRetryExchange(queueName); err != nil {
		return err
	}

	// Step 2: Declare retry queue with TTL
	if err := c.declareRetryQueue(queueName); err != nil {
		return err
	}

	// Step 3: Bind retry queue to retry exchange
	if err := c.bindRetryQueue(queueName); err != nil {
		return err
	}

	// Step 4: Declare main queue with DLX
	if err := c.declareMainQueue(queueName); err != nil {
		return err
	}

	// Step 5: Declare DLQ
	if err := c.declareDLQ(queueName); err != nil {
		return err
	}

	return nil
}

func (c *RabbitMQConsumer) declareRetryExchange(queueName string) error {
	err := c.channel.ExchangeDeclare(
		queueName+".retry.exchange", // name
		"direct",                    // type
		true,                        // durable
		false,                       // auto-deleted
		false,                       // internal
		false,                       // no-wait
		nil,                         // args
	)
	if err != nil {
		return fmt.Errorf("failed to declare retry exchange: %w", err)
	}
	return nil
}

func (c *RabbitMQConsumer) declareRetryQueue(queueName string) error {
	retryArgs := amqp.Table{
		"x-message-ttl":             5000,      // 5 seconds
		"x-dead-letter-exchange":    "",        // default exchange
		"x-dead-letter-routing-key": queueName, // back to main queue
	}

	_, err := c.channel.QueueDeclare(
		queueName+".retry", // name
		true,               // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		retryArgs,          // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare retry queue: %w", err)
	}
	return nil
}

func (c *RabbitMQConsumer) bindRetryQueue(queueName string) error {
	err := c.channel.QueueBind(
		queueName+".retry",          // queue name
		queueName+".retry",          // routing key
		queueName+".retry.exchange", // exchange
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind retry queue: %w", err)
	}
	return nil
}

func (c *RabbitMQConsumer) declareMainQueue(queueName string) error {
	mainArgs := amqp.Table{
		"x-dead-letter-exchange":    queueName + ".retry.exchange", // to retry
		"x-dead-letter-routing-key": queueName + ".retry",          // routing key
	}

	_, err := c.channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		mainArgs,  // arguments with DLX
	)
	if err != nil {
		return fmt.Errorf("failed to declare main queue: %w", err)
	}
	return nil
}

func (c *RabbitMQConsumer) declareDLQ(queueName string) error {
	_, err := c.channel.QueueDeclare(
		queueName+".dlq", // name
		true,             // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // no arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare DLQ: %w", err)
	}
	return nil
}
