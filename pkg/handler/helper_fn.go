package handler

import (
	"factorial-cal-services/pkg/dto"

	"github.com/gin-gonic/gin"
)

// sendAPIResponse sends a standardized API response
func sendAPIResponse(c *gin.Context, code int, status string, message string, data any) {
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
