package consumer

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"

	"factorial-cal-services/pkg/domain"
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
	s3Service             service.S3Service
	repository            repository.FactorialRepository
	maxRequestRepo        repository.MaxRequestRepository
	currentCalculatedRepo repository.CurrentCalculatedRepository
	incrementalService    service.IncrementalFactorialService
}

// NewFactorialMessageHandler creates a new factorial message handler
func NewFactorialMessageHandler(
	factorialService service.FactorialService,
	redisService service.RedisService,
	s3Service service.S3Service,
	repository repository.FactorialRepository,
	maxRequestRepo repository.MaxRequestRepository,
	currentCalculatedRepo repository.CurrentCalculatedRepository,
	incrementalService service.IncrementalFactorialService,
) MessageHandler {
	handler := &FactorialMessageHandler{
		factorialService:      factorialService,
		redisService:          redisService,
		s3Service:             s3Service,
		repository:            repository,
		maxRequestRepo:        maxRequestRepo,
		currentCalculatedRepo: currentCalculatedRepo,
		incrementalService:    incrementalService,
	}
	return handler.Handle
}

// Handle processes a factorial calculation message
func (h *FactorialMessageHandler) Handle(ctx context.Context, body []byte) error {
	// Parse message
	var message struct {
		Number string `json:"number"`
	}

	if err := json.Unmarshal(body, &message); err != nil {
		return fmt.Errorf("failed to parse message: %w", err)
	}

	log.Printf("Processing factorial calculation for number: %s", message.Number)

	// Check if already calculated
	existing, err := h.repository.FindByNumber(message.Number)
	if err != nil {
		return fmt.Errorf("failed to check existing calculation: %w", err)
	}

	if existing != nil && existing.Status == domain.StatusDone {
		log.Printf("Calculation for %s already completed, skipping", message.Number)
		return nil
	}

	// Create or update DB record with status=Calculating
	if existing == nil {
		calc := &domain.FactorialCalculation{
			Number: message.Number,
			Status: domain.StatusCalculating,
			S3Key:  "", // Will be updated after upload
		}
		if err := h.repository.Create(calc); err != nil {
			return fmt.Errorf("failed to create calculation record: %w", err)
		}
	} else {
		if err := h.repository.UpdateStatus(message.Number, domain.StatusCalculating); err != nil {
			return fmt.Errorf("failed to update status: %w", err)
		}
	}

	// Calculate factorial
	result, err := h.factorialService.CalculateFactorial(message.Number)
	if err != nil {
		// Update status to failed
		h.repository.UpdateStatus(message.Number, domain.StatusFailed)
		return fmt.Errorf("failed to calculate factorial: %w", err)
	}

	log.Printf("Factorial calculated for %s (result length: %d characters)", message.Number, len(result))

	// Calculate checksum
	checksum := h.Checksum(result)
	size := int64(len(result))

	// Upload to S3 (all results go to S3)
	s3Key, err := h.s3Service.UploadFactorial(ctx, message.Number, result)
	if err != nil {
		h.repository.UpdateStatus(message.Number, domain.StatusFailed)
		return fmt.Errorf("failed to upload to S3: %w", err)
	}

	log.Printf("Uploaded to S3: %s", s3Key)

	// Update DB with S3 key, checksum, size, and status=Done
	if err := h.repository.UpdateS3KeyWithChecksum(message.Number, s3Key, checksum, size, domain.StatusDone); err != nil {
		return fmt.Errorf("failed to update S3 key: %w", err)
	}

	// Cache to Redis if applicable (only for small numbers)
	if h.redisService.ShouldCache(message.Number) {
		if err := h.redisService.Set(ctx, message.Number, result); err != nil {
			log.Printf("Warning: Failed to cache result in Redis: %v", err)
		} else {
			log.Printf("Cached result in Redis for %s", message.Number)
		}
	} else {
		log.Printf("Skipping Redis cache for large number: %s", message.Number)
	}

	log.Printf("Successfully processed factorial calculation for %s", message.Number)
	return nil
}

func (h *FactorialMessageHandler) ConsumeTopic(ctx context.Context) error {
	return nil
}

// Checksum calculates SHA256 checksum for the given data
func (s *FactorialMessageHandler) Checksum(data string) string {
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}
