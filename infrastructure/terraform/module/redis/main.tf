resource "aws_elasticache_subnet_group" "redis_subnet_group" {
  name       = "${var.cluster_name}-redis-subnet-group"
  subnet_ids = var.subnet_ids

  tags = var.tags
}

resource "aws_security_group" "redis" {
  name        = "${var.cluster_name}-redis-sg"
  description = "Security group for Redis ElastiCache cluster"
  vpc_id      = var.vpc_id

  ingress {
    description = "Redis from VPC"
    from_port   = 6379
    to_port     = 6379
    protocol    = "tcp"
    cidr_blocks = [data.aws_vpc.this.cidr_block]
  }

  egress {
    description = "Allow all outbound"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(var.tags, {
    Name = "${var.cluster_name}-redis-sg"
  })
}

data "aws_vpc" "this" {
  id = var.vpc_id
}

resource "aws_elasticache_replication_group" "redis" {
  replication_group_id       = "${var.cluster_name}-redis"
  description                = "Redis replication group for ${var.cluster_name}"
  engine                     = "redis"
  engine_version             = "7.0"
  node_type                  = var.node_type
  num_cache_clusters         = var.num_nodes
  parameter_group_name       = "default.redis7"
  automatic_failover_enabled = var.num_nodes > 1
  auth_token                 = var.redis_password != "" ? var.redis_password : null
  port                       = 6379
  subnet_group_name          = aws_elasticache_subnet_group.redis_subnet_group.name
  security_group_ids         = [aws_security_group.redis.id]
  transit_encryption_enabled = var.redis_password != "" ? true : false

  tags = var.tags
}

# resource "aws_elasticache_cluster" "redis" {
#   cluster_id           = "${var.cluster_name}-redis"
#   engine               = "redis"
#   node_type            = var.node_type
#   num_cache_nodes      = var.num_nodes
#   parameter_group_name = "default.redis7"
#   subnet_group_name    = aws_elasticache_subnet_group.redis_subnet_group.name
#   port                 = 6379
# }
