resource "aws_vpc_endpoint" "s3" {
  count             = 3
  vpc_id            = module.network.vpc.id
  service_name      = "com.amazonaws.${data.aws_region.current.name}.s3"
  route_table_ids   = [data.aws_route_tables.public.ids[count.index], data.aws_route_tables.application.ids[count.index]]
  vpc_endpoint_type = "Gateway"
  policy            = data.aws_iam_policy_document.s3_vpc_endpoint.json
  tags = merge(
    local.default_tags,
    { Name = "${local.default_tags.application}-S3-endpoint" },
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

data "aws_route_tables" "application" {
  vpc_id = module.network.vpc.id

  filter {
    name   = "tag:Name"
    values = ["application-*"]
  }
}

data "aws_route_tables" "public" {
  vpc_id = module.network.vpc.id

  filter {
    name   = "tag:Name"
    values = ["public-*"]
  }
}

resource "aws_security_group" "vpc_endpoints_access" {
  name   = "${local.default_tags.application}-vpc-endpoint-access-subnets"
  vpc_id = module.network.vpc.id
  tags   = { Name = "${local.default_tags.application}-vpc-endpoint-access-subnets" }
}

resource "aws_security_group_rule" "vpc_endpoints_application_subnet_ingress" {
  from_port         = 443
  to_port           = 443
  protocol          = "tcp"
  security_group_id = aws_security_group.vpc_endpoints_access.id
  type              = "ingress"
  cidr_blocks       = module.network.application_subnets[*].cidr_block
  description       = "Allow Services in Application Subnets to connect to VPC Interface Endpoints"
}

resource "aws_security_group_rule" "vpc_endpoints_public_subnet_ingress" {
  from_port         = 443
  to_port           = 443
  protocol          = "tcp"
  security_group_id = aws_security_group.vpc_endpoints_access.id
  type              = "ingress"
  cidr_blocks       = module.network.public_subnets[*].cidr_block
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

resource "aws_vpc_endpoint" "application" {
  for_each = local.interface_endpoint

  vpc_id              = module.network.vpc.id
  service_name        = "com.amazonaws.${data.aws_region.current.name}.${each.value}"
  vpc_endpoint_type   = "Interface"
  private_dns_enabled = true
  security_group_ids  = [aws_security_group.vpc_endpoints_access.id]
  subnet_ids          = data.aws_subnets.application.ids
  tags = merge(
    local.default_tags,
    { Name = "${local.default_tags.application}-${each.value}-application" },
  )
}

data "aws_subnets" "application" {
  filter {
    name   = "vpc-id"
    values = [module.network.vpc.id]
  }

  filter {
    name   = "tag:Name"
    values = ["application-*"]
  }
}
