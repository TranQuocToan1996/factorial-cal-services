output "vpc_id" {
  value       = aws_vpc.main.id
  description = "VPC ID"
}

output "igw_id" {
  value       = aws_internet_gateway.this.id
  description = "Internet Gateway ID"
}

output "public_subnet_ids" {
  value       = [aws_subnet.public_a.id, aws_subnet.public_b.id]
  description = "Public subnet IDs"
}

output "private_subnet_ids" {
  value       = [aws_subnet.private_a.id, aws_subnet.private_b.id]
  description = "Private subnet IDs"
}

output "rds_subnet_ids" {
  value       = [for s in aws_subnet.rds : s.id]
  description = "RDS private subnet IDs"
}

output "security_groups" {
  value = {
    pub      = aws_security_group.pub.id
    app      = aws_security_group.app.id
    dev      = aws_security_group.dev.id
    ec2_rds  = aws_security_group.ec2_rds.id
    postgres = aws_security_group.postgres.id
    rabbitmq = aws_security_group.rabbitmq.id
    rds_ec2  = aws_security_group.rds_ec2.id
  }
  description = "Security group IDs"
}

output "route_table_ids" {
  value = {
    public    = aws_route_table.public.id
    private_a = aws_route_table.private_a.id
    private_b = aws_route_table.private_b.id
    rds       = aws_route_table.rds.id
  }
  description = "Route table IDs"
}


