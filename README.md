# Factorial Calculation Service

A distributed factorial calculation service built with Go, using RabbitMQ for message queuing, Redis for caching, PostgreSQL for persistence, and AWS S3 for storing large factorial results.

## Features

- **REST API**: Submit factorial calculation requests
- **Async Processing**: Message queue-based processing with RabbitMQ
- **Caching**: Redis caching for frequently accessed results
- **Sequential Calculator**: Background service for continuous factorial calculations
- **S3 Storage**: Large factorial results stored in AWS S3
- **Health Checks**: Built-in health endpoints for monitoring

## Prerequisites

### For Docker Compose
- Docker 20.10+
- Docker Compose 2.0+

### For Kubernetes
- Kubernetes cluster (minikube, kind, or EKS)
- kubectl configured
- Helm 3.0+

## Quick Start - Docker Compose

### 1. Clone the Repository

```bash
git clone <repository-url>
cd factorial-cal-services
```

### 2. Start All Services

```bash
# Build and start all services
docker compose -f Docker-compose.yml up -d --build

# Or using Makefile
make up
```

This will start:
- PostgreSQL (port 5432)
- Redis (port 6379)
- RabbitMQ (ports 5672, 15672)
- API service (port 8080)
- Worker service
- Calculator service

### 3. Verify Services

```bash
# Check all containers are running
docker compose ps

# Check API health
curl http://localhost:8080/health

# View logs
docker compose logs -f api
docker compose logs -f worker
docker compose logs -f calculator
```

### 4. Access Services

- **API**: http://localhost:8080
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **RabbitMQ Management**: http://localhost:15672 (guest/guest)

### 5. Stop Services

```bash
# Stop all services
docker compose down

# Stop and remove volumes
docker compose down -v

# Or using Makefile
make down
```

## Quick Start - Kubernetes (Local)

### Option 1: Using Minikube

#### 1. Start Minikube

```bash
minikube start
kubectl create namespace factorial
```

#### 2. Deploy Infrastructure Services

First, deploy PostgreSQL, Redis, and RabbitMQ:

```bash
# Install PostgreSQL (using Bitnami chart)
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install postgres bitnami/postgresql -n factorial \
  --set auth.postgresPassword=password \
  --set auth.database=factorial-cal-services

# Install Redis
helm install redis bitnami/redis -n factorial \
  --set auth.enabled=false

# Install RabbitMQ
helm install rabbitmq bitnami/rabbitmq -n factorial \
  --set auth.username=guest \
  --set auth.password=guest
```

#### 3. Build and Load Images

```bash
# Build images
docker build --target api -t factorial-api:latest .
docker build --target worker -t factorial-worker:latest .
docker build --target calculator -t factorial-calculator:latest .

# Load into minikube
minikube image load factorial-api:latest
minikube image load factorial-worker:latest
minikube image load factorial-calculator:latest
```

#### 4. Create Secrets

```bash
# Get PostgreSQL password
POSTGRES_PASSWORD=$(kubectl get secret postgres-postgresql -n factorial -o jsonpath="{.data.postgres-password}" | base64 -d)

# Get Redis password (if enabled)
REDIS_PASSWORD=""

# Get RabbitMQ password
RABBITMQ_PASSWORD=$(kubectl get secret rabbitmq -n factorial -o jsonpath="{.data.rabbitmq-password}" | base64 -d)

# Create application secrets
kubectl create secret generic factorial-service-secret -n factorial \
  --from-literal=db-password="$POSTGRES_PASSWORD" \
  --from-literal=rabbitmq-password="$RABBITMQ_PASSWORD" \
  --from-literal=redis-password="$REDIS_PASSWORD"
```

#### 5. Update Helm Values

Update `infrastructure/helm/values.yaml`:

```yaml
database:
  host: postgres-postgresql.factorial.svc.cluster.local
  port: "5432"
  name: factorial-cal-services
  user: postgres

redis:
  host: redis-master.factorial.svc.cluster.local
  port: "6379"

rabbitmq:
  host: rabbitmq.factorial.svc.cluster.local
  port: "5672"
  user: guest

api:
  image:
    repository: factorial-api
    tag: latest
    pullPolicy: Never  # Use local images in minikube

worker:
  image:
    repository: factorial-worker
    tag: latest
    pullPolicy: Never

calculator:
  image:
    repository: factorial-calculator
    tag: latest
    pullPolicy: Never
```

