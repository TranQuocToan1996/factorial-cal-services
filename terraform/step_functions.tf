terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

variable "api_endpoint" {
  description = "The API endpoint URL for factorial service"
  type        = string
}

variable "lambda_function_name" {
  description = "Name of the Lambda function"
  type        = string
  default     = "factorial-trigger"
}

variable "state_machine_name" {
  description = "Name of the Step Functions state machine"
  type        = string
  default     = "factorial-calculator"
}

# IAM Role for Lambda
resource "aws_iam_role" "lambda_role" {
  name = "${var.lambda_function_name}-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}

# Lambda basic execution policy
resource "aws_iam_role_policy_attachment" "lambda_basic" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
  role       = aws_iam_role.lambda_role.name
}

# Lambda function
resource "aws_lambda_function" "factorial_trigger" {
  filename      = "lambda-deployment.zip"
  function_name = var.lambda_function_name
  role          = aws_iam_role.lambda_role.arn
  handler       = "main"
  runtime       = "go1.x"
  timeout       = 30

  environment {
    variables = {
      API_ENDPOINT = var.api_endpoint
    }
  }
}

# IAM Role for Step Functions
resource "aws_iam_role" "sfn_role" {
  name = "${var.state_machine_name}-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "states.amazonaws.com"
        }
      }
    ]
  })
}

# Policy for Step Functions to invoke Lambda
resource "aws_iam_role_policy" "sfn_lambda_policy" {
  name = "${var.state_machine_name}-lambda-policy"
  role = aws_iam_role.sfn_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "lambda:InvokeFunction"
        ]
        Resource = aws_lambda_function.factorial_trigger.arn
      }
    ]
  })
}

# Step Functions State Machine
resource "aws_sfn_state_machine" "factorial_calculator" {
  name     = var.state_machine_name
  role_arn = aws_iam_role.sfn_role.arn

  definition = jsonencode({
    Comment = "Factorial calculation workflow"
    StartAt = "TriggerFactorialCalculation"
    States = {
      TriggerFactorialCalculation = {
        Type     = "Task"
        Resource = aws_lambda_function.factorial_trigger.arn
        End      = true
      }
    }
  })
}

# Outputs
output "state_machine_arn" {
  description = "ARN of the Step Functions state machine"
  value       = aws_sfn_state_machine.factorial_calculator.arn
}

output "lambda_function_arn" {
  description = "ARN of the Lambda function"
  value       = aws_lambda_function.factorial_trigger.arn
}

