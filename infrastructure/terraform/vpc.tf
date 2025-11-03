resource "aws_vpc" "main" {
  cidr_block           = var.vpc_cidr
  enable_dns_support   = true
  enable_dns_hostnames = true

  tags = merge(local.common_tags, {
    Name = "express-nodejs-demo-vpc"
    name = "express-nodejs-demo-express-app"
  })
}


