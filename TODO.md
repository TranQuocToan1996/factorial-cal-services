1. README
2. clean unit test cicd

## History

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