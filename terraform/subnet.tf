resource "aws_subnet" "public_subnet" {
  vpc_id     = aws_vpc.my_vpc.id
  cidr_block = "10.0.1.0/24"  # Altere conforme necessário
  availability_zone = "us-east-1a"  # Altere para a zona desejada

  tags = {
    Name = "public-subnet"
  }
}

resource "aws_route_table_association" "subnet_association" {
  subnet_id      = aws_subnet.public_subnet.id
  route_table_id = aws_route_table.my_route_table.id
}
