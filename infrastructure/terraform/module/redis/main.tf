resource "aws_elasticache_subnet_group" "redis_subnet_group" {
  name       = "${var.cluster_name}-redis-subnet-group"
  subnet_ids = var.subnet_ids
}

resource "aws_elasticache_cluster" "redis" {
  cluster_id           = "${var.cluster_name}-redis"
  engine               = "redis"
  node_type            = var.node_type
  num_cache_nodes      = var.num_nodes
  parameter_group_name = "default.redis7"
  subnet_group_name    = aws_elasticache_subnet_group.redis_subnet_group.name
  port                 = 6379
}
