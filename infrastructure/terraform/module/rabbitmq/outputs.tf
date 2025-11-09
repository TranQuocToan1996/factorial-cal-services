output "rabbitmq_endpoints" {
  description = "RabbitMQ broker endpoints"
  value       = aws_mq_broker.rabbitmq.instances[0].endpoints
}

output "rabbitmq_broker_id" {
  description = "RabbitMQ broker ID"
  value       = aws_mq_broker.rabbitmq.id
}

output "rabbitmq_broker_arn" {
  description = "RabbitMQ broker ARN"
  value       = aws_mq_broker.rabbitmq.arn
}
