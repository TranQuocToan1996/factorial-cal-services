terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
  required_version = ">= 1.5.0"
}

provider "aws" {
  region = "us-east-1"
}

# Data source to read existing secret
data "aws_secretsmanager_secret" "factorial_service" {
  arn = "arn:aws:secretsmanager:us-east-1:218435950768:secret:dev/factorial-service-PuUw9d"
}

data "aws_secretsmanager_secret_version" "factorial_service" {
  secret_id = data.aws_secretsmanager_secret.factorial_service.id
}

locals {
  secrets = jsondecode(data.aws_secretsmanager_secret_version.factorial_service.secret_string)
}

# Redis module
module "redis" {
  source       = "./module/redis"
  cluster_name = "infra"
  node_type    = "cache.t3.micro"
  num_nodes    = 1
  # Redis connection info from secret (if needed for outputs)
  redis_host     = local.secrets.REDIS_HOST
  redis_password = local.secrets.REDIS_PASSWORD
}

# RabbitMQ module
module "rabbitmq" {
  source            = "./module/rabbitmq"
  cluster_name      = "infra"
  instance_type     = "mq.t3.micro"
  secret_arn        = data.aws_secretsmanager_secret.factorial_service.arn
  rabbitmq_user     = local.secrets.RABBITMQ_USER
  rabbitmq_password = local.secrets.RABBITMQ_PASSWORD
}
