# Factorial Service Helm Chart

This Helm chart deploys the Factorial Calculation Service on Kubernetes, including:
- **API Service**: REST API for submitting and retrieving factorial calculations
- **Worker Service**: Processes messages from RabbitMQ queue
- **Calculator Service**: Continuously calculates factorials sequentially

## Prerequisites

- Kubernetes 1.19+
- Helm 3.0+
- PostgreSQL database (or use provided postgres chart)
- Redis instance
- RabbitMQ instance
- AWS S3 bucket (for storing factorial results)
- AWS IAM role with S3 permissions (for EKS IRSA)

## Installation

### 1. Configure Values

Update `values.yaml` with your environment-specific values:

```yaml
database:
  host: "your-postgres-host"
  name: "factorial-cal-services"
  user: "postgres"

redis:
  host: "your-redis-host"

rabbitmq:
  host: "your-rabbitmq-host"

aws:
  region: "us-east-1"
  s3Bucket: "your-s3-bucket-name"

serviceAccount:
  annotations:
    eks.amazonaws.com/role-arn: "arn:aws:iam::ACCOUNT:role/YOUR_ROLE"
```

### 2. Create Secrets

Update the secret template or create secrets manually:

```bash
kubectl create secret generic factorial-service-secret \
  --from-literal=db-password='your-db-password' \
  --from-literal=rabbitmq-password='your-rabbitmq-password' \
  --from-literal=redis-password='your-redis-password'
```

### 3. Install Chart

```bash
helm install factorial-service ./infrastructure/helm \
  --namespace factorial \
  --create-namespace \
  --set database.host=your-postgres-host \
  --set redis.host=your-redis-host \
  --set rabbitmq.host=your-rabbitmq-host
```

## Configuration

### API Service

- **Replicas**: Default 2
- **Service Type**: LoadBalancer (configurable)
- **Port**: 8080
- **Health Check**: `/health` endpoint

### Worker Service

- **Replicas**: Default 3
- **Autoscaling**: Enabled by default (3-10 replicas)
- **Batch Processing**: Configurable via `WORKER_BATCH_SIZE` and `WORKER_MAX_BATCHES`

### Calculator Service

- **Replicas**: Default 1 (single instance recommended for sequential processing)
- **Autoscaling**: Disabled by default (can be enabled if needed)
- **Purpose**: Continuously calculates factorials in sequence

## Environment Variables

All services share common environment variables:
- Database connection settings
- Redis connection settings
- AWS S3 configuration
- Application configuration (max factorial, thresholds, etc.)

## Resource Limits

Default resource requests/limits:
- **API**: 256Mi-512Mi memory, 200m-500m CPU
- **Worker**: 512Mi-1Gi memory, 500m-1000m CPU
- **Calculator**: 512Mi-1Gi memory, 500m-1000m CPU

Adjust in `values.yaml` based on your workload.

## Autoscaling

Worker service supports Horizontal Pod Autoscaling:
- Enabled by default
- Scales based on CPU utilization (70% target)
- Min: 3 replicas, Max: 10 replicas

Calculator service HPA is disabled by default but can be enabled if needed.

## Upgrading

```bash
helm upgrade factorial-service ./infrastructure/helm \
  --namespace factorial \
  --set database.host=your-postgres-host
```

## Uninstallation

```bash
helm uninstall factorial-service --namespace factorial
```

## Troubleshooting

### Check Pod Status

```bash
kubectl get pods -n factorial -l app=factorial-service
```

### View Logs

```bash
# API logs
kubectl logs -n factorial -l component=api

# Worker logs
kubectl logs -n factorial -l component=worker

# Calculator logs
kubectl logs -n factorial -l component=calculator
```

### Check Services

```bash
kubectl get svc -n factorial
```

### Verify Environment Variables

```bash
kubectl describe pod <pod-name> -n factorial | grep -A 50 "Environment:"
```

## Notes

- Calculator service does NOT require RabbitMQ (only needs DB and Redis)
- Ensure AWS IAM role is properly configured for S3 access (IRSA for EKS)
- Database migrations run automatically on API service startup
- Secrets should be managed via external secret management (AWS Secrets Manager, Vault, etc.) in production

