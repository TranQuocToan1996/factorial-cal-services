package dto

import "time"

// APIResponse represents the standard API response wrapper
type APIResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// CalculateRequest represents the request to calculate a factorial
type CalculateRequest struct {
	Number int64 `json:"number" binding:"required"`
}

// CalculateResponseData represents the data payload for calculate response
type CalculateResponseData struct {
	Number  int64  `json:"number,omitempty"`
	Message string `json:"message,omitempty"`
}

// ResultResponseData represents the data payload for result response
type ResultResponseData struct {
	Number          string `json:"number"`
	FactorialResult string `json:"factorial_result"`
}

// MetadataResponseData represents the data payload for metadata response
type MetadataResponseData struct {
	ID        int64     `json:"id"`
	Number    int64     `json:"number"`
	S3Key     string    `json:"s3_key,omitempty"`
	Checksum  string    `json:"checksum,omitempty"`
	Status    string    `json:"status"`
	Bucket    string    `json:"bucket"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Legacy DTOs for backward compatibility (deprecated, use APIResponse wrapper)
// CalculateResponse represents the response after submitting a factorial calculation
type CalculateResponse struct {
	Number string `json:"number"`
	Status string `json:"status"` // "accepted"
}

// ResultResponse represents the response containing the factorial result
type ResultResponse struct {
	Number string `json:"number"`
	Result string `json:"result"`
	Status string `json:"status"` // "done"
}

// MetadataResponse represents the metadata of a factorial calculation
type MetadataResponse struct {
	Number    string    `json:"number"`
	Status    string    `json:"status"`
	S3Key     string    `json:"s3_key,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
