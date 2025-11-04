package consumer

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"time"

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
) MessageHandler {
	handler := &FactorialMessageHandler{
		factorialService: factorialService,
		redisService:     redisService,
		s3Service:        s3Service,
		repository:       repository,
		maxRequestRepo:   maxRequestRepo,
	}
	return handler.Handle
}

// NewFactorialBatchHandler creates a new batch handler for factorial calculations
func NewFactorialBatchHandler(
	factorialService service.FactorialService,
	redisService service.RedisService,
	s3Service service.S3Service,
	repository repository.FactorialRepository,
	maxRequestRepo repository.MaxRequestRepository,
	currentCalculatedRepo repository.CurrentCalculatedRepository,
	incrementalService service.IncrementalFactorialService,
) BatchMessageHandler {
	handler := &FactorialMessageHandler{
		factorialService:      factorialService,
		redisService:          redisService,
		s3Service:             s3Service,
		repository:            repository,
		maxRequestRepo:        maxRequestRepo,
		currentCalculatedRepo: currentCalculatedRepo,
		incrementalService:    incrementalService,
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
	checksum := h.Calculate(result)
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

	// Get current calculated number from DB
	currentNumber, err := h.currentCalculatedRepo.GetCurrentNumber()
	if err != nil {
		log.Printf("Warning: Failed to get current number: %v", err)
		currentNumber = "0" // Default to 0 if error
	}

	log.Printf("Current calculated number: %s, Target max number: %s", currentNumber, maxNumber)

	// Process incrementally from current_number to max_number
	if err := h.processIncrementalBatch(ctx, currentNumber, maxNumber, numbers); err != nil {
		log.Printf("Error in incremental batch processing: %v", err)
		return err
	}

	log.Printf("Batch processing completed: %d numbers processed", len(numbers))
	return nil
}

// findMaxNumber finds the maximum number from a slice of number strings
func (h *FactorialMessageHandler) findMaxNumber(numbers []string) string {
	if len(numbers) == 0 {
		return "0"
	}

	var maxNum string
	var maxValue int64 = -1
	hasValidNumber := false

	for _, num := range numbers {
		// Parse as integer for accurate comparison
		value, err := strconv.ParseInt(num, 10, 64)
		if err != nil {
			// Skip invalid strings, only consider valid numbers
			continue
		}

		if !hasValidNumber || value > maxValue {
			maxValue = value
			maxNum = num
			hasValidNumber = true
		}
	}

	// If no valid numbers found, return first string (fallback to lexicographic)
	if !hasValidNumber {
		return numbers[0]
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
	checksum := h.Calculate(result)
	size := int64(len(result))

	// Upload to S3 (all results go to S3)
	s3Ctx, s3Cancel := context.WithTimeout(ctx, 30*time.Second)
	s3Key, err := h.s3Service.UploadFactorial(s3Ctx, number, result)
	s3Cancel()
	if err != nil {
		if err := h.repository.UpdateStatus(number, domain.StatusFailed); err != nil {
			log.Printf("Error updating status to failed: %v", err)
		}
		return fmt.Errorf("failed to upload to S3: %w", err)
	}

	log.Printf("Uploaded to S3: %s", s3Key)

	// Update DB with S3 key, checksum, size, and status=Done
	if err := h.repository.UpdateS3KeyWithChecksum(number, s3Key, checksum, size, domain.StatusDone); err != nil {
		return fmt.Errorf("failed to update S3 key: %w", err)
	}

	// Cache to Redis if applicable (only for small numbers)
	if h.redisService.ShouldCache(number) {
		redisCtx, redisCancel := context.WithTimeout(ctx, 5*time.Second)
		if err := h.redisService.Set(redisCtx, number, result); err != nil {
			log.Printf("Warning: Failed to cache result in Redis: %v", err)
		} else {
			log.Printf("Cached result in Redis for %s", number)
		}
		redisCancel()
	} else {
		log.Printf("Skipping Redis cache for large number: %s", number)
	}

	log.Printf("Successfully processed factorial calculation for %s", number)
	return nil
}

// processIncrementalBatch calculates factorials incrementally from current_number to max_number
// This implements the Step Functions logic in the worker (bottom-up calculation)
func (h *FactorialMessageHandler) processIncrementalBatch(ctx context.Context, currentNumber string, maxNumber string, requestedNumbers []string) error {
	curNum, err := strconv.ParseInt(currentNumber, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid current number: %w", err)
	}

	maxNum, err := strconv.ParseInt(maxNumber, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid max number: %w", err)
	}

	// If current >= max, nothing to calculate
	if curNum >= maxNum {
		log.Printf("Current number %d >= max number %d, nothing to calculate", curNum, maxNum)
		// Still process requested numbers individually in case they're not done
		for _, number := range requestedNumbers {
			if err := h.processSingleNumber(ctx, number); err != nil {
				log.Printf("Error processing number %s: %v", number, err)
			}
		}
		return nil
	}

	// Get current factorial from Redis or S3
	var curFactorial string
	if curNum > 0 {
		curNumberStr := strconv.FormatInt(curNum, 10)
		redisCtx, redisCancel := context.WithTimeout(ctx, 5*time.Second)
		if h.redisService.ShouldCache(curNumberStr) {
			var err error
			curFactorial, err = h.redisService.Get(redisCtx, curNumberStr)
			if err != nil {
				log.Printf("Warning: Failed to get current factorial from Redis: %v", err)
			}
		}
		redisCancel()
		if curFactorial == "" {
			// Try S3
			calc, err := h.repository.FindByNumber(curNumberStr)
			if err == nil && calc != nil && calc.Status == domain.StatusDone {
				s3Ctx, s3Cancel := context.WithTimeout(ctx, 30*time.Second)
				var err error
				curFactorial, err = h.s3Service.DownloadFactorial(s3Ctx, calc.S3Key)
				s3Cancel()
				if err != nil {
					log.Printf("Warning: Failed to get current factorial from S3: %v", err)
				}
			}
		}
	}

	// If no current factorial found and curNum > 0, calculate from scratch
	if curFactorial == "" && curNum > 0 {
		log.Printf("Current factorial not found for %d, calculating from scratch", curNum)
		curFactorial, err = h.factorialService.CalculateFactorial(strconv.FormatInt(curNum, 10))
		if err != nil {
			return fmt.Errorf("failed to calculate current factorial: %w", err)
		}
	}

	// If still no factorial, start from 1
	if curFactorial == "" {
		curFactorial = "1"
		curNum = 0
	}

	log.Printf("Starting incremental calculation from %d to %d", curNum+1, maxNum)

	// Calculate incrementally: next_fac = (cur_number + 1) * cur_fac
	result := new(big.Int)
	result, ok := result.SetString(curFactorial, 10)
	if !ok {
		return fmt.Errorf("invalid current factorial format")
	}

	// Process each number incrementally from current+1 to max
	for i := curNum + 1; i <= maxNum; i++ {
		numberStr := strconv.FormatInt(i, 10)

		// Check if already calculated
		existing, err := h.repository.FindByNumber(numberStr)
		if err != nil {
			log.Printf("Error checking existing calculation for %s: %v", numberStr, err)
			continue
		}
		if existing != nil && existing.Status == domain.StatusDone {
			log.Printf("Number %s already calculated, skipping", numberStr)
			// Still update current factorial for next iteration
			multiplier := big.NewInt(i)
			result.Mul(result, multiplier)
			continue
		}

		// Create or update DB record with status=Calculating
		if existing == nil {
			calc := &domain.FactorialCalculation{
				Number: numberStr,
				Status: domain.StatusCalculating,
				S3Key:  "",
			}
			if err := h.repository.Create(calc); err != nil {
				log.Printf("Error creating calculation record for %s: %v", numberStr, err)
				continue
			}
		} else {
			if err := h.repository.UpdateStatus(numberStr, domain.StatusCalculating); err != nil {
				log.Printf("Error updating status for %s: %v", numberStr, err)
				continue
			}
		}

		// Calculate next factorial: next_fac = (cur_number + 1) * cur_fac
		multiplier := big.NewInt(i)
		result.Mul(result, multiplier)
		factorialResult := result.String()

		log.Printf("Calculated factorial for %d (result length: %d characters)", i, len(factorialResult))

		// Calculate checksum
		checksum := h.Calculate(factorialResult)
		size := int64(len(factorialResult))

		// Upload to S3
		s3Ctx, s3Cancel := context.WithTimeout(ctx, 30*time.Second)
		s3Key, err := h.s3Service.UploadFactorial(s3Ctx, numberStr, factorialResult)
		s3Cancel()
		if err != nil {
			log.Printf("Error uploading to S3 for %s: %v", numberStr, err)
			if err := h.repository.UpdateStatus(numberStr, domain.StatusFailed); err != nil {
				log.Printf("Error updating status to failed for %s: %v", numberStr, err)
			}
			continue
		}

		// Cache to Redis if applicable
		if h.redisService.ShouldCache(numberStr) {
			redisCtx, redisCancel := context.WithTimeout(ctx, 5*time.Second)
			if err := h.redisService.Set(redisCtx, numberStr, factorialResult); err != nil {
				log.Printf("Warning: Failed to cache result in Redis for %s: %v", numberStr, err)
			}
			redisCancel()
		}

		// Update metadata and current_number atomically
		if err := h.repository.UpdateWithCurrentNumber(numberStr, s3Key, checksum, size, domain.StatusDone, numberStr); err != nil {
			log.Printf("Error updating metadata for %s: %v", numberStr, err)
			continue
		}

		log.Printf("Successfully processed factorial calculation for %d", i)
	}

	log.Printf("Incremental batch processing completed: calculated from %d to %d", curNum+1, maxNum)
	return nil
}

// Calculate calculates SHA256 checksum for the given data
func (s *FactorialMessageHandler) Calculate(data string) string {
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}
