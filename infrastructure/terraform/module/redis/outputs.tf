output "redis_endpoint" {
  description = "Redis configuration endpoint address"
  value       = aws_elasticache_replication_group.redis.configuration_endpoint_address
}

output "redis_port" {
  description = "Redis port"
  value       = aws_elasticache_replication_group.redis.port
}

output "redis_address" {
  description = "Redis address (from secret or cluster endpoint)"
  value       = var.redis_host != "" ? var.redis_host : aws_elasticache_replication_group.redis.configuration_endpoint_address
  sensitive   = false
}

output "redis_replication_group_id" {
  description = "Redis replication group ID"
  value       = aws_elasticache_replication_group.redis.replication_group_id
}
