## History

### 2025-11-01: Complete Factorial Calculation Service Implementation
- ✅ Implemented full factorial calculation service with string-based I/O
- ✅ Created PostgreSQL migration with proper schema (number as VARCHAR(100))
- ✅ Built domain models and repository layer (GORM)
- ✅ Implemented core services:
  - Factorial calculation service with big.Int support (0-20,000)
  - Redis caching service (LRU for numbers < 10,000)
  - S3 storage service (all results stored in S3)
- ✅ Updated configuration to support Redis, S3, AWS Step Functions
- ✅ Created REST API with Gin framework:
  - POST /api/v1/factorial (submit calculation)
  - GET /api/v1/factorial/:number (get result)
  - GET /api/v1/factorial/metadata/:number (get status)
- ✅ Built worker service for async processing
- ✅ Updated RabbitMQ consumer handler with factorial message processing
- ✅ Implemented AWS Step Functions client and Lambda trigger
- ✅ Created Terraform configuration for Step Functions infrastructure
- ✅ Built complete Helm charts for Kubernetes deployment:
  - API deployment with LoadBalancer service
  - Worker deployment with HPA (3-10 replicas)
  - ConfigMaps and Secrets management
  - ServiceAccount for AWS IRSA
- ✅ Setup CI/CD pipelines:
  - GitHub Actions for CI (test, lint, build)
  - GitHub Actions for CD (ECR push, Helm update)
  - ArgoCD application manifest for GitOps
- ✅ Updated Dockerfile with multi-stage builds (API and Worker targets)
- ✅ Wrote comprehensive unit tests:
  - Factorial service tests (edge cases, validation)
  - Redis service tests (with miniredis mock)
  - Repository tests (with SQLite in-memory DB)
- ✅ Created extensive documentation:
  - Comprehensive README with setup instructions
  - ARCHITECTURE.md with detailed system design
  - API documentation in code (Swagger annotations)

**Architecture Summary:**
- All results stored in S3 (plain text)
- Small numbers (< 10,000) cached in Redis (24h TTL)
- Large numbers (>= 10,000) skip Redis cache
- PostgreSQL stores only metadata (status, s3_key, timestamps, checksum, size)
- Multiple independent workers for horizontal scaling
- Event-driven architecture with RabbitMQ
- AWS Step Functions for orchestration
- Full Kubernetes support with Helm and ArgoCD

### 2025-10-31: Simplified RabbitMQ Queue Architecture
- Removed dead letter queue (DLQ) functionality
- Removed retry exchange and retry queue
- Simplified to single main queue architecture
- Messages that fail processing are rejected without requeue (Nack with requeue=false)
- Deleted `pkg/consumer/rabbitmq_dlq.go` (no longer needed)
- Updated `pkg/consumer/rabbitmq_queue_setup.go` - now only declares main queue
- Updated `pkg/consumer/rabbitmq_message_handler.go` - simplified error handling without retry logic
- Updated `pkg/consumer/rabbitmq_order_consumer.go` - removed DLQ references from logs
- Producer remains unchanged (already simple)