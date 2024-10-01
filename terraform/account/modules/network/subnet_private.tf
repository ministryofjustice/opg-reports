resource "aws_subnet" "private" {
  count                           = 3
  vpc_id                          = aws_vpc.main.id
  cidr_block                      = cidrsubnet(aws_vpc.main.cidr_block, 7, count.index + 95)
  availability_zone               = data.aws_availability_zones.all.names[count.index]
  map_public_ip_on_launch         = false
  assign_ipv6_address_on_creation = false

  tags = merge(
    var.tags,
    { 
      Name = "${var.tags.application}-private-${data.aws_availability_zones.all.names[count.index]}"
      Private = "true" 
    },
  )
}

resource "aws_route_table_association" "private" {
  count          = 3
  subnet_id      = aws_subnet.private[count.index].id
  route_table_id = aws_route_table.private[count.index].id
}

resource "aws_route_table" "private" {
  count  = 3
  vpc_id = aws_vpc.main.id

  tags = merge(
    var.tags,
    { Name = "${var.tags.application}-private-route-table" },
  )
}

resource "aws_route" "private_nat_gateway" {
  count = 3

  route_table_id         = aws_route_table.private[count.index].id
  destination_cidr_block = "0.0.0.0/0"
  nat_gateway_id         = aws_nat_gateway.gw[count.index].id

  timeouts {
    create = "5m"
  }
}
