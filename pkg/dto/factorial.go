package dto

// CalculateRequest represents the request to calculate a factorial
type CalculateRequest struct {
	Number string `json:"number" binding:"required"`
}

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
	Number    string `json:"number"`
	Status    string `json:"status"`
	S3Key     string `json:"s3_key,omitempty"`
	CreatedAt string `json:"created_at"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

