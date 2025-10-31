package producer

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

// RabbitMQProducer implements Producer for RabbitMQ
type RabbitMQProducer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewRabbitMQProducer creates a new RabbitMQ producer
func NewRabbitMQProducer(amqpURL string) (*RabbitMQProducer, error) {
	var conn *amqp.Connection
	var err error

	// Check if TLS is enabled by looking for RABBITMQ_CA
	rootCA := os.Getenv("RABBITMQ_CA")
	if rootCA != "" {
		// Replace literal \n with actual newlines
		rootCA = strings.ReplaceAll(rootCA, "\\n", "\n")

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

	fmt.Println("Connected to RabbitMQ successfully!")

	return &RabbitMQProducer{
		conn:    conn,
		channel: ch,
	}, nil
}

// Publish sends a message to the specified queue
func (p *RabbitMQProducer) Publish(ctx context.Context, queueName string, payload []byte) error {
	// Publish message
	err := p.channel.PublishWithContext(
		ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         payload,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Published message to queue %s: %s", queueName, string(payload))
	return nil
}

// Close closes the RabbitMQ connection and channel
func (p *RabbitMQProducer) Close() error {
	if p.channel != nil {
		if err := p.channel.Close(); err != nil {
			return err
		}
	}
	if p.conn != nil {
		if err := p.conn.Close(); err != nil {
			return err
		}
	}
	return nil
}

// IsHealthy checks if the RabbitMQ connection is alive
func (p *RabbitMQProducer) IsHealthy() bool {
	return p.conn != nil && !p.conn.IsClosed() && p.channel != nil
}
