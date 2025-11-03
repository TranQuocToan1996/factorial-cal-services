package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"factorial-cal-services/pkg/domain"
	"factorial-cal-services/pkg/dto"
	"factorial-cal-services/pkg/producer"
	"factorial-cal-services/pkg/service"

	"github.com/gin-gonic/gin"
)

// sendAPIResponse sends a standardized API response
func sendAPIResponse(c *gin.Context, code int, status string, message string, data interface{}) {
	c.JSON(code, dto.APIResponse{
		Code:    code,
		Status:  status,
		Message: message,
		Data:    data,
	})
}

// sendErrorResponse sends an error response in the new format
func sendErrorResponse(c *gin.Context, code int, status string, message string) {
	sendAPIResponse(c, code, status, message, dto.ErrorResponse{
		Error:   status,
		Message: message,
	})
}

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
		sendErrorResponse(c, http.StatusBadRequest, "fail", err.Error())
		return
	}

	// Validate number
	_, err := h.factorialService.ValidateNumber(req.Number)
	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "fail", err.Error())
		return
	}

	// Check if already calculated
	existing, err := h.repository.FindByNumber(req.Number)
	if err != nil {
		log.Printf("Error checking existing calculation: %v", err)
		sendErrorResponse(c, http.StatusInternalServerError, "fail", "Failed to check calculation status")
		return
	}

	if existing != nil && existing.Status == domain.StatusDone {
		// Return result if already calculated
		ctx := context.Background()
		var factorialResult string
		if h.redisService.ShouldCache(req.Number) {
			factorialResult, _ = h.redisService.Get(ctx, req.Number)
		}
		if factorialResult == "" {
			factorialResult, _ = h.s3Service.DownloadFactorial(ctx, existing.S3Key)
		}

		sendAPIResponse(c, http.StatusOK, "ok", "done", dto.ResultResponseData{
			Number:          req.Number,
			FactorialResult: factorialResult,
		})
		return
	}

	// Publish to RabbitMQ
	message := map[string]string{"number": req.Number}
	messageBytes, _ := json.Marshal(message)

	err = h.producer.Publish(context.Background(), h.queueName, messageBytes)
	if err != nil {
		log.Printf("Error publishing message: %v", err)
		sendErrorResponse(c, http.StatusInternalServerError, "fail", "Failed to submit calculation")
		return
	}

	// Return calculating status
	sendAPIResponse(c, http.StatusOK, "ok", "calculating", dto.CalculateResponseData{})
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
		sendErrorResponse(c, http.StatusBadRequest, "fail", err.Error())
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
			sendAPIResponse(c, http.StatusOK, "ok", "done", dto.ResultResponseData{
				Number:          number,
				FactorialResult: result,
			})
			return
		}
	}

	// Not in cache or large number, check DB for S3 key
	calc, err := h.repository.FindByNumber(number)
	if err != nil {
		log.Printf("Error finding calculation: %v", err)
		sendErrorResponse(c, http.StatusInternalServerError, "fail", "Failed to retrieve calculation")
		return
	}

	if calc == nil {
		// Not found - publish event and return calculating status (per flow.md)
		message := map[string]string{"number": number}
		messageBytes, _ := json.Marshal(message)
		if err := h.producer.Publish(ctx, h.queueName, messageBytes); err != nil {
			log.Printf("Error publishing event: %v", err)
		}
		sendAPIResponse(c, http.StatusOK, "ok", "calculating", dto.CalculateResponseData{})
		return
	}

	if calc.Status != domain.StatusDone {
		sendAPIResponse(c, http.StatusOK, "ok", "calculating", dto.CalculateResponseData{})
		return
	}

	// Download from S3
	result, err := h.s3Service.DownloadFactorial(ctx, calc.S3Key)
	if err != nil {
		log.Printf("Error downloading from S3: %v", err)
		sendErrorResponse(c, http.StatusInternalServerError, "fail", "Failed to retrieve result from storage")
		return
	}

	// Cache to Redis if applicable
	if h.redisService.ShouldCache(number) {
		if err := h.redisService.Set(ctx, number, result); err != nil {
			log.Printf("Failed to cache result: %v", err)
		}
	}

	sendAPIResponse(c, http.StatusOK, "ok", "done", dto.ResultResponseData{
		Number:          number,
		FactorialResult: result,
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
		sendErrorResponse(c, http.StatusBadRequest, "fail", err.Error())
		return
	}

	// Query DB
	calc, err := h.repository.FindByNumber(number)
	if err != nil {
		log.Printf("Error finding calculation: %v", err)
		sendErrorResponse(c, http.StatusInternalServerError, "fail", "Failed to retrieve metadata")
		return
	}

	if calc == nil {
		sendAPIResponse(c, http.StatusOK, "ok", "calculating", dto.CalculateResponseData{})
		return
	}

	sendAPIResponse(c, http.StatusOK, "ok", "done", dto.MetadataResponseData{
		ID:        strconv.FormatInt(calc.ID, 10),
		Number:    calc.Number,
		S3Key:     calc.S3Key,
		Checksum:  calc.Checksum,
		Status:    calc.Status,
		CreatedAt: calc.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: calc.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}
