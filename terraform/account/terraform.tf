terraform {

  backend "s3" {
    bucket         = "opg.terraform.state"
    key            = "opg-reports-account/terraform.tfstate"
    encrypt        = true
    region         = "eu-west-1"
    role_arn       = "arn:aws:iam::311462405659:role/opg-reporting-state-access"
    dynamodb_table = "remote_lock"
  }

}

provider "aws" {
  region = "eu-west-1"

  default_tags {
    tags = local.default_tags
  }
  assume_role {
    role_arn     = "arn:aws:iam::${local.environment.account_id}:role/${var.default_role}"
    session_name = "terraform-session"
  }
}
