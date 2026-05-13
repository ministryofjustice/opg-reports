terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "6.43.0"
      configuration_aliases = [
        aws.management
      ]
    }
  }
  required_version = "1.15.2"
}
