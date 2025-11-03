package consumer

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQConsumer implements Consumer for RabbitMQ
type RabbitMQConsumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewRabbitMQConsumer creates a new RabbitMQ consumer
func NewRabbitMQConsumer(amqpURL string) (*RabbitMQConsumer, error) {
	var conn *amqp.Connection
	var err error

	// Check if TLS is enabled by looking for RABBITMQ_CA
	rootCA := os.Getenv("RABBITMQ_CA")
	if rootCA != "" {
		// Replace literal \n with actual newlines
		rootCA = strings.ReplaceAll(rootCA, "\\n", "\n")
		fmt.Println("rootCA", rootCA)

		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM([]byte(rootCA)) {
			return nil, fmt.Errorf("failed to append root CA cert")
		}

		tlsConfig := &tls.Config{
			RootCAs:            caCertPool,
			InsecureSkipVerify: false, // keep secure!
		}

		conn, err = amqp.DialTLS(amqpURL, tlsConfig)
	} else {
		// Use plain connection for local development
		conn, err = amqp.Dial(amqpURL)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Set QoS to allow prefetching for batch processing
	// Increase prefetch count to batch size for better throughput
	err = ch.Qos(
		100,   // prefetch count (allow multiple messages for batching)
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to set QoS: %w", err)
	}

	return &RabbitMQConsumer{
		conn:    conn,
		channel: ch,
	}, nil
}

// Consume starts consuming messages from the specified queue
func (c *RabbitMQConsumer) Consume(ctx context.Context, queueName string, handler MessageHandler) error {
	// Setup queue
	if err := c.setupQueues(queueName); err != nil {
		return err
	}

	// Start consuming messages
	msgs, err := c.channel.Consume(
		queueName, // queue
		"",        // consumer tag
		false,     // auto-ack = false (manual ack)
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	log.Printf("Started consuming from queue: %s", queueName)

	// Process messages in goroutine
	go c.consumeLoop(ctx, msgs, handler)

	return nil
}

// consumeLoop handles the main message consumption loop
func (c *RabbitMQConsumer) consumeLoop(ctx context.Context, msgs <-chan amqp.Delivery, handler MessageHandler) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Consumer context cancelled, stopping...")
			return
		case msg, ok := <-msgs:
			if !ok {
				log.Println("Message channel closed")
				return
			}
			c.handleMessageAndAck(ctx, msg, handler)
	}
}

// ConsumeBatch starts consuming messages in batches from the specified queue
func (c *RabbitMQConsumer) ConsumeBatch(ctx context.Context, queueName string, batchSize int, handler BatchMessageHandler) error {
	// Setup queue
	if err := c.setupQueues(queueName); err != nil {
		return err
	}

	// Start consuming messages
	msgs, err := c.channel.Consume(
		queueName, // queue
		"",        // consumer tag
		false,     // auto-ack = false (manual ack)
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	log.Printf("Started batch consuming from queue: %s (batch size: %d)", queueName, batchSize)

	// Process messages in goroutine
	go c.consumeBatchLoop(ctx, msgs, batchSize, handler)

	return nil
}

// consumeBatchLoop handles batch message consumption
func (c *RabbitMQConsumer) consumeBatchLoop(ctx context.Context, msgs <-chan amqp.Delivery, batchSize int, handler BatchMessageHandler) {
	batch := make([][]byte, 0, batchSize)
	deliveries := make([]amqp.Delivery, 0, batchSize)
	
	for {
		select {
		case <-ctx.Done():
			// Process remaining messages in batch before stopping
			if len(batch) > 0 {
				c.processBatch(ctx, batch, deliveries, handler)
			}
			log.Println("Consumer context cancelled, stopping...")
			return
		case msg, ok := <-msgs:
			if !ok {
				// Process remaining messages in batch before closing
				if len(batch) > 0 {
					c.processBatch(ctx, batch, deliveries, handler)
				}
				log.Println("Message channel closed")
				return
			}
			
			batch = append(batch, msg.Body)
			deliveries = append(deliveries, msg)
			
			// Process batch when it reaches the desired size
			if len(batch) >= batchSize {
				c.processBatch(ctx, batch, deliveries, handler)
				batch = batch[:0]
				deliveries = deliveries[:0]
			}
		}
	}
}

// processBatch processes a batch of messages and acks/nacks them accordingly
func (c *RabbitMQConsumer) processBatch(ctx context.Context, batch [][]byte, deliveries []amqp.Delivery, handler BatchMessageHandler) {
	if len(batch) == 0 {
		return
	}
	
	log.Printf("Processing batch of %d messages", len(batch))
	
	if err := handler(ctx, batch); err != nil {
		log.Printf("Error processing batch: %v", err)
		// Nack all messages in batch (don't requeue)
		for _, msg := range deliveries {
			msg.Nack(false, false)
		}
	} else {
		// Ack all messages in batch
		for _, msg := range deliveries {
			msg.Ack(false)
		}
		log.Printf("Batch processed successfully: %d messages", len(batch))
	}
}

// Close closes the RabbitMQ connection and channel
func (c *RabbitMQConsumer) Close() error {
	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			return err
		}
	}
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return err
		}
	}
	return nil
}

// IsHealthy checks if the RabbitMQ connection is alive
func (c *RabbitMQConsumer) IsHealthy() bool {
	return c.conn != nil && !c.conn.IsClosed() && c.channel != nil
}
