package consumer

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"factorial-cal-services/pkg/utils/patterns"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQConsumer implements Consumer for RabbitMQ
type RabbitMQConsumer struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	semaphore *patterns.Semaphore
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
		conn:      conn,
		channel:   ch,
		semaphore: patterns.NewSemaphore(runtime.NumCPU() * 4),
	}, nil
}

// Consume starts consuming messages from the specified queue
func (c *RabbitMQConsumer) Consume(ctx context.Context, queueName string, handler MessageHandler) error {
	// Setup queue
	if err := c.setupQueues(queueName); err != nil {
		return err
	}
	autoAck := false          // false: Manual acknowledgment (ack after processing), true: Auto-ack on delivery
	exclusiveConsume := false // false: Multiple consumers can share the queue (load balancing), true: Only this consumer can consume from the queue
	noLocal := false          // false: Can receive messages published on the same connection, true: Prevents receiving messages published on the same connection
	nowait := false           // false: Wait for server confirmation, true: Fire-and-forget (no response)

	msgs, err := c.channel.Consume(
		queueName,
		"",
		autoAck,
		exclusiveConsume,
		noLocal,
		nowait,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	log.Printf("Started %v consuming from queue: %s", runtime.NumCPU(), queueName)
	for msg := range msgs {
		c.semaphore.Submit(func() error {
			c.handleMessageAndAck(ctx, msg, handler)
			return nil
		})
	}

	return nil
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
