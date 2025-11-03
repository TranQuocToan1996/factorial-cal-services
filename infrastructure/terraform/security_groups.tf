# express-nodejs-demo-pub (sg-029853ea536615004)
resource "aws_security_group" "pub" {
  name        = "express-nodejs-demo-pub"
  description = "express-nodejs-demo publish SG"
  vpc_id      = aws_vpc.main.id

  # Port 80 from admin CIDR
  ingress {
    description = "HTTP"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = [var.admin_cidr]
  }

  # Port 81 from admin CIDR
  ingress {
    description = "Port 81"
    from_port   = 81
    to_port     = 81
    protocol    = "tcp"
    cidr_blocks = [var.admin_cidr]
  }

  # Port 8080 from admin CIDR and self-reference
  ingress {
    description     = "App 8080"
    from_port       = 8080
    to_port         = 8080
    protocol        = "tcp"
    cidr_blocks     = [var.admin_cidr]
    security_groups = [aws_security_group.app.id]
  }

  # Port 3000 from admin CIDR and self-reference
  ingress {
    description     = "App 3000"
    from_port       = 3000
    to_port         = 3000
    protocol        = "tcp"
    cidr_blocks     = [var.admin_cidr]
    security_groups = [aws_security_group.app.id]
  }

  # Port 443 from admin CIDR, 0.0.0.0/0, and self-reference
  ingress {
    description     = "HTTPS"
    from_port       = 443
    to_port         = 443
    protocol        = "tcp"
    cidr_blocks     = [var.admin_cidr, "0.0.0.0/0"]
    security_groups = [aws_security_group.app.id]
  }

  # Port 8081 self-reference only
  ingress {
    description     = "App 8081"
    from_port       = 8081
    to_port         = 8081
    protocol        = "tcp"
    security_groups = [aws_security_group.app.id]
  }

  # Egress to 8080 and 3000
  egress {
    description = "App 8080"
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    description = "App 3000"
    from_port   = 3000
    to_port     = 3000
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(local.common_tags, {
    Name = "express-nodejs-demo-pub"
  })
}

# express-nodejs-demo-application (sg-0eda5df6af283ac2f)
resource "aws_security_group" "app" {
  name        = "express-nodejs-demo-application"
  description = "SG for application"
  vpc_id      = aws_vpc.main.id

  # Port 8080 from pub SG
  ingress {
    description     = "App 8080"
    from_port       = 8080
    to_port         = 8080
    protocol        = "tcp"
    security_groups = [aws_security_group.pub.id]
  }

  # Port 22 from dev SG
  ingress {
    description     = "SSH"
    from_port       = 22
    to_port         = 22
    protocol        = "tcp"
    security_groups = [aws_security_group.dev.id]
  }

  # Port 3000 from pub SG
  ingress {
    description     = "App 3000"
    from_port       = 3000
    to_port         = 3000
    protocol        = "tcp"
    security_groups = [aws_security_group.pub.id]
  }

  # Port 443 from pub SG
  ingress {
    description     = "HTTPS"
    from_port       = 443
    to_port         = 443
    protocol        = "tcp"
    security_groups = [aws_security_group.pub.id]
  }

  # Port 8081 from pub SG
  ingress {
    description     = "App 8081"
    from_port       = 8081
    to_port         = 8081
    protocol        = "tcp"
    security_groups = [aws_security_group.pub.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(local.common_tags, {
    Name = "express-nodejs-demo-application"
  })
}

# express-nodejs-demo-dev (sg-0826779e45be79bf9)
resource "aws_security_group" "dev" {
  name        = "express-nodejs-demo-dev"
  description = "dev instance"
  vpc_id      = aws_vpc.main.id

  # SSH from admin CIDR
  ingress {
    description = "SSH"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = [var.admin_cidr]
  }

  # No egress rules in JSON
  egress = []

  tags = merge(local.common_tags, {
    Name = "express-nodejs-demo-dev"
  })
}

# ec2-rds-1 (sg-08cc0e98b2464baa9)
resource "aws_security_group" "ec2_rds" {
  name        = "ec2-rds-1"
  description = "Security group attached to instances to securely connect to order-simple-postgres-db. Modification could lead to connection loss."
  vpc_id      = aws_vpc.main.id

  # No ingress rules

  # Egress to RDS SG on port 5432
  egress {
    description     = "Postgres 5432"
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [aws_security_group.rds_ec2.id]
  }

  tags = merge(local.common_tags, {
    Name = "ec2-rds-1"
  })
}

# express-nodejs-demo-postgres (sg-08deee793218845f6)
resource "aws_security_group" "postgres" {
  name        = "express-nodejs-demo-postgres"
  description = "Postgres"
  vpc_id      = aws_vpc.main.id

  # Port 5432 from admin CIDR and app SG
  ingress {
    description     = "Postgres 5432"
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    cidr_blocks     = [var.admin_cidr]
    security_groups = [aws_security_group.app.id]
  }

  # Port 22 from dev SG
  ingress {
    description     = "SSH"
    from_port       = 22
    to_port         = 22
    protocol        = "tcp"
    security_groups = [aws_security_group.dev.id]
  }

  # No egress rules in JSON
  egress = []

  tags = merge(local.common_tags, {
    Name = "express-nodejs-demo-postgres"
  })
}

# queue-stream-sg (sg-05c48ab7b3f4244be)
resource "aws_security_group" "rabbitmq" {
  name        = "queue-stream-sg"
  description = "queue-stream-sg"
  vpc_id      = aws_vpc.main.id

  # Port 5672 from admin CIDR and app SG
  ingress {
    description     = "AMQP 5672"
    from_port       = 5672
    to_port         = 5672
    protocol        = "tcp"
    cidr_blocks     = [var.admin_cidr]
    security_groups = [aws_security_group.app.id]
  }

  # Port 5671 from admin CIDR, app SG, and dev SG
  ingress {
    description     = "AMQPS 5671"
    from_port       = 5671
    to_port         = 5671
    protocol        = "tcp"
    cidr_blocks     = [var.admin_cidr]
    security_groups = [aws_security_group.app.id, aws_security_group.dev.id]
  }

  # Port 15672 from admin CIDR and app SG
  ingress {
    description     = "Mgmt 15672"
    from_port       = 15672
    to_port         = 15672
    protocol        = "tcp"
    cidr_blocks     = [var.admin_cidr]
    security_groups = [aws_security_group.app.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(local.common_tags, {
    Name = "queue-stream-sg"
  })
}

# rds-ec2-1 (sg-0ce4e5e1efef20f2d)
resource "aws_security_group" "rds_ec2" {
  name        = "rds-ec2-1"
  description = "Security group attached to order-simple-postgres-db to allow EC2 instances with specific security groups attached to connect to the database. Modification could lead to connection loss."
  vpc_id      = aws_vpc.main.id

  # Port 5432 from ec2-rds-1 SG
  ingress {
    description     = "Rule to allow connections from EC2 instances with sg-08cc0e98b2464baa9 attached"
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [aws_security_group.ec2_rds.id]
  }

  # No egress rules in JSON
  egress = []

  tags = merge(local.common_tags, {
    Name = "rds-ec2-1"
  })
}
