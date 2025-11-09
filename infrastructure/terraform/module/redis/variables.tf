variable "cluster_name" {
  type = string
}

variable "node_type" {
  type    = string
  default = "cache.t3.micro"
}

variable "num_nodes" {
  type    = number
  default = 1
}

variable "subnet_ids" {
  type    = list(string)
  default = []
}

variable "redis_host" {
  type        = string
  description = "Redis host from Secrets Manager"
  default     = ""
}

variable "redis_password" {
  type        = string
  description = "Redis password from Secrets Manager"
  sensitive   = true
  default     = ""
}
