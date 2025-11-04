package aws

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
)

// StepFunctionsClient handles AWS Step Functions operations
type StepFunctionsClient interface {
	StartExecution(ctx context.Context, number string) (string, error)
}

type stepFunctionsClient struct {
	client          *sfn.Client
	stateMachineArn string
}

// NewStepFunctionsClient creates a new Step Functions client
func NewStepFunctionsClient(cfg aws.Config, stateMachineArn string) StepFunctionsClient {
	return &stepFunctionsClient{
		client:          sfn.NewFromConfig(cfg),
		stateMachineArn: stateMachineArn,
	}
}

// StartExecution starts a Step Functions execution for factorial calculation
func (c *stepFunctionsClient) StartExecution(ctx context.Context, number string) (string, error) {
	// Create input JSON
	input := map[string]string{
		"number": number,
	}
	inputJSON, err := json.Marshal(input)
	if err != nil {
		return "", fmt.Errorf("failed to marshal input: %w", err)
	}

	// Start execution
	result, err := c.client.StartExecution(ctx, &sfn.StartExecutionInput{
		StateMachineArn: aws.String(c.stateMachineArn),
		Input:           aws.String(string(inputJSON)),
	})
	if err != nil {
		return "", fmt.Errorf("failed to start execution: %w", err)
	}

	return *result.ExecutionArn, nil
}