#### 6. Deploy Application

```bash
helm install factorial-service ./infrastructure/helm \
  --namespace factorial \
  --set database.host=postgres-postgresql.factorial.svc.cluster.local \
  --set redis.host=redis-master.factorial.svc.cluster.local \
  --set rabbitmq.host=rabbitmq.factorial.svc.cluster.local
```

#### 7. Access Services

```bash
# Port forward API service
kubectl port-forward -n factorial svc/factorial-service-api 8080:8080

# Access API
curl http://localhost:8080/health
```

### Option 2: Using Kind

#### 1. Create Kind Cluster

```bash
kind create cluster --name factorial
kubectl create namespace factorial
```

#### 2. Build and Load Images

```bash
# Build images
docker build --target api -t factorial-api:latest .
docker build --target worker -t factorial-worker:latest .
docker build --target calculator -t factorial-calculator:latest .

# Load into kind
kind load docker-image factorial-api:latest --name factorial
kind load docker-image factorial-worker:latest --name factorial
kind load docker-image factorial-calculator:latest --name factorial
```

#### 3. Follow steps 2-7 from Minikube section above

## Configuration

### Environment Variables

All services support the following environment variables:

#### Database
- `DB_HOST`: PostgreSQL host (default: localhost)
- `DB_PORT`: PostgreSQL port (default: 5432)
- `DB_NAME`: Database name (default: factorial-cal-services)
- `DB_USER`: Database user (default: postgres)
- `DB_PASSWORD`: Database password
- `DB_SSLMODE`: SSL mode (default: disable)

#### Redis
- `REDIS_HOST`: Redis host (default: localhost)
- `REDIS_PORT`: Redis port (default: 6379)
- `REDIS_PASSWORD`: Redis password (optional)
- `REDIS_THRESHOLD`: Cache threshold for numbers (default: 1000)

#### RabbitMQ
- `RABBITMQ_HOST`: RabbitMQ host (default: localhost)
- `RABBITMQ_PORT`: RabbitMQ port (default: 5672)
- `RABBITMQ_USER`: RabbitMQ user (default: guest)
- `RABBITMQ_PASSWORD`: RabbitMQ password
- `FACTORIAL_CAL_SERVICES_QUEUE_NAME`: Queue name (default: factorial-cal-queue)
- `RABBITMQ_CA`: TLS CA certificate (optional)

#### AWS
- `AWS_REGION`: AWS region (default: us-east-1)
- `S3_BUCKET_NAME`: S3 bucket name
- `STORAGE_TYPE`: Storage type - s3 or local (default: s3)

#### Application
- `SERVER_PORT`: API server port (default: :8080)
- `MAX_FACTORIAL`: Maximum factorial number (default: 10000)
- `WORKER_BATCH_SIZE`: Worker batch size (default: 100)
- `WORKER_MAX_BATCHES`: Maximum worker batches (default: 16)

## API Endpoints

### Health Check
```bash
GET /health
```

### Submit Calculation
```bash
POST /api/v1/factorial
Content-Type: application/json

{
  "number": 10
}
```

### Get Result
```bash
GET /api/v1/factorial/{number}
```

### Get Metadata
```bash
GET /api/v1/factorial/metadata/{number}
```

See Swagger documentation at `/swagger/index.html` for detailed API documentation.

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

## Development

### Run Locally (without Docker)

```bash
# Install dependencies
make deps

# Run migrations
export DATABASE_URL="postgres://postgres:password@localhost:5432/factorial-cal-services?sslmode=disable"
make migrate-up

# Run API
go run cmd/api/main.go

# Run Worker (in another terminal)
go run cmd/worker/main.go

# Run Calculator (in another terminal)
go run cmd/calculator/main.go
```

### Run Tests

```bash
# Run all tests
make test

# Run tests with coverage
go test -v ./... -cover
```

### Generate Swagger Documentation

```bash
make swagger
```

## Troubleshooting

### Docker Compose Issues

```bash
# View logs
docker compose logs -f

# Restart a specific service
docker compose restart api

# Rebuild and restart
docker compose up -d --build api
```

### Kubernetes Issues

