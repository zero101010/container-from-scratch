resource "aws_route_table" "my_route_table" {
  vpc_id = aws_vpc.my_vpc.id

  tags = {
    Name = "my-route-table"
  }
}

resource "aws_route" "route_to_igw" {
  route_table_id         = aws_route_table.my_route_table.id
  destination_cidr_block = "0.0.0.0/0"  # Rota para todo o tráfego

  gateway_id = aws_internet_gateway.my_igw.id
}