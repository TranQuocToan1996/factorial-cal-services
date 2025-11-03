resource "aws_subnet" "public_a" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = var.public_subnet_cidrs.a
  availability_zone       = var.availability_zones.a
  map_public_ip_on_launch = false

  tags = merge(local.common_tags, {
    Name = "express-nodejs-demo-subnet-public1-us-east-1a"
    name = "express-nodejs-demo-express-app"
  })
}

resource "aws_subnet" "public_b" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = var.public_subnet_cidrs.b
  availability_zone       = var.availability_zones.b
  map_public_ip_on_launch = false

  tags = merge(local.common_tags, {
    Name = "express-nodejs-demo-subnet-public2-us-east-1b"
    name = "express-nodejs-demo-express-app"
  })
}

resource "aws_subnet" "private_a" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = var.private_subnet_cidrs.a
  availability_zone = var.availability_zones.a

  tags = merge(local.common_tags, {
    Name = "express-nodejs-demo-subnet-private1-us-east-1a"
    name = "express-nodejs-demo-express-app"
  })
}

resource "aws_subnet" "private_b" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = var.private_subnet_cidrs.b
  availability_zone = var.availability_zones.b

  tags = merge(local.common_tags, {
    Name = "express-nodejs-demo-subnet-private2-us-east-1b"
    name = "express-nodejs-demo-express-app"
  })
}

# RDS private subnets with proper AZ mapping
resource "aws_subnet" "rds" {
  for_each          = var.rds_subnet_az_mapping
  vpc_id            = aws_vpc.main.id
  cidr_block        = each.key
  availability_zone = each.value

  tags = merge(local.common_tags, {
    Name = local.rds_subnet_names[each.key]
  })
}


