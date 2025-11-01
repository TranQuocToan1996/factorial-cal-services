# Factorial Calculation Service

A highly scalable, asynchronous factorial calculation service built with Go, designed to handle large factorial computations efficiently using a distributed architecture.

## Architecture

### System Overview

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│   Client    │────▶│   API (Go)   │────▶│  RabbitMQ   │
└─────────────┘     └──────────────┘     └─────────────┘
                            │                     │
                            ▼                     ▼
                    ┌──────────────┐     ┌─────────────┐
                    │    Redis     │     │   Worker    │
                    │   (Cache)    │     │   (Go)      │
                    └──────────────┘     └─────────────┘
                            │                     │
                            │                     ▼
                    ┌──────────────┐     ┌─────────────┐
                    │    MySQL     │◀────│   S3        │
                    │  (Metadata)  │     │  (Results)  │
                    └──────────────┘     └─────────────┘
```

### Components

- **API Service**: REST API for submitting calculations and retrieving results
- **Worker Service**: Processes factorial calculations asynchronously
- **RabbitMQ**: Message queue for async job processing
- **Redis**: LRU cache for small factorial results (< 10,000)
- **MySQL**: Stores metadata (status, S3 keys, timestamps)
- **S3**: Stores all factorial results as plain text
- **AWS Step Functions**: Orchestration layer for triggering calculations

### Storage Strategy

- **Small numbers (< 10,000)**: Results stored in both S3 and Redis
- **Large numbers (≥ 10,000)**: Results stored only in S3 (skip Redis)
- **Database**: Only metadata (number, status, s3_key, created_at)

## Features

- ✅ Async factorial calculation for numbers 0-20,000
- ✅ String-based input/output for handling large numbers
- ✅ Redis caching for frequently accessed results
- ✅ S3 storage for all results
- ✅ Multiple independent workers for horizontal scaling
- ✅ Kubernetes-ready with Helm charts
- ✅ CI/CD with GitHub Actions and ArgoCD
- ✅ AWS Step Functions integration
- ✅ Comprehensive unit tests
- ✅ API documentation with Swagger

## API Endpoints

### POST /api/v1/factorial

Submit a factorial calculation request.

**Request:**
```json
{
  "number": "10"
}
```

**Response (202 Accepted):**
```json
{
  "number": "10",
  "status": "accepted"
}
```

### GET /api/v1/factorial/:number

Get the factorial result for a number.

**Response (200 OK):**
```json
{
  "number": "10",
  "result": "3628800",
  "status": "done"
}
```

**Response (202 Accepted - Still Processing):**
```json
{
  "error": "processing",
  "message": "Calculation is still in progress (status: calculating)"
}
```

### GET /api/v1/factorial/metadata/:number

Get metadata for a calculation.

**Response (200 OK):**
```json
{
  "number": "10",
  "status": "done",
  "s3_key": "factorials/10.txt",
  "created_at": "2025-11-01T12:00:00Z"
}
```

### GET /health

Health check endpoint.

**Response (200 OK):**
```json
{
  "status": "ok"
}
```

## Quick Start

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- MySQL
- RabbitMQ
- Redis
- AWS Account (for S3 and Step Functions)

### Local Development

1. **Clone the repository:**
```bash
git clone https://github.com/your-org/factorial-cal-services.git
cd factorial-cal-services
```

2. **Set up environment variables:**
```bash
cp .env.example .env
# Edit .env with your configuration
```

3. **Run with Docker Compose:**
```bash
docker-compose up -d
```

4. **Run migrations:**
```bash
make migrate-up
```

5. **Start the API:**
```bash
make run-api
```

6. **Start the Worker:**
```bash
make run-worker
```

### Using Docker

**Build images:**
```bash
docker build -f Dockerfile --target api -t factorial-api:latest .
docker build -f Dockerfile --target worker -t factorial-worker:latest .
```

**Run API:**
```bash
docker run -p 8080:8080 --env-file .env factorial-api:latest
```

**Run Worker:**
```bash
docker run --env-file .env factorial-worker:latest
```

## Kubernetes Deployment

### Using Helm

1. **Install the Helm chart:**
```bash
helm install factorial-service ./helm/factorial-service \
  --namespace factorial-service \
  --create-namespace \
  --values ./helm/factorial-service/values.yaml
