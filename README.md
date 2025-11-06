# Factorial Calculation Service

A distributed factorial calculation service built with Go, using RabbitMQ for message queuing, Redis for caching, PostgreSQL for persistence, and AWS S3 for storing large factorial results.


# How to start LOCAL:
- rename and fix the .env or the env in the docker compose file
- Storage: Create S3 bucket or use local disk. Or you can implement memory storage layer. Add base case in the /infrastructure/base_case
- docker compose -f Docker-compose.yml up -d --build
- Refer the docs/swagger.json for API info. There are 3 APIs:
    - For request the calculate number only.
    POST /api/v1/factorial
    {
        "number": 4
    }

    - GET the result of the number
    GET /api/v1/factorial/{number}

    - GET metadata: Get the key, bucket, checksum, status of the request calculate. May use for client call S3 get the factorial of big numbers result
    GET /api/v1/factorial/metadata/{number}
- Call API request number. After that call GET the result of the number.


# How to deploy to AWS ECS:
- AWS configure ID, Key, region. 
- Refer context some info setup.
- Storage: Create S3 bucket or use local disk. Or you can implement memory storage layer. Add base case in the /infrastructure/base_case
- Use step 1.1 or 1.2
1.1 Refer to the /infrastructure/terraform if u you Terraform
1.2 refer to /infrastructure/ecs and add to chatbot for generate aws commands to create
- Base on the .env file. Set the env files so the service can get it.
- Call API request number. After that call GET the result of the number.


# For the deployed version, for simple no Auth require:
Swagger: https://express-nodejs-demo-alb-1541480660.us-east-1.elb.amazonaws.com/swagger/index.html

- POST https://express-nodejs-demo-alb-1541480660.us-east-1.elb.amazonaws.com/api/v1/factorial
{
    "number": 4
}

- GET https://express-nodejs-demo-alb-1541480660.us-east-1.elb.amazonaws.com/api/v1/factorial/4

- GET https://express-nodejs-demo-alb-1541480660.us-east-1.elb.amazonaws.com/api/v1/factorial/metadata/4


## Prerequisites

### For Docker Compose
- Docker
- RabbitMQ
- Redis
- AWS

## Project Structure

```
factorial-cal-services/
├── cmd/                          # Application entry points
│   ├── api/                      # API server
│   │   └── main.go
│   ├── calculator/               # Calculator service (continuous calculation)
│   │   └── main.go
│   ├── migrate/                  # Database migration tool
│   │   └── main.go
│   └── worker/                   # Worker service (message queue consumer)
│       └── main.go
│
├── pkg/                          # Application packages
│   ├── aws/                      # AWS clients (Step Functions, etc.)
│   │   └── stepfunctions_client.go
│   ├── config/                   # Configuration management
│   │   └── config.go
│   ├── consumer/                 # Message queue consumers
│   │   ├── interface.go
│   │   ├── rabbitmq_consumer.go
│   │   ├── rabbitmq_handler.go
│   │   └── rabbitmq_queue_setup.go
│   ├── db/                       # Database connection
│   │   └── gorm.go
│   ├── domain/                   # Domain models/entities
│   │   ├── factorial.go
│   │   ├── factorial_max.go
│   │   └── factorial_request.go
│   ├── dto/                      # Data transfer objects
│   │   ├── factorial.go
│   │   └── factorial_message.go
│   ├── handler/                  # HTTP handlers
│   │   ├── factorial_handler.go
│   │   └── helper_fn.go
│   ├── producer/                 # Message queue producers
│   │   ├── interface.go
│   │   └── rabbitmq_producer.go
│   ├── repository/               # Data access layer
│   │   ├── current_calculated_repository.go
│   │   ├── factorial_repository.go
│   │   └── max_request_repository.go
│   ├── service/                  # Business logic layer
│   │   ├── factorial_service.go
│   │   ├── redis_service.go
│   │   ├── storage_interface.go
│   │   ├── storage_local_service.go
│   │   └── storage_s3_service.go
│   └── utils/                    # Utility functions
│       └── patterns/
│           ├── semaphore.go
│           └── worker_pool.go
│
├── migrations/                   # Database migrations
│   ├── 000001_init.up.sql
│   ├── 000001_init.down.sql
│   ├── 000002_additional_tables.up.sql
│   ├── 000002_additional_tables.down.sql
│   └── migration.go
│
├── docs/                         # API documentation (Swagger)
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
│
├── infrastructure/               # Infrastructure as Code
│   ├── helm/                     # Helm charts
│   │   ├── Chart.yaml
│   │   ├── values.yaml
│   │   ├── README.md
│   │   └── templates/
│   │       ├── api-deployment.yaml
│   │       ├── calculator-deployment.yaml
│   │       ├── worker-deployment.yaml
│   │       ├── service.yaml
│   │       ├── secret.yaml
│   │       ├── serviceaccount.yaml
│   │       ├── hpa.yaml
│   │       └── _helpers.tpl
│   ├── terraform/                # Terraform configurations
│   │   ├── vpc.tf
│   │   ├── subnets.tf
│   │   ├── security_groups.tf
│   │   └── ...
│   └── argocd/                   # ArgoCD configurations
│       └── application.yaml
│
├── context/                      # Architecture and design docs
│   ├── arch.md
│   ├── api_design.md
│   ├── database.md
│   └── flow.md
│
├── Dockerfile                    # Multi-stage Dockerfile
├── Docker-compose.yml            # Docker Compose configuration
├── Makefile                      # Build and deployment commands
├── go.mod                        # Go module dependencies
└── README.md                     # This file
```

## Architecture Overview

### Services

1. **API Service** (`cmd/api`)
   - REST API for factorial calculation requests
   - Publishes messages to RabbitMQ queue
   - Retrieves results from Redis cache or S3
   - Runs database migrations on startup

2. **Worker Service** (`cmd/worker`)
   - Consumes messages from RabbitMQ queue
   - Updates max request number in database
   - Processes messages in batches

3. **Calculator Service** (`cmd/calculator`)
   - Continuously calculates factorials sequentially
   - Reads current and max numbers from database
   - Stores results in S3
   - Updates database with calculation results

### Data Flow

```
User Request → API → RabbitMQ Queue → Worker → Database
                                          ↓
                                    Calculator → S3 → Database
                                          ↓
                                    Redis Cache ← API Response
```