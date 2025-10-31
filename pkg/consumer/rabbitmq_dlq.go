package consumer

import (
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// sendToDLQ publishes message to final DLQ with failure reason
func (c *RabbitMQConsumer) sendToDLQ(dlqName string, msg amqp.Delivery, failureErr error) error {
	failureReason := "unknown"
	if failureErr != nil {
		failureReason = failureErr.Error()
	}

	return c.channel.Publish(
		"",      // exchange (default)
		dlqName, // routing key (DLQ name)
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType:  msg.ContentType,
			Body:         msg.Body,
			DeliveryMode: amqp.Persistent,
			Headers: amqp.Table{
				"x-original-queue": msg.RoutingKey,
				"x-failed-at":      time.Now().Unix(),
				"x-retry-count":    c.getRetryCount(msg.Headers),
				"x-failure-reason": failureReason,
			},
		},
	)
}

// getRetryCount extracts retry count from x-death header
func (c *RabbitMQConsumer) getRetryCount(headers amqp.Table) int {
	if headers == nil {
		return 0
	}

	// RabbitMQ adds x-death header when message goes through DLX
	if xDeath, ok := headers["x-death"].([]any); ok && len(xDeath) > 0 {
		if death, ok := xDeath[0].(amqp.Table); ok {
			if count, ok := death["count"].(int64); ok {
				return int(count)
			}
		}
	}

	return 0
}
