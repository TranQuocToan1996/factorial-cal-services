variable "cluster_name" {
  type        = string
  description = "Cluster name for RabbitMQ broker"
}

variable "instance_type" {
  type        = string
  default     = "mq.t3.micro"
  description = "RabbitMQ broker instance type"
}

variable "secret_arn" {
  type        = string
  description = "ARN of the Secrets Manager secret storing RabbitMQ credentials"
}

variable "rabbitmq_user" {
  type        = string
  description = "RabbitMQ username (from secret)"
  default     = ""
  sensitive   = true
}

variable "rabbitmq_password" {
  type        = string
  description = "RabbitMQ password (from secret)"
  default     = ""
  sensitive   = true
}

variable "subnet_ids" {
  type        = list(string)
  description = "List of subnet IDs for RabbitMQ broker"
}

variable "vpc_id" {
  type        = string
  description = "VPC ID for security group"
}

variable "tags" {
  type        = map(string)
  description = "Tags to apply to RabbitMQ broker"
  default     = {}
}
