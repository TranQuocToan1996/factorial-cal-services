data "aws_secretsmanager_secret" "rabbitmq_credentials" {
  arn = var.secret_arn
}

data "aws_secretsmanager_secret_version" "rabbitmq_credentials" {
  secret_id = data.aws_secretsmanager_secret.rabbitmq_credentials.id
}

locals {
  # Use provided credentials or fallback to secret
  username = var.rabbitmq_user != "" ? var.rabbitmq_user : try(jsondecode(data.aws_secretsmanager_secret_version.rabbitmq_credentials.secret_string).RABBITMQ_USER, "admin")
  password = var.rabbitmq_password != "" ? var.rabbitmq_password : try(jsondecode(data.aws_secretsmanager_secret_version.rabbitmq_credentials.secret_string).RABBITMQ_PASSWORD, "")
}

resource "aws_mq_broker" "rabbitmq" {
  broker_name         = "${var.cluster_name}-rabbitmq"
  engine_type         = "RabbitMQ"
  engine_version      = "3.13.1"
  host_instance_type  = var.instance_type
  publicly_accessible = false
  deployment_mode     = "SINGLE_INSTANCE"

  user {
    username = local.username
    password = local.password
  }

  logs {
    general = true
  }

  tags = var.tags
}
