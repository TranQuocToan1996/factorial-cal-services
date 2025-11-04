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
	"sync"

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
	autoAck := false
	exclusiveConsume := true
	nowait := false

	// Start consuming messages
	msgs, err := c.channel.Consume(
		queueName,        // queue
		"",               // consumer tag
		autoAck,          // auto-ack = false (manual ack)
		exclusiveConsume, // exclusive
		false,            // no-local
		nowait,           // no-wait
		nil,              // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	log.Printf("Started %v consuming from queue: %s", runtime.NumCPU(), queueName)
	for i := 0; i < runtime.NumCPU(); i++ {
		go c.consumeLoop(ctx, msgs, handler)
	}

	return nil
}

// consumeLoop handles the main message consumption loop
func (c *RabbitMQConsumer) consumeLoop(ctx context.Context, msgs <-chan amqp.Delivery, handler MessageHandler) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, 16)
	for msg := range msgs {
		sem <- struct{}{}
		wg.Add(1)
		go func(msg amqp.Delivery) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("Recovered from panic in message handler: %v\n", r)
				}
			}()
			defer func() { <-sem }()
			c.handleMessageAndAck(ctx, msg, handler)
		}(msg)
	}
	wg.Wait()
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

func (c *RabbitMQConsumer) IsHealthy() bool {
	return false
}
