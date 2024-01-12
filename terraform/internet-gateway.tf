resource "aws_internet_gateway" "my_igw" {
  vpc_id = aws_vpc.my_vpc.id  # Substitua pelo ID da sua VPC

  tags = {
    Name = "my-internet-gateway"
  }
}

