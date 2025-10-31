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

	// Set QoS to process one message at a time
	err = ch.Qos(
		1,     // prefetch count
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

// Consume starts consuming messages from the specified queue with retry and DLQ support
func (c *RabbitMQConsumer) Consume(ctx context.Context, queueName string, handler MessageHandler) error {
	// Setup all queues and exchanges
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

	log.Printf("Started consuming from queue: %s (with retry and DLQ)", queueName)

	// Process messages in goroutine
	go c.consumeLoop(ctx, msgs, queueName, handler)

	return nil
}

// consumeLoop handles the main message consumption loop
func (c *RabbitMQConsumer) consumeLoop(ctx context.Context, msgs <-chan amqp.Delivery, queueName string, handler MessageHandler) {
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
			c.handleMessage(ctx, msg, queueName, handler)
		}
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
