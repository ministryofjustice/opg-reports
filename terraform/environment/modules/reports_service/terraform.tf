terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "6.69.1"
      configuration_aliases = [
        aws.management
      ]
    }
  }
  required_version = "1.17.0"
}
