data "aws_region" "current" {}

data "aws_ecr_repository" "reports_api" {
  name     = "opg-reports/api"
  provider = aws.management
}

data "aws_ecr_repository" "reports_frontend" {
  name     = "opg-reports/front"
  provider = aws.management
}

data "aws_vpc" "reports" {
  filter {
    name   = "tag:Name"
    values = ["${var.tags.application}-vpc"]
  }
}

data "aws_subnets" "private" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.reports.id]
  }

  filter {
    name   = "tag:Name"
    values = ["${var.tags.application}-private-*"]
  }
}

data "aws_subnets" "public" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.reports.id]
  }

  filter {
    name   = "tag:Name"
    values = ["${var.tags.application}-public-*"]
  }
}