# System Architecture

## Overview

The Factorial Calculation Service is designed as a distributed, event-driven system for computing factorials asynchronously. It follows microservices principles with clear separation of concerns.

## Architecture Diagram

```
                                     ┌─────────────────────┐
                                     │   AWS Step         │
                                     │   Functions        │
                                     └──────────┬──────────┘
                                                │
                                                ▼
┌──────────────┐                      ┌─────────────────────┐
│              │   POST /factorial    │                     │
│   Client     │─────────────────────▶│    API Service      │
│              │                      │    (Port 8080)      │
│              │   GET /factorial/:n  │                     │
│              │◀─────────────────────│  - Gin Framework    │
└──────────────┘                      │  - REST Endpoints   │
                                      │  - Swagger Docs     │
                                      └──────────┬──────────┘
                                                 │
                        ┌────────────────────────┼────────────────────────┐
                        │                        │                        │
                        ▼                        ▼                        ▼
              ┌──────────────────┐    ┌──────────────────┐    ┌──────────────────┐
              │   RabbitMQ       │    │   Redis          │    │   MySQL          │
              │   Message Queue  │    │   Cache          │    │   Database       │
              │                  │    │                  │    │                  │
              │ - factorial      │    │ Key Format:      │    │ Tables:          │
              │   queue          │    │ factorial:{num}  │    │ - factorial_     │
              │                  │    │                  │    │   calculations   │
              └─────────┬────────┘    │ TTL: 24h         │    │                  │
                        │             │ Only < 10K       │    │ Indexes:         │
                        │             └──────────────────┘    │ - number         │
                        │                                     │ - status         │
                        ▼                                     │ - created_at     │
              ┌──────────────────┐                           └──────────────────┘
              │   Worker         │                                     ▲
              │   Service        │                                     │
              │                  │                                     │
              │ - Consume msgs   │────────────────────────────────────┘
              │ - Calculate      │
              │ - Upload S3      │
              │ - Cache Redis    │
              └─────────┬────────┘
                        │
                        ▼
              ┌──────────────────┐
              │   AWS S3         │
              │                  │
              │ Bucket:          │
              │ factorial-results│
              │                  │
              │ Key Format:      │
              │ factorials/{n}.txt│
              └──────────────────┘
```

## Components

### 1. API Service

**Responsibilities:**
- Accept HTTP requests for factorial calculations
- Validate input (0-20,000)
- Publish messages to RabbitMQ
- Query results from Redis/S3/MySQL
- Serve Swagger documentation

**Technology Stack:**
- Go 1.21
- Gin Web Framework
- GORM for database access
- AWS SDK for S3
- Redis client

**Scaling:**
- Horizontal pod autoscaling (HPA)
- Kubernetes deployment with 2+ replicas
- LoadBalancer service for traffic distribution

### 2. Worker Service

**Responsibilities:**
- Consume messages from RabbitMQ queue
- Calculate factorials using `big.Int`
- Upload all results to S3
- Cache small results (< 10,000) to Redis
- Update MySQL metadata

**Technology Stack:**
- Go 1.21
- RabbitMQ consumer
- AWS SDK for S3
- Redis client
- GORM for database access

**Scaling:**
- Multiple independent worker pods (3-10)
- No coordination required between workers
- HPA based on CPU utilization (70%)

### 3. RabbitMQ

**Configuration:**
- Queue: `factorial_queue`
- Durable: Yes
- Auto-delete: No
- Prefetch: 1 message per worker
- No DLQ (messages rejected without requeue)

**Message Format:**
```json
{
  "number": "10"
}
```

### 4. Redis

**Purpose:** Cache frequently accessed small factorial results

**Configuration:**
- Key format: `factorial:{number}`
- TTL: 24 hours
- Eviction policy: LRU
- Only caches numbers < 10,000

**Why Redis for small numbers only:**
- Small factorials are queried frequently
- Large factorial results can be multi-MB in size
- Cost-effective to store large results only in S3
- Redis memory is more expensive than S3 storage

### 5. MySQL

**Purpose:** Store metadata for all calculations

**Schema:**
```sql
CREATE TABLE factorial_calculations (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    number VARCHAR(100) NOT NULL UNIQUE,
    status VARCHAR(20) NOT NULL,
    s3_key VARCHAR(512) NOT NULL,
    checksum VARCHAR(64),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_number (number),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
);
```

**Status Values:**
- `calculating`: Worker is computing the result
- `uploading`: Uploading result to S3
- `done`: Calculation complete and available
- `failed`: Calculation failed

### 6. AWS S3

**Purpose:** Long-term storage for all factorial results

**Configuration:**
- Bucket: `factorial-results`
- Key format: `factorials/{number}.txt`
- Storage class: Standard (can be optimized to Intelligent-Tiering)
- Versioning: Disabled (results are immutable)
- Encryption: Server-side (SSE-S3)