```

2. **Update values:**
```bash
helm upgrade factorial-service ./helm/factorial-service \
  --namespace factorial-service \
  --values ./helm/factorial-service/values.yaml
```

3. **Uninstall:**
```bash
helm uninstall factorial-service --namespace factorial-service
```

### Using ArgoCD

1. **Apply ArgoCD application:**
```bash
kubectl apply -f argocd/application.yaml
```

2. **Sync the application:**
```bash
argocd app sync factorial-service
```

## AWS Infrastructure

### Terraform Setup

1. **Initialize Terraform:**
```bash
cd terraform
terraform init
```

2. **Plan deployment:**
```bash
terraform plan -var="api_endpoint=https://your-api-endpoint.com"
```

3. **Apply configuration:**
```bash
terraform apply -var="api_endpoint=https://your-api-endpoint.com"
```

### Step Functions Integration

The service includes AWS Step Functions integration for orchestrating factorial calculations:

- State machine triggers Lambda function
- Lambda calls the POST /factorial API endpoint
- Supports event-driven architecture

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_PORT` | API server port | `:8080` |
| `DB_HOST` | MySQL host | `localhost` |
| `DB_PORT` | MySQL port | `3306` |
| `DB_USER` | MySQL username | `root` |
| `DB_PASSWORD` | MySQL password | - |
| `DB_NAME` | Database name | `factorial_db` |
| `RABBITMQ_HOST` | RabbitMQ host | `localhost` |
| `RABBITMQ_PORT` | RabbitMQ port | `5672` |
| `RABBITMQ_USER` | RabbitMQ username | `guest` |
| `RABBITMQ_PASSWORD` | RabbitMQ password | `guest` |
| `REDIS_HOST` | Redis host | `localhost` |
| `REDIS_PORT` | Redis port | `6379` |
| `REDIS_PASSWORD` | Redis password | - |
| `REDIS_DB` | Redis database | `0` |
| `AWS_REGION` | AWS region | `us-east-1` |
| `S3_BUCKET_NAME` | S3 bucket for results | - |
| `MAX_FACTORIAL` | Maximum allowed number | `20000` |
| `REDIS_THRESHOLD` | Cache threshold | `10000` |

## Testing

### Run Unit Tests

```bash
go test -v ./...
```

### Run Tests with Coverage

```bash
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Run Integration Tests (Optional)

```bash
go test -v -tags=integration ./...
```

## CI/CD Pipeline

### GitHub Actions

- **CI Pipeline** (`.github/workflows/ci.yaml`):
  - Runs on PRs and pushes to develop/main
  - Executes unit tests
  - Runs linting with golangci-lint
  - Builds binaries

- **CD Pipeline** (`.github/workflows/cd.yaml`):
  - Triggers on push to main
  - Builds and pushes Docker images to ECR
  - Updates Helm chart with new image tags
  - ArgoCD auto-syncs the changes

## Performance Considerations

- **Redis Caching**: Small numbers (< 10,000) are cached for 24 hours
- **Worker Scaling**: Multiple independent workers for parallel processing
- **S3 Storage**: All results stored for durability and cost-effectiveness
- **Database Indexing**: Optimized queries with proper indexes
- **Connection Pooling**: Efficient resource utilization

## Limitations

- Maximum factorial calculation: 20,000
- Input/output as strings to handle large numbers
- Redis cache limited to numbers < 10,000
- No support for negative numbers

## Troubleshooting

### API is not responding

Check if all services are running:
```bash
kubectl get pods -n factorial-service
```

### Worker not processing messages

Check RabbitMQ connection:
```bash
kubectl logs -n factorial-service deployment/factorial-service-worker
```

### Redis connection errors

Verify Redis is accessible:
```bash
redis-cli -h <redis-host> ping
```

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Commit changes: `git commit -am 'Add feature'`
4. Push to branch: `git push origin feature/my-feature`
5. Submit a pull request

## License

Apache 2.0 License

## Contact

For support or questions, please contact: support@example.com
