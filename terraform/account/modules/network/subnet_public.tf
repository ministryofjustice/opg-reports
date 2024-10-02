resource "aws_subnet" "public" {
  count                           = 3
  vpc_id                          = aws_vpc.main.id
  cidr_block                      = cidrsubnet(aws_vpc.main.cidr_block, 7, count.index + 45)
  availability_zone               = data.aws_availability_zones.all.names[count.index]
  map_public_ip_on_launch         = false
  assign_ipv6_address_on_creation = false

  tags = merge(
    var.tags,
    {
      Name    = "${var.tags.application}-public-${data.aws_availability_zones.all.names[count.index]}"
      Private = "false"
    },
  )
}

resource "aws_route_table_association" "public" {
  count          = 3
  subnet_id      = aws_subnet.public[count.index].id
  route_table_id = aws_route_table.public[count.index].id
}

resource "aws_route_table" "public" {
  count  = 3
  vpc_id = aws_vpc.main.id

  tags = merge(
    var.tags,
    { Name = "${var.tags.application}-public" },
  )
}

resource "aws_route" "public_internet_gateway" {
  count = 3

  route_table_id         = aws_route_table.public[count.index].id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.gw.id

  timeouts {
    create = "5m"
  }
}
