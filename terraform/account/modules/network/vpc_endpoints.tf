resource "aws_vpc_endpoint" "s3" {
  count             = 3
  vpc_id            = aws_vpc.main.id
  service_name      = "com.amazonaws.${data.aws_region.current.name}.s3"
  route_table_ids   = [aws_route_table.public[count.index].id, aws_route_table.private[count.index].id]
  vpc_endpoint_type = "Gateway"
  policy            = data.aws_iam_policy_document.s3_vpc_endpoint.json
  tags = merge(
    var.tags,
    { Name = "${var.tags.application}-S3-endpoint" },
  )
}

data "aws_iam_policy_document" "s3_vpc_endpoint" {
  statement {
    sid       = "ReportsS3VpcEndpointPolicy"
    actions   = ["*"]
    resources = ["*"]
    principals {
      type        = "*"
      identifiers = ["*"]
    }
  }
}

resource "aws_security_group" "vpc_endpoints_private" {
  name   = "${var.tags.application}-vpc-endpoint-access-private-subnets"
  vpc_id = aws_vpc.main.id
  tags   = { Name = "${var.tags.application}-vpc-endpoint-access-private-subnets" }
}

resource "aws_security_group_rule" "vpc_endpoints_private_subnet_ingress" {
  from_port         = 443
  to_port           = 443
  protocol          = "tcp"
  security_group_id = aws_security_group.vpc_endpoints_private.id
  type              = "ingress"
  cidr_blocks       = aws_subnet.private[*].cidr_block
  description       = "Allow Services in Private Subnets to connect to VPC Interface Endpoints"
}

resource "aws_security_group_rule" "vpc_endpoints_public_subnet_ingress" {
  from_port         = 443
  to_port           = 443
  protocol          = "tcp"
  security_group_id = aws_security_group.vpc_endpoints_private.id
  type              = "ingress"
  cidr_blocks       = aws_subnet.public[*].cidr_block
  description       = "Allow Services in Public Subnets to connect to VPC Interface Endpoints"
}

locals {
  interface_endpoint = toset([
    "ecr.api",
    "ecr.dkr",
    "logs",
    "secretsmanager",
    "ssm"
  ])
}

resource "aws_vpc_endpoint" "private" {
  for_each = local.interface_endpoint

  vpc_id              = aws_vpc.main.id
  service_name        = "com.amazonaws.${data.aws_region.current.name}.${each.value}"
  vpc_endpoint_type   = "Interface"
  private_dns_enabled = true
  security_group_ids  = aws_security_group.vpc_endpoints_private[*].id
  subnet_ids          = data.aws_subnets.private.ids
  tags = merge(
    var.tags,
    { Name = "${var.tags.application}-${each.value}-private" },
  )
}