```bash
# Check pod status
kubectl get pods -n factorial

# View logs
kubectl logs -n factorial -l component=api

# Describe pod for events
kubectl describe pod <pod-name> -n factorial

# Check service endpoints
kubectl get endpoints -n factorial
```

### Database Connection Issues

- Verify database is running: `docker compose ps postgres`
- Check connection string format
- Ensure migrations have run: Check logs for migration errors

### Redis Connection Issues

- Verify Redis is running: `docker compose ps redis`
- Test connection: `docker compose exec redis redis-cli ping`

### RabbitMQ Connection Issues

- Verify RabbitMQ is running: `docker compose ps rabbitmq`
- Check management UI: http://localhost:15672
- Verify queue exists in management UI

## How to Start in AWS (EKS)

This guide covers deploying the Factorial Calculation Service to Amazon EKS (Elastic Kubernetes Service).

### Prerequisites

- AWS CLI configured with appropriate permissions
- kubectl installed and configured
- Helm 3.0+ installed
- eksctl installed (for cluster creation) or existing EKS cluster
- Docker installed (for building images)
- AWS IAM permissions for:
  - EKS cluster management
  - ECR (Elastic Container Registry)
  - S3 bucket creation/access
  - IAM role creation for IRSA

### Step 1: Create EKS Cluster

#### Option A: Using eksctl (Recommended)

```bash
# Create EKS cluster with managed node group
eksctl create cluster \
  --name factorial-cluster \
  --region us-east-1 \
  --node-type t3.medium \
  --nodes 3 \
  --nodes-min 2 \
  --nodes-max 5 \
  --managed

# Configure kubectl
aws eks update-kubeconfig --name factorial-cluster --region us-east-1
```

#### Option B: Using AWS Console or Terraform

If using Terraform, see `infrastructure/terraform/` for VPC and networking setup.

```bash
# After cluster creation, configure kubectl
aws eks update-kubeconfig --name your-cluster-name --region us-east-1
```

### Step 2: Create ECR Repository

```bash
# Set variables
AWS_REGION=us-east-1
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
ECR_REPO_NAME=factorial-cal-services

# Create ECR repositories
aws ecr create-repository \
  --repository-name factorial-api \
  --region $AWS_REGION

aws ecr create-repository \
  --repository-name factorial-worker \
  --region $AWS_REGION

aws ecr create-repository \
  --repository-name factorial-calculator \
  --region $AWS_REGION

# Get ECR login token
aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com
```

### Step 3: Build and Push Docker Images

```bash
# Set ECR base URL
ECR_BASE=$AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com

# Build and push API image
docker build --target api -t factorial-api:latest .
docker tag factorial-api:latest $ECR_BASE/factorial-api:latest
docker push $ECR_BASE/factorial-api:latest

# Build and push Worker image
docker build --target worker -t factorial-worker:latest .
docker tag factorial-worker:latest $ECR_BASE/factorial-worker:latest
docker push $ECR_BASE/factorial-worker:latest

# Build and push Calculator image
docker build --target calculator -t factorial-calculator:latest .
docker tag factorial-calculator:latest $ECR_BASE/factorial-calculator:latest
docker push $ECR_BASE/factorial-calculator:latest
```

### Step 4: Set Up AWS Infrastructure Services

#### Option A: Managed Services (Recommended for Production)

**RDS PostgreSQL:**
```bash
# Create RDS PostgreSQL instance
aws rds create-db-instance \
  --db-instance-identifier factorial-postgres \
  --db-instance-class db.t3.micro \
  --engine postgres \
  --engine-version 15.4 \
  --master-username postgres \
  --master-user-password YourSecurePassword123! \
  --allocated-storage 20 \
  --vpc-security-group-ids sg-xxxxx \
  --db-name factorial-cal-services \
  --backup-retention-period 7
```

**ElastiCache Redis:**
```bash
# Create ElastiCache Redis cluster
aws elasticache create-cache-cluster \
  --cache-cluster-id factorial-redis \
  --cache-node-type cache.t3.micro \
  --engine redis \
  --num-cache-nodes 1 \
  --engine-version 7.0
```

