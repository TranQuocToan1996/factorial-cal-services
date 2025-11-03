resource "aws_internet_gateway" "this" {
  vpc_id = aws_vpc.main.id

  tags = merge(local.common_tags, {
    Name = "express-nodejs-demo-igw"
  })
}


