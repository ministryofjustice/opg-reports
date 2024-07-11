terraform {

  backend "s3" {
    bucket         = "opg.terraform.state"
    key            = "opg-reports-account/terraform.tfstate"
    encrypt        = true
    region         = "eu-west-1"
    role_arn       = "arn:aws:iam::311462405659:role/docs-and-metadata-ci"
    dynamodb_table = "remote_lock"
  }

}
