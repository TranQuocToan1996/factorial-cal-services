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

  vpc_id = "vpc-0956e9914efe15691"

  # Use specific private subnets for ElastiCache and MQ
  private_subnet_ids = [
    "subnet-0911715bbcfeaf63f", # express-nodejs-demo-subnet-private1-us-east-1a
    "subnet-0484b1c030d49a464"  # express-nodejs-demo-subnet-private2-us-east-1b
  ]
}


data "aws_subnets" "default" {
  filter {
    name   = "vpc-id"
    values = ["vpc-0956e9914efe15691"]
  }
}

# Redis module
module "redis" {
  source         = "./module/redis"
  cluster_name   = "infra"
  node_type      = "cache.t3.micro"
  num_nodes      = 1
  subnet_ids     = local.private_subnet_ids
  vpc_id         = "vpc-0956e9914efe15691"
  redis_host     = try(local.secrets.REDIS_HOST, "")
  redis_password = try(local.secrets.REDIS_PASSWORD, "")

  tags = {
    Environment = "Development"
    Project     = "FactorialService"
  }
}

# RabbitMQ module
module "rabbitmq" {
  source            = "./module/rabbitmq"
  cluster_name      = "infra"
  instance_type     = "mq.t3.micro"
  secret_arn        = data.aws_secretsmanager_secret.factorial_service.arn
  rabbitmq_user     = try(local.secrets.RABBITMQ_USER, "")
  rabbitmq_password = try(local.secrets.RABBITMQ_PASSWORD, "")
  subnet_ids        = data.aws_subnets.default.ids
  vpc_id            = "vpc-0956e9914efe15691"

  tags = {
    Environment = "Development"
    Project     = "FactorialService"
  }
}