**File Format:** Plain text containing the factorial result as a string

### 7. AWS Step Functions

**Purpose:** Orchestrate factorial calculations via event-driven architecture

**Workflow:**
1. Step Functions receives event with `{number: "10"}`
2. Invokes Lambda function
3. Lambda calls POST /factorial API
4. Returns execution result

## Data Flow

### Calculation Flow

1. **Client submits request:**
   ```
   POST /api/v1/factorial
   {"number": "10"}
   ```

2. **API validates and publishes:**
   - Validates: 0 ≤ number ≤ 20,000
   - Checks if already calculated (MySQL)
   - Publishes to RabbitMQ
   - Returns 202 Accepted

3. **Worker processes:**
   - Consumes message from queue
   - Checks if already done (MySQL)
   - Creates/updates DB record (status=calculating)
   - Calculates factorial using `big.Int`
   - Uploads result to S3
   - Updates DB (status=done, s3_key)
   - If < 10,000: caches to Redis
   - Acks message

4. **Client retrieves result:**
   ```
   GET /api/v1/factorial/10
   ```
   - If < 10,000: Check Redis → S3 (if miss)
   - If ≥ 10,000: Get from S3 directly
   - Return result

### Retrieval Flow

```
GET /factorial/:number
         │
         ▼
   Parse number
         │
         ▼
    < 10,000?
    /        \
  Yes         No
  │           │
  ▼           ▼
Check      Query DB
Redis      for s3_key
  │           │
Found?        ▼
/    \    Download
Yes   No   from S3
│     │       │
│     └───────┤
│             │
└─────────────┤
              ▼
         Return result
              │
              ▼
      Cache to Redis
      (if < 10,000)
```

## Scalability Considerations

### Horizontal Scaling

- **API**: Stateless, can scale to 10+ pods
- **Worker**: Independent consumers, can scale to 10+ pods
- **RabbitMQ**: Single instance (can use clustering for HA)
- **Redis**: Single instance (can use Redis Cluster for scale)
- **MySQL**: Single RDS instance (can use read replicas)

### Performance Optimizations

1. **Caching Strategy:**
   - Redis for hot data (< 10K)
   - S3 for cold data (all results)
   - Reduces database load

2. **Database Indexing:**
   - Index on `number` for fast lookups
   - Index on `status` for filtering
   - Index on `created_at` for time-based queries

3. **Connection Pooling:**
   - GORM connection pooling for MySQL
   - Redis connection pooling
   - RabbitMQ channel pooling

4. **Worker Design:**
   - Prefetch count = 1 (fair distribution)
   - Manual acks (reliability)
   - No message requeuing (avoid infinite loops)

## Security

### API Security

- No authentication (can add JWT/OAuth2)
- Input validation (prevent injection)
- Rate limiting (can add with middleware)

### AWS Security

- IAM roles for EKS (IRSA)
- S3 bucket policies (private access)
- Encryption at rest (SSE-S3)
- Encryption in transit (TLS)

### Database Security

- Connection over TLS
- Credentials in Kubernetes secrets
- No public access

### RabbitMQ Security

- TLS for connections
- Credentials in secrets
- Virtual host isolation

## Monitoring and Observability

### Metrics

- API request rate and latency
- Worker processing rate
- Queue depth
- Cache hit rate
- Error rates

### Logging

- Structured logging (JSON format)
- Log levels: INFO, WARN, ERROR
- Contextual information (request ID, number)

### Tracing

- Distributed tracing (can add OpenTelemetry)
- Request flow visualization

### Alerting

- High queue depth (> 1000 messages)
- Worker failures
- API errors (5xx)
- Database connection issues

## Disaster Recovery

### Backup Strategy

- MySQL: Daily automated backups (RDS)
- S3: Cross-region replication (optional)
- Configuration: GitOps (ArgoCD)

### Recovery Procedures

1. **Database failure:** Restore from RDS backup
2. **S3 failure:** AWS handles durability (99.999999999%)
3. **RabbitMQ failure:** Messages may be lost (can add persistence)
4. **Worker failure:** Kubernetes restarts pods automatically

## Future Enhancements

1. **Authentication & Authorization**
   - JWT tokens
   - API keys
   - Rate limiting per user

2. **Result Streaming**
   - WebSocket for real-time updates
   - Server-sent events (SSE)

3. **Batch Processing**
   - Bulk calculation requests
   - Parallel computation

4. **Advanced Caching**
   - CDN for API responses
   - Redis Cluster for scale

5. **Analytics**
   - Usage statistics
   - Popular numbers
   - Performance metrics dashboard

6. **Cost Optimization**
   - S3 Intelligent-Tiering
   - Spot instances for workers
   - Reserved capacity for RDS

## Conclusion

This architecture provides a robust, scalable solution for factorial calculations with clear separation of concerns, efficient resource utilization, and room for future growth.