**Amazon MQ (RabbitMQ):**
```bash
# Create Amazon MQ broker
aws mq create-broker \
  --broker-name factorial-rabbitmq \
  --broker-deployment-mode SINGLE_INSTANCE \
  --engine-type rabbitmq \
  --engine-version 3.11.20 \
  --host-instance-type mq.t3.micro \
  --publicly-accessible false \
  --users Username=guest,Password=YourSecurePassword123!
```

#### Option B: Deploy in EKS (For Development/Testing)

```bash
# Create namespace
kubectl create namespace factorial

# Install PostgreSQL using Bitnami chart
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install postgres bitnami/postgresql -n factorial \
  --set auth.postgresPassword=YourSecurePassword123! \
  --set auth.database=factorial-cal-services

# Install Redis
helm install redis bitnami/redis -n factorial \
  --set auth.enabled=false

# Install RabbitMQ
helm install rabbitmq bitnami/rabbitmq -n factorial \
  --set auth.username=guest \
  --set auth.password=YourSecurePassword123!
```

### Step 5: Create S3 Bucket

```bash
# Create S3 bucket
aws s3 mb s3://factorial-cal-services-prod --region $AWS_REGION

# Enable versioning (optional but recommended)
aws s3api put-bucket-versioning \
  --bucket factorial-cal-services-prod \
  --versioning-configuration Status=Enabled

# Create bucket policy for application access (if using IRSA)
cat > s3-bucket-policy.json <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::$AWS_ACCOUNT_ID:role/factorial-service-role"
      },
      "Action": [
        "s3:GetObject",
        "s3:PutObject",
        "s3:DeleteObject"
      ],
      "Resource": "arn:aws:s3:::factorial-cal-services-prod/*"
    },
    {
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::$AWS_ACCOUNT_ID:role/factorial-service-role"
      },
      "Action": "s3:ListBucket",
      "Resource": "arn:aws:s3:::factorial-cal-services-prod"
    }
  ]
}
EOF

aws s3api put-bucket-policy \
  --bucket factorial-cal-services-prod \
  --policy file://s3-bucket-policy.json
```

### Step 6: Set Up IAM Role for Service Account (IRSA)

```bash
# Create IAM policy for S3 access
cat > s3-policy.json <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:GetObject",
        "s3:PutObject",
        "s3:DeleteObject"
      ],
      "Resource": "arn:aws:s3:::factorial-cal-services-prod/*"
    },
    {
      "Effect": "Allow",
      "Action": "s3:ListBucket",
      "Resource": "arn:aws:s3:::factorial-cal-services-prod"
    }
  ]
}
EOF

aws iam create-policy \
  --policy-name factorial-s3-policy \
  --policy-document file://s3-policy.json

# Get policy ARN
POLICY_ARN=$(aws iam list-policies --query 'Policies[?PolicyName==`factorial-s3-policy`].Arn' --output text)

# Create IAM role with trust policy for EKS service account
OIDC_PROVIDER=$(aws eks describe-cluster --name factorial-cluster --query "cluster.identity.oidc.issuer" --output text | sed -e "s/^https:\/\///")
NAMESPACE=factorial
SERVICE_ACCOUNT_NAME=factorial-service-sa

cat > trust-policy.json <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Federated": "arn:aws:iam::$AWS_ACCOUNT_ID:oidc-provider/$OIDC_PROVIDER"
      },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "StringEquals": {
          "${OIDC_PROVIDER}:sub": "system:serviceaccount:$NAMESPACE:$SERVICE_ACCOUNT_NAME",
          "${OIDC_PROVIDER}:aud": "sts.amazonaws.com"
        }
      }
    }
  ]
}
EOF

# Create IAM role
aws iam create-role \
  --role-name factorial-service-role \
  --assume-role-policy-document file://trust-policy.json

# Attach S3 policy to role
aws iam attach-role-policy \
  --role-name factorial-service-role \
  --policy-arn $POLICY_ARN

# Get role ARN
ROLE_ARN=$(aws iam get-role --role-name factorial-service-role --query 'Role.Arn' --output text)
```

### Step 7: Configure Helm Values for EKS

Create `values-eks.yaml`:

