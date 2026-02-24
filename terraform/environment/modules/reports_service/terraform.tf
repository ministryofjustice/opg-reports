terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "6.32.1"
      configuration_aliases = [
        aws.management
      ]
    }
  }
  required_version = "1.14.5"
}
