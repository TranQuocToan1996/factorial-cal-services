variable "cluster_name" {
  type        = string
  description = "Cluster name for Redis"
}

variable "node_type" {
  type        = string
  default     = "cache.t3.micro"
  description = "ElastiCache node type"
}

variable "num_nodes" {
  type        = number
  default     = 1
  description = "Number of cache nodes"
}

variable "subnet_ids" {
  type        = list(string)
  description = "List of subnet IDs for ElastiCache subnet group"
}

variable "vpc_id" {
  type        = string
  description = "VPC ID for security group"
}

variable "redis_host" {
  type        = string
  description = "Redis host from Secrets Manager"
  default     = ""
}

variable "redis_password" {
  type        = string
  description = "Redis password/auth token from Secrets Manager"
  sensitive   = true
  default     = ""
}

variable "tags" {
  type        = map(string)
  description = "Tags to apply to resources"
  default     = {}
}
