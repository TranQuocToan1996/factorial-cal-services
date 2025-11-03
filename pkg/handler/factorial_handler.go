package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

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
// @Success      200  {object}  dto.APIResponse{data=dto.ResultResponseData} "Calculation already completed, returns result"
// @Success      200  {object}  dto.APIResponse{data=dto.CalculateResponseData} "Calculation submitted successfully"
// @Failure      400  {object}  dto.APIResponse{data=dto.ErrorResponse} "Invalid request - number missing or invalid format"
// @Failure      500  {object}  dto.APIResponse{data=dto.ErrorResponse} "Internal server error - database or queue failure"
// @Router       /factorial [post]
// @Example      400 {"code":400,"status":"fail","message":"invalid number format","data":{"error":"fail","message":"invalid number format"}}
// @Example      500 {"code":500,"status":"fail","message":"Failed to submit calculation","data":{"error":"fail","message":"Failed to submit calculation"}}
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
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		var factorialResult string
		var err error
		if h.redisService.ShouldCache(req.Number) {
			redisCtx, redisCancel := context.WithTimeout(ctx, 5*time.Second)
			factorialResult, err = h.redisService.Get(redisCtx, req.Number)
			redisCancel()
			if err != nil {
				log.Printf("Redis error: %v", err)
			}
		}
		if factorialResult == "" && h.s3Service != nil {
			factorialResult, err = h.s3Service.DownloadFactorial(ctx, existing.S3Key)
			if err != nil {
				log.Printf("S3 download error: %v", err)
				sendErrorResponse(c, http.StatusInternalServerError, "fail", "Failed to retrieve result")
				return
			}
		}

		sendAPIResponse(c, http.StatusOK, "ok", "done", dto.ResultResponseData{
			Number:          req.Number,
			FactorialResult: factorialResult,
		})
		return
	}

	// Publish to RabbitMQ
	message := map[string]string{"number": req.Number}
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		sendErrorResponse(c, http.StatusInternalServerError, "fail", "Failed to prepare calculation request")
		return
	}

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
// @Success      200  {object}  dto.APIResponse{data=dto.ResultResponseData} "Result retrieved successfully"
// @Success      200  {object}  dto.APIResponse{data=dto.CalculateResponseData} "Calculation in progress"
// @Failure      400  {object}  dto.APIResponse{data=dto.ErrorResponse} "Invalid number format"
// @Failure      500  {object}  dto.APIResponse{data=dto.ErrorResponse} "Internal server error - database or storage failure"
// @Router       /factorial/{number} [get]
// @Example      200 {"code":200,"status":"ok","message":"done","data":{"number":"10","factorial_result":"3628800"}}
// @Example      200 {"code":200,"status":"ok","message":"calculating","data":{}}
// @Example      400 {"code":400,"status":"fail","message":"number must be between 0 and 10000","data":{"error":"fail","message":"number must be between 0 and 10000"}}
// @Example      500 {"code":500,"status":"fail","message":"Failed to retrieve result from storage","data":{"error":"fail","message":"Failed to retrieve result from storage"}}
func (h *FactorialHandler) GetResult(c *gin.Context) {
	number := c.Param("number")

	// Validate number format
	_, err := h.factorialService.ValidateNumber(number)
	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "fail", err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Check if S3 service is available
	if h.s3Service == nil {
		sendErrorResponse(c, http.StatusInternalServerError, "fail", "S3 service not configured")
		return
	}

	// Check if should use Redis cache
	if h.redisService.ShouldCache(number) {
		// Try Redis first
		redisCtx, redisCancel := context.WithTimeout(ctx, 5*time.Second)
		result, err := h.redisService.Get(redisCtx, number)
		redisCancel()
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
		messageBytes, err := json.Marshal(message)
		if err != nil {
			log.Printf("Error marshaling message: %v", err)
			sendErrorResponse(c, http.StatusInternalServerError, "fail", "Failed to prepare calculation request")
			return
		}
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
	s3Ctx, s3Cancel := context.WithTimeout(ctx, 30*time.Second)
	defer s3Cancel()
	result, err := h.s3Service.DownloadFactorial(s3Ctx, calc.S3Key)
	if err != nil {
		log.Printf("Error downloading from S3: %v", err)
		sendErrorResponse(c, http.StatusInternalServerError, "fail", "Failed to retrieve result from storage")
		return
	}

	// Cache to Redis if applicable
	if h.redisService.ShouldCache(number) {
		redisCtx, redisCancel := context.WithTimeout(ctx, 5*time.Second)
		if err := h.redisService.Set(redisCtx, number, result); err != nil {
			log.Printf("Failed to cache result: %v", err)
		}
		redisCancel()
	}

	sendAPIResponse(c, http.StatusOK, "ok", "done", dto.ResultResponseData{
		Number:          number,
		FactorialResult: result,
	})
}

// GetMetadata godoc
// @Summary      Get factorial calculation metadata
// @Description  Get the metadata of a factorial calculation (status, S3 key, checksum, etc.)
// @Tags         factorial
// @Produce      json
// @Param        number path string true "Number"
// @Success      200  {object}  dto.APIResponse{data=dto.MetadataResponseData} "Metadata retrieved successfully"
// @Success      200  {object}  dto.APIResponse{data=dto.CalculateResponseData} "Calculation not found or in progress"
// @Failure      400  {object}  dto.APIResponse{data=dto.ErrorResponse} "Invalid number format"
// @Failure      500  {object}  dto.APIResponse{data=dto.ErrorResponse} "Internal server error - database failure"
// @Router       /factorial/metadata/{number} [get]
// @Example      200 {"code":200,"status":"ok","message":"done","data":{"id":"1","number":"10","factorial_result":"3628800","s3_key":"factorials/10.txt","checksum":"abc123...","status":"done","created_at":"2025-01-01T00:00:00Z","updated_at":"2025-01-01T00:00:00Z"}}
// @Example      200 {"code":200,"status":"ok","message":"calculating","data":{}}
// @Example      400 {"code":400,"status":"fail","message":"invalid number format","data":{"error":"fail","message":"invalid number format"}}
// @Example      500 {"code":500,"status":"fail","message":"Failed to retrieve metadata","data":{"error":"fail","message":"Failed to retrieve metadata"}}
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

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var factorialResult string
	if calc.Status == domain.StatusDone && h.s3Service != nil {
		// Download factorial result if status is done
		var err error
		factorialResult, err = h.s3Service.DownloadFactorial(ctx, calc.S3Key)
		if err != nil {
			log.Printf("Error downloading factorial result: %v", err)
			// Continue without factorial result if download fails
		}
	}

	sendAPIResponse(c, http.StatusOK, "ok", "done", dto.MetadataResponseData{
		ID:              strconv.FormatInt(calc.ID, 10),
		Number:          calc.Number,
		FactorialResult: factorialResult,
		S3Key:           calc.S3Key,
		Checksum:        calc.Checksum,
		Status:          calc.Status,
		CreatedAt:       calc.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:       calc.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}
