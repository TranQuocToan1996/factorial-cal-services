package consumer

import (
	"context"
	"fmt"
	"log"

	"factorial-cal-services/pkg/dto"
	"factorial-cal-services/pkg/repository"
	"factorial-cal-services/pkg/service"

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

// FactorialMessageHandler handles factorial calculation messages
type FactorialMessageHandler struct {
	factorialService      service.FactorialService
	redisService          service.RedisService
	storage               service.StorageService
	factRepository        repository.FactorialRepository
	maxRequestRepo        repository.MaxRequestRepository
	currentCalculatedRepo repository.CurrentCalculatedRepository
}

// NewFactorialMessageHandler creates a new factorial message handler
func NewFactorialMessageHandler(
	factorialService service.FactorialService,
	redisService service.RedisService,
	storage service.StorageService,
	repository repository.FactorialRepository,
	maxRequestRepo repository.MaxRequestRepository,
	currentCalculatedRepo repository.CurrentCalculatedRepository,
) MessageHandler {
	handler := &FactorialMessageHandler{
		factorialService:      factorialService,
		redisService:          redisService,
		storage:               storage,
		factRepository:        repository,
		maxRequestRepo:        maxRequestRepo,
		currentCalculatedRepo: currentCalculatedRepo,
	}
	return handler.HandleRequestCalculateFactorial
}

func (h *FactorialMessageHandler) HandleRequestCalculateFactorial(ctx context.Context, body []byte) error {
	var message dto.FactorialMessage

	if err := message.Unmarshal(body); err != nil {
		return fmt.Errorf("failed to parse message: %w", err)
	}

	log.Printf("Update max request number for number: %v", message.Number)

	// Check if already calculated
	rowAffected, err := h.maxRequestRepo.SetMaxNumberIfGreater(message.Number)
	if rowAffected == 0 || err != nil {
		return fmt.Errorf("max number %v is not greater than the current max number: %v - %v", message.Number, rowAffected, err)
	}

	return nil
}
