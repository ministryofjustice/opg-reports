data "aws_availability_zones" "all" {
  state = "available"
}

data "aws_region" "current" {}

