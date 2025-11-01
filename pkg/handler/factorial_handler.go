package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"factorial-cal-services/pkg/domain"
	"factorial-cal-services/pkg/dto"
	"factorial-cal-services/pkg/producer"
	"factorial-cal-services/pkg/service"

	"github.com/gin-gonic/gin"
)

// FactorialHandler handles factorial calculation HTTP requests
type FactorialHandler struct {
	factorialService service.FactorialService
	redisService     service.RedisService
	s3Service        service.S3Service
	repository       domain.FactorialRepository
	producer         producer.Producer
	queueName        string
}

// NewFactorialHandler creates a new factorial handler
func NewFactorialHandler(
	factorialService service.FactorialService,
	redisService service.RedisService,
	s3Service service.S3Service,
	repository domain.FactorialRepository,
	producer producer.Producer,
	queueName string,
) *FactorialHandler {
	return &FactorialHandler{
		factorialService: factorialService,
		redisService:     redisService,
		s3Service:        s3Service,
		repository:       repository,
		producer:         producer,
		queueName:        queueName,
	}
}

// SubmitCalculation godoc
// @Summary      Submit factorial calculation
// @Description  Submit a number for factorial calculation (async processing)
// @Tags         factorial
// @Accept       json
// @Produce      json
// @Param        request body dto.CalculateRequest true "Calculation Request"
// @Success      202  {object}  dto.CalculateResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /factorial [post]
func (h *FactorialHandler) SubmitCalculation(c *gin.Context) {
	var req dto.CalculateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	// Validate number
	_, err := h.factorialService.ValidateNumber(req.Number)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_number",
			Message: err.Error(),
		})
		return
	}

	// Check if already calculated
	existing, err := h.repository.FindByNumber(req.Number)
	if err != nil {
		log.Printf("Error checking existing calculation: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to check calculation status",
		})
		return
	}

	if existing != nil && existing.Status == domain.StatusDone {
		c.JSON(http.StatusOK, dto.CalculateResponse{
			Number: req.Number,
			Status: "already_calculated",
		})
		return
	}

	// Publish to RabbitMQ
	message := map[string]string{"number": req.Number}
	messageBytes, _ := json.Marshal(message)

	err = h.producer.Publish(context.Background(), h.queueName, messageBytes)
	if err != nil {
		log.Printf("Error publishing message: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to submit calculation",
		})
		return
	}

	c.JSON(http.StatusAccepted, dto.CalculateResponse{
		Number: req.Number,
		Status: "accepted",
	})
}

// GetResult godoc
// @Summary      Get factorial result
// @Description  Get the factorial calculation result for a number
// @Tags         factorial
// @Produce      json
// @Param        number path string true "Number"
// @Success      200  {object}  dto.ResultResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /factorial/{number} [get]
func (h *FactorialHandler) GetResult(c *gin.Context) {
	number := c.Param("number")

	// Validate number format
	_, err := h.factorialService.ValidateNumber(number)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_number",
			Message: err.Error(),
		})
		return
	}

	ctx := context.Background()

	// Check if should use Redis cache
	if h.redisService.ShouldCache(number) {
		// Try Redis first
		result, err := h.redisService.Get(ctx, number)
		if err != nil {
			log.Printf("Redis error: %v", err)
		} else if result != "" {
			c.JSON(http.StatusOK, dto.ResultResponse{
				Number: number,
				Result: result,
				Status: "done",
			})
			return
		}
	}

	// Not in cache or large number, check DB for S3 key
	calc, err := h.repository.FindByNumber(number)
	if err != nil {
		log.Printf("Error finding calculation: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to retrieve calculation",
		})
		return
	}

	if calc == nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "not_found",
			Message: "Calculation not found or not yet completed",
		})
		return
	}

	if calc.Status != domain.StatusDone {
		c.JSON(http.StatusAccepted, dto.ErrorResponse{
			Error:   "processing",
			Message: fmt.Sprintf("Calculation is still in progress (status: %s)", calc.Status),
		})
		return
	}

	// Download from S3
	result, err := h.s3Service.DownloadFactorial(ctx, calc.S3Key)
	if err != nil {
		log.Printf("Error downloading from S3: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to retrieve result from storage",
		})
		return
	}

	// Cache to Redis if applicable
	if h.redisService.ShouldCache(number) {
		if err := h.redisService.Set(ctx, number, result); err != nil {
			log.Printf("Failed to cache result: %v", err)
		}
	}

	c.JSON(http.StatusOK, dto.ResultResponse{
		Number: number,
		Result: result,
		Status: "done",
	})
}

// GetMetadata godoc
// @Summary      Get factorial calculation metadata
// @Description  Get the metadata of a factorial calculation (status, S3 key, etc.)
// @Tags         factorial
// @Produce      json
// @Param        number path string true "Number"
// @Success      200  {object}  dto.MetadataResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /factorial/metadata/{number} [get]
func (h *FactorialHandler) GetMetadata(c *gin.Context) {
	number := c.Param("number")

	// Validate number format
	_, err := h.factorialService.ValidateNumber(number)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_number",
			Message: err.Error(),
		})
		return
	}

	// Query DB
	calc, err := h.repository.FindByNumber(number)
	if err != nil {
		log.Printf("Error finding calculation: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to retrieve metadata",
		})
		return
	}

	if calc == nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "not_found",
			Message: "Calculation not found",
		})
		return
	}

	c.JSON(http.StatusOK, dto.MetadataResponse{
		Number:    calc.Number,
		Status:    calc.Status,
		S3Key:     calc.S3Key,
		CreatedAt: calc.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

