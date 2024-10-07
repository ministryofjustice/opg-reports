module "network" {
  source = "git@github.com:ministryofjustice/opg-terraform-aws-network.git?ref=v1.4.0"

  cidr                 = "10.1.0.0/16"
  enable_dns_hostnames = true

  providers = {
    aws = aws
  }
}
