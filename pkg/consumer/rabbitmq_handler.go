package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"factorial-cal-services/pkg/domain"
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
	factorialService service.FactorialService
	redisService     service.RedisService
	s3Service        service.S3Service
	repository       domain.FactorialRepository
	maxRequestRepo   domain.MaxRequestRepository
	checksumService  service.ChecksumService
}

// NewFactorialMessageHandler creates a new factorial message handler
func NewFactorialMessageHandler(
	factorialService service.FactorialService,
	redisService service.RedisService,
	s3Service service.S3Service,
	repository domain.FactorialRepository,
	maxRequestRepo domain.MaxRequestRepository,
	checksumService service.ChecksumService,
) MessageHandler {
	handler := &FactorialMessageHandler{
		factorialService: factorialService,
		redisService:     redisService,
		s3Service:        s3Service,
		repository:       repository,
		maxRequestRepo:   maxRequestRepo,
		checksumService:  checksumService,
	}
	return handler.Handle
}

// NewFactorialBatchHandler creates a new batch handler for factorial calculations
func NewFactorialBatchHandler(
	factorialService service.FactorialService,
	redisService service.RedisService,
	s3Service service.S3Service,
	repository domain.FactorialRepository,
	maxRequestRepo domain.MaxRequestRepository,
	checksumService service.ChecksumService,
) BatchMessageHandler {
	handler := &FactorialMessageHandler{
		factorialService: factorialService,
		redisService:     redisService,
		s3Service:        s3Service,
		repository:       repository,
		maxRequestRepo:   maxRequestRepo,
		checksumService:  checksumService,
	}
	return handler.HandleBatch
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
	checksum := h.checksumService.Calculate(result)
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

// HandleBatch processes a batch of factorial calculation messages
func (h *FactorialMessageHandler) HandleBatch(ctx context.Context, payloads [][]byte) error {
	// Parse all messages
	numbers := make([]string, 0, len(payloads))
	for _, payload := range payloads {
		var message struct {
			Number string `json:"number"`
		}
		if err := json.Unmarshal(payload, &message); err != nil {
			log.Printf("Failed to parse message: %v", err)
			continue
		}
		numbers = append(numbers, message.Number)
	}

	if len(numbers) == 0 {
		return fmt.Errorf("no valid messages in batch")
	}

	log.Printf("Processing batch of %d factorial calculations", len(numbers))

	// Extract max value from batch
	maxNumber := h.findMaxNumber(numbers)
	log.Printf("Max number in batch: %s", maxNumber)

	// Update factorial_max_request_numbers if DB value < batch max
	if err := h.maxRequestRepo.SetMaxNumberIfGreater(maxNumber); err != nil {
		log.Printf("Warning: Failed to update max request number: %v", err)
		// Continue processing even if max update fails
	}

	// Process each number in the batch
	for _, number := range numbers {
		if err := h.processSingleNumber(ctx, number); err != nil {
			log.Printf("Error processing number %s: %v", number, err)
			// Continue with other numbers even if one fails
		}
	}

	log.Printf("Batch processing completed: %d numbers processed", len(numbers))
	return nil
}

// findMaxNumber finds the maximum number from a slice of number strings
func (h *FactorialMessageHandler) findMaxNumber(numbers []string) string {
	if len(numbers) == 0 {
		return "0"
	}

	maxNum := numbers[0]
	for _, num := range numbers[1:] {
		// Compare as integers for accurate comparison
		current, err1 := strconv.ParseInt(num, 10, 64)
		max, err2 := strconv.ParseInt(maxNum, 10, 64)
		
		if err1 != nil || err2 != nil {
			// If parsing fails, compare as strings (lexicographic)
			if len(num) > len(maxNum) || (len(num) == len(maxNum) && num > maxNum) {
				maxNum = num
			}
			continue
		}
		
		if current > max {
			maxNum = num
		}
	}
	return maxNum
}

// processSingleNumber processes a single factorial calculation
func (h *FactorialMessageHandler) processSingleNumber(ctx context.Context, number string) error {
	log.Printf("Processing factorial calculation for number: %s", number)

	// Check if already calculated
	existing, err := h.repository.FindByNumber(number)
	if err != nil {
		return fmt.Errorf("failed to check existing calculation: %w", err)
	}

	if existing != nil && existing.Status == domain.StatusDone {
		log.Printf("Calculation for %s already completed, skipping", number)
		return nil
	}

	// Create or update DB record with status=Calculating
	if existing == nil {
		calc := &domain.FactorialCalculation{
			Number: number,
			Status: domain.StatusCalculating,
			S3Key:  "", // Will be updated after upload
		}
		if err := h.repository.Create(calc); err != nil {
			return fmt.Errorf("failed to create calculation record: %w", err)
		}
	} else {
		if err := h.repository.UpdateStatus(number, domain.StatusCalculating); err != nil {
			return fmt.Errorf("failed to update status: %w", err)
		}
	}

	// Calculate factorial
	result, err := h.factorialService.CalculateFactorial(number)
	if err != nil {
		// Update status to failed
		h.repository.UpdateStatus(number, domain.StatusFailed)
		return fmt.Errorf("failed to calculate factorial: %w", err)
	}

	log.Printf("Factorial calculated for %s (result length: %d characters)", number, len(result))

	// Calculate checksum
	checksum := h.checksumService.Calculate(result)
	size := int64(len(result))

	// Upload to S3 (all results go to S3)
	s3Key, err := h.s3Service.UploadFactorial(ctx, number, result)
	if err != nil {
		h.repository.UpdateStatus(number, domain.StatusFailed)
		return fmt.Errorf("failed to upload to S3: %w", err)
	}

	log.Printf("Uploaded to S3: %s", s3Key)

	// Update DB with S3 key, checksum, size, and status=Done
	if err := h.repository.UpdateS3KeyWithChecksum(number, s3Key, checksum, size, domain.StatusDone); err != nil {
		return fmt.Errorf("failed to update S3 key: %w", err)
	}

	// Cache to Redis if applicable (only for small numbers)
	if h.redisService.ShouldCache(number) {
		if err := h.redisService.Set(ctx, number, result); err != nil {
			log.Printf("Warning: Failed to cache result in Redis: %v", err)
		} else {
			log.Printf("Cached result in Redis for %s", number)
		}
	} else {
		log.Printf("Skipping Redis cache for large number: %s", number)
	}

	log.Printf("Successfully processed factorial calculation for %s", number)
	return nil
}
