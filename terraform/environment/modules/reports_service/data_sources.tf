data "aws_ecr_repository" "reports_api" {
  name     = "opg-reports/api"
  provider = aws.management
}

data "aws_ecr_repository" "reports_frontend" {
  name     = "opg-reports/front"
  provider = aws.management
}

data "aws_prefix_list" "s3" {
  name = "com.amazonaws.${data.aws_region.current.name}.s3"
}

data "aws_region" "current" {}

data "aws_security_group" "vpc_regional_endpoints" {
  vpc_id = data.aws_vpc.reports.id
  name   = "${var.tags.application}-vpc-endpoint-access-subnets"
}

data "aws_subnets" "private" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.reports.id]
  }

  filter {
    name   = "tag:Name"
    values = ["application-*"]
  }
}

data "aws_subnets" "public" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.reports.id]
  }

  filter {
    name   = "tag:Name"
    values = ["public-*"]
  }
}

data "aws_vpc" "reports" {
  filter {
    name   = "tag:Name"
    values = ["${local.name_prefix}-vpc"]
  }
}
