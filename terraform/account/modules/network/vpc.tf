resource "aws_vpc" "main" {
  cidr_block           = var.cidr
  enable_dns_hostnames = true
  tags = merge(
    var.tags,
    { Name = "${var.tags.application}-vpc" },
  )
}

resource "aws_vpc_dhcp_options" "dns_resolver" {
  domain_name         = "${data.aws_region.current.name}.compute.internal"
  domain_name_servers = ["AmazonProvidedDNS"]
  tags = merge(
    var.tags,
    { Name = "${var.tags.application}-dns-resolver" },
  )
}

resource "aws_vpc_dhcp_options_association" "dns_resolver" {
  vpc_id          = aws_vpc.main.id
  dhcp_options_id = aws_vpc_dhcp_options.dns_resolver.id
}

resource "aws_internet_gateway" "gw" {
  vpc_id = aws_vpc.main.id
  tags = merge(
    var.tags,
    { Name = "${var.tags.application}-internet-gateway" },
  )
}

resource "aws_eip" "nat" {
  count  = 3
  domain = "vpc"
  tags = merge(
    var.tags,
    { Name = "${var.tags.application}-eip-nat-gateway" },
  )
}

resource "aws_nat_gateway" "gw" {
  count         = 3
  allocation_id = aws_eip.nat[count.index].id
  subnet_id     = aws_subnet.public[count.index].id
  tags = merge(
    var.tags,
    { Name = "${var.tags.application}-nat-gateway-${data.aws_availability_zones.all.names[count.index]}" },
  )
}

resource "aws_default_security_group" "default" {
  vpc_id = aws_vpc.main.id
  tags = merge(
    var.tags,
    { Name = "${var.tags.application}-vpc-default-security-group" }
  )
}