```yaml
api:
  replicaCount: 3
  image:
    repository: $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/factorial-api
    tag: latest
    pullPolicy: Always
  resources:
    requests:
      memory: "512Mi"
      cpu: "500m"
    limits:
      memory: "1Gi"
      cpu: "1000m"
  service:
    type: ClusterIP  # Use ALB ingress instead
    port: 8080
    targetPort: 8080

worker:
  replicaCount: 5
  image:
    repository: $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/factorial-worker
    tag: latest
    pullPolicy: Always
  resources:
    requests:
      memory: "512Mi"
      cpu: "500m"
    limits:
      memory: "1Gi"
      cpu: "1000m"
  autoscaling:
    enabled: true
    minReplicas: 5
    maxReplicas: 20
    targetCPUUtilizationPercentage: 70

calculator:
  replicaCount: 1
  image:
    repository: $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/factorial-calculator
    tag: latest
    pullPolicy: Always
  resources:
    requests:
      memory: "512Mi"
      cpu: "500m"
    limits:
      memory: "1Gi"
      cpu: "1000m"

config:
  serverPort: ":8080"
  maxFactorial: 20000
  redisThreshold: 10000
  queueName: "factorial-cal-queue"
  workerBatchSize: 100
  workerMaxBatches: 16
  storageType: "s3"
  queueType: "rabbitmq"

database:
  host: "factorial-postgres.xxxxx.us-east-1.rds.amazonaws.com"  # RDS endpoint
  port: "5432"
  name: "factorial-cal-services"
  user: "postgres"
  sslmode: "require"  # Use SSL for RDS

rabbitmq:
  host: "amqps://factorial-rabbitmq-xxxxx.mq.us-east-1.amazonaws.com"  # Amazon MQ endpoint
  port: "5671"  # SSL port
  user: "guest"

redis:
  host: "factorial-redis.xxxxx.0001.use1.cache.amazonaws.com"  # ElastiCache endpoint
  port: "6379"
  db: 0

aws:
  region: "us-east-1"
  s3Bucket: "factorial-cal-services-prod"
  stepFunctionsArn: ""

serviceAccount:
  create: true
  annotations:
    eks.amazonaws.com/role-arn: "arn:aws:iam::$AWS_ACCOUNT_ID:role/factorial-service-role"
  name: "factorial-service-sa"
```

**Important:** Replace placeholders with actual values:
- `$AWS_ACCOUNT_ID` with your AWS account ID
- `$AWS_REGION` with your region
- Database/RabbitMQ/Redis hosts with actual endpoints from AWS console

### Step 8: Create Kubernetes Secrets

```bash
# Create namespace
kubectl create namespace factorial

# Create secrets
kubectl create secret generic factorial-service-secret -n factorial \
  --from-literal=db-password='YourSecurePassword123!' \
  --from-literal=rabbitmq-password='YourSecurePassword123!' \
  --from-literal=redis-password=''

# Or use AWS Secrets Manager (recommended for production)
# First, create secret in AWS Secrets Manager
aws secretsmanager create-secret \
  --name factorial-service-secrets \
  --secret-string '{"db-password":"YourSecurePassword123!","rabbitmq-password":"YourSecurePassword123!","redis-password":""}'

# Then use External Secrets Operator or similar to sync
```

### Step 9: Deploy with Helm

```bash
# Add any required Helm repositories
helm repo add aws-efs-csi-driver https://kubernetes-sigs.github.io/aws-efs-csi-driver/
helm repo update

# Install the application
helm install factorial-service ./infrastructure/helm \
  --namespace factorial \
  --values values-eks.yaml \
  --wait \
  --timeout 10m
```

### Step 10: Set Up Ingress/ALB

#### Option A: Using AWS Load Balancer Controller

```bash
# Install AWS Load Balancer Controller
helm repo add eks https://aws.github.io/eks-charts
helm install aws-load-balancer-controller eks/aws-load-balancer-controller \
  -n kube-system \
  --set clusterName=factorial-cluster \
  --set serviceAccount.create=false \
  --set serviceAccount.name=aws-load-balancer-controller

# Create Ingress resource
cat > ingress.yaml <<EOF
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: factorial-api-ingress
  namespace: factorial
  annotations:
    kubernetes.io/ingress.class: alb
    alb.ingress.kubernetes.io/scheme: internet-facing
    alb.ingress.kubernetes.io/target-type: ip
    alb.ingress.kubernetes.io/listen-ports: '[{"HTTP": 80}, {"HTTPS": 443}]'
    alb.ingress.kubernetes.io/ssl-redirect: '443'
spec:
  rules:
  - host: api.factorial.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: factorial-service-api
            port:
              number: 8080
EOF

kubectl apply -f ingress.yaml
```

