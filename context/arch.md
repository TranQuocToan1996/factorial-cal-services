## Architecture Overview

Factorial calculation service that handles small and large factorials using caching, async processing, and AWS services.

### Core Components

1. API layer
   - Handles HTTP requests (POST/GET)
   - Validates input and routes requests

2. Message queue (RabbitMQ)
   - Receives calculation requests asynchronously
   - Decouples API from computation

3. Worker service
   - Consumes messages from RabbitMQ
   - Checks `max_factorial` and `cur_factorial` from DB
   - Triggers AWS Step Functions for large calculations

4. AWS Step Functions
   - Orchestrates calculation and storage for large factorials
   - Handles multipart uploads to S3
   - Updates metadata in PostgreSQL

### Data Stores

1. Redis cache
   - Stores small factorial results (< 1 billion)
   - LRU eviction
   - Fast retrieval for common queries

2. PostgreSQL (factorial DB)
   - Stores metadata:
     - `n` (bigint, primary key)
     - `s3_key` (not null)
     - `size` (bigint)
     - `checksum` (varchar(64))
     - `status` (enum: calculating, uploading, done)
     - `created_at` (timestamp)
   - Tracks `max_factorial` and `cur_factorial` for worker coordination

3. S3
   - Stores large factorial results as objects
   - Multipart uploads for very large files
   - Accessed via `s3_key` from metadata

### Request Flows

#### 1. POST /factorial (submit calculation)
```
Client → API → RabbitMQ Topic → Worker (checks max_factorial) → AWS Step Function
```

#### 2. GET /factorial?number=X (small factorial, < 1 billion)
```
Client → API → Redis (direct lookup) → Return cached result
```

#### 3. GET /factorial/metadata?number=X (large factorial)
```
Client → API → PostgreSQL (metadata lookup) → S3 (retrieve via s3_key) → Return result
```

### Design Decisions

- Small vs large: Small results cached in Redis; large results stored in S3
- Async processing: RabbitMQ decouples request submission from computation
- Step Functions: Orchestrates multi-step workflows for large calculations
- Lazy calculation: Large factorials computed on-demand
- Data integrity: Checksums stored in metadata
- Scalability: Redis for speed, S3 for large objects, PostgreSQL for metadata

### Infrastructure Assumptions

- AWS hosted
- RAM constraints limit in-memory calculation
- Redis uses LRU eviction
- Clients frequently access small factorials
- Large factorials are queried less often

This architecture separates concerns, scales caching and storage, and uses Step Functions for orchestration of complex workflows.


List of arch now:
- API Golang (EKS)
- RabbitMQ (Helm -> EKS)
- Worker Golang (EKS)
- Redis (Helm -> EKS)
- Storage Engine (Now S3-AWS)
- MySQL (Now RDS-AWS)