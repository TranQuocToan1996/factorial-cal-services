resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.this.id
  }

  tags = merge(local.common_tags, {
    Name = "express-nodejs-demo-rtb-public"
    name = "express-nodejs-demo-express-app"
  })
}

resource "aws_route_table_association" "public_a" {
  subnet_id      = aws_subnet.public_a.id
  route_table_id = aws_route_table.public.id
}

resource "aws_route_table_association" "public_b" {
  subnet_id      = aws_subnet.public_b.id
  route_table_id = aws_route_table.public.id
}

resource "aws_route_table" "private_a" {
  vpc_id = aws_vpc.main.id

  tags = merge(local.common_tags, {
    Name = "express-nodejs-demo-rtb-private1-us-east-1a"
    name = "express-nodejs-demo-express-app"
  })
}

resource "aws_route_table_association" "private_a" {
  subnet_id      = aws_subnet.private_a.id
  route_table_id = aws_route_table.private_a.id
}

resource "aws_route_table" "private_b" {
  vpc_id = aws_vpc.main.id

  tags = merge(local.common_tags, {
    Name = "express-nodejs-demo-rtb-private2-us-east-1b"
    name = "express-nodejs-demo-express-app"
  })
}

resource "aws_route_table_association" "private_b" {
  subnet_id      = aws_subnet.private_b.id
  route_table_id = aws_route_table.private_b.id
}

# Separate RDS route table
resource "aws_route_table" "rds" {
  vpc_id = aws_vpc.main.id

  tags = merge(local.common_tags, {
    Name = "RDS-Pvt-rt"
  })
}

# Associate all RDS subnets to the RDS route table
resource "aws_route_table_association" "rds" {
  for_each       = aws_subnet.rds
  subnet_id      = each.value.id
  route_table_id = aws_route_table.rds.id
}