#### Option B: Using NGINX Ingress

```bash
# Install NGINX Ingress
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm install ingress-nginx ingress-nginx/ingress-nginx \
  --namespace ingress-nginx \
  --create-namespace \
  --set controller.service.type=LoadBalancer

# Create Ingress (similar to above but with nginx annotations)
```

### Step 11: Verify Deployment

```bash
# Check all pods are running
kubectl get pods -n factorial

# Check services
kubectl get svc -n factorial

# Check ingress
kubectl get ingress -n factorial

# View logs
kubectl logs -n factorial -l component=api --tail=50
kubectl logs -n factorial -l component=worker --tail=50
kubectl logs -n factorial -l component=calculator --tail=50

# Test API endpoint
# Get ALB endpoint
ALB_ENDPOINT=$(kubectl get ingress factorial-api-ingress -n factorial -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')

# Test health endpoint
curl http://$ALB_ENDPOINT/health

# Or port-forward for testing
kubectl port-forward -n factorial svc/factorial-service-api 8080:8080
curl http://localhost:8080/health
```

### Step 12: Configure Auto-Scaling

HPA is already configured in the Helm chart. Monitor scaling:

```bash
# Check HPA status
kubectl get hpa -n factorial

# Watch HPA
watch kubectl get hpa -n factorial
```

### Step 13: Set Up Monitoring (Optional)

```bash
# Install Prometheus and Grafana
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm install prometheus prometheus-community/kube-prometheus-stack -n monitoring --create-namespace

# Configure service monitors for your application
```

### Step 14: Set Up Logging (Optional)

```bash
# Install Fluent Bit for CloudWatch Logs
helm repo add eks https://aws.github.io/eks-charts
helm install aws-for-fluent-bit eks/aws-for-fluent-bit \
  -n kube-system \
  --set cloudWatch.enabled=true \
  --set cloudWatch.region=us-east-1
```

### Troubleshooting EKS Deployment

```bash
# Check pod events
kubectl describe pod <pod-name> -n factorial

# Check node status
kubectl get nodes

# Check if IRSA is working
kubectl describe sa factorial-service-sa -n factorial

# Test S3 access from pod
kubectl exec -it -n factorial <pod-name> -- aws s3 ls s3://factorial-cal-services-prod/

# Check CloudWatch logs
aws logs tail /aws/eks/factorial-cluster/cluster --follow
```

### Production Recommendations

1. **Use Managed Services**: RDS, ElastiCache, Amazon MQ for better reliability
2. **Enable SSL/TLS**: Use SSL for all database and message queue connections
3. **Secrets Management**: Use AWS Secrets Manager instead of Kubernetes secrets
4. **Monitoring**: Set up CloudWatch, Prometheus, and Grafana
5. **Logging**: Use Fluent Bit to send logs to CloudWatch
6. **Backup**: Enable automated backups for RDS
7. **Disaster Recovery**: Set up multi-AZ deployments for critical services
8. **Cost Optimization**: Use Spot instances for worker nodes, Reserved instances for database
9. **Security**: Enable Pod Security Standards, network policies
10. **CI/CD**: Set up GitHub Actions or GitLab CI for automated deployments

### Cleanup

```bash
# Uninstall Helm release
helm uninstall factorial-service -n factorial

# Delete namespace
kubectl delete namespace factorial

# Delete ECR repositories
aws ecr delete-repository --repository-name factorial-api --force
aws ecr delete-repository --repository-name factorial-worker --force
aws ecr delete-repository --repository-name factorial-calculator --force

# Delete IAM role and policy
aws iam detach-role-policy --role-name factorial-service-role --policy-arn $POLICY_ARN
aws iam delete-role --role-name factorial-service-role
aws iam delete-policy --policy-arn $POLICY_ARN

# Delete S3 bucket
aws s3 rb s3://factorial-cal-services-prod --force

# Delete EKS cluster (if created with eksctl)
eksctl delete cluster --name factorial-cluster --region us-east-1
```

For more details, see `infrastructure/helm/README.md`.

## License

[Add your license here]
