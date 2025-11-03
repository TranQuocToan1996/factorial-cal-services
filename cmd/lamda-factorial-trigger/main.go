package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

// Event represents the input from Step Functions
type Event struct {
	Number string `json:"number"`
}

// Response represents the Lambda response
type Response struct {
	StatusCode int               `json:"statusCode"`
	Body       string            `json:"body"`
	Headers    map[string]string `json:"headers"`
}

// HandleRequest handles the Lambda invocation
func HandleRequest(ctx context.Context, event Event) (Response, error) {
	apiEndpoint := os.Getenv("API_ENDPOINT")
	if apiEndpoint == "" {
		return Response{
			StatusCode: 500,
			Body:       `{"error": "API_ENDPOINT not configured"}`,
		}, fmt.Errorf("API_ENDPOINT environment variable not set")
	}

	// Prepare request payload
	payload := map[string]string{
		"number": event.Number,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return Response{
			StatusCode: 500,
			Body:       fmt.Sprintf(`{"error": "Failed to marshal payload: %v"}`, err),
		}, err
	}

	// Make HTTP POST request to API
	url := fmt.Sprintf("%s/api/v1/factorial", apiEndpoint)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return Response{
			StatusCode: 500,
			Body:       fmt.Sprintf(`{"error": "Failed to create request: %v"}`, err),
		}, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Response{
			StatusCode: 500,
			Body:       fmt.Sprintf(`{"error": "Failed to call API: %v"}`, err),
		}, err
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{
			StatusCode: 500,
			Body:       fmt.Sprintf(`{"error": "Failed to read response: %v"}`, err),
		}, err
	}

	return Response{
		StatusCode: resp.StatusCode,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
