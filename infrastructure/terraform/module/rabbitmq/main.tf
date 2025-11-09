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

resource "aws_security_group" "rabbitmq" {
  name        = "${var.cluster_name}-rabbitmq-sg"
  description = "Security group for RabbitMQ broker"
  vpc_id      = var.vpc_id

  ingress {
    description = "RabbitMQ AMQP"
    from_port   = 5672
    to_port     = 5672
    protocol    = "tcp"
    cidr_blocks = [data.aws_vpc.this.cidr_block]
  }

  ingress {
    description = "RabbitMQ Management"
    from_port   = 15672
    to_port     = 15672
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
    Name = "${var.cluster_name}-rabbitmq-sg"
  })
}
data "aws_vpc" "this" {
  id = var.vpc_id
}

resource "aws_mq_broker" "rabbitmq" {
  broker_name                = "${var.cluster_name}-rabbitmq"
  engine_type                = "RabbitMQ"
  engine_version             = "3.13"
  auto_minor_version_upgrade = true
  host_instance_type         = var.instance_type
  publicly_accessible        = false
  deployment_mode            = "SINGLE_INSTANCE"
  security_groups            = [aws_security_group.rabbitmq.id]
  subnet_ids                 = [var.subnet_ids[0]]

  user {
    username = local.username
    password = local.password
  }

  logs {
    general = true
  }

  tags = var.tags
}
