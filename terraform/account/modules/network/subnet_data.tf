resource "aws_subnet" "data" {
  count                           = 3
  vpc_id                          = aws_vpc.main.id
  cidr_block                      = cidrsubnet(aws_vpc.main.cidr_block, 7, count.index + 115)
  availability_zone               = data.aws_availability_zones.all.names[count.index]
  map_public_ip_on_launch         = false
  assign_ipv6_address_on_creation = false

  tags = merge(
    var.tags,
    { 
      Name = "${var.tags.application}-data-${data.aws_availability_zones.all.names[count.index]}"
      Private = "true" 
    },
  )
}

resource "aws_route_table_association" "data" {
  count          = 3
  subnet_id      = aws_subnet.data[count.index].id
  route_table_id = aws_route_table.data[count.index].id
}

resource "aws_route_table" "data" {
  count  = 3
  vpc_id = aws_vpc.main.id
  
  tags = merge(
    var.tags,
    { Name = "${var.tags.application}-data-route-table" },
  )
}
