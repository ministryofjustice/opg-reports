terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.0.0"
      configuration_aliases = [
        aws.management
      ]
    }
  }
  required_version = ">= 1.0.0"
}
