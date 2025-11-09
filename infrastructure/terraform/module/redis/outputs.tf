output "redis_endpoint" {
  description = "Redis configuration endpoint"
  value       = aws_elasticache_cluster.redis.configuration_endpoint
}

output "redis_address" {
  description = "Redis address (from secret or cluster endpoint)"
  value       = var.redis_host != "" ? var.redis_host : aws_elasticache_cluster.redis.configuration_endpoint
  sensitive   = false
}
