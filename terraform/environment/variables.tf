locals {
  environment       = contains(keys(var.environments), terraform.workspace) ? var.environments[terraform.workspace] : var.environments["default"]
  environment_name  = terraform.workspace == "production" ? "production" : "development"

  mandatory_moj_tags = {
    business-unit    = "OPG"
    application      = "opg-reports"
    account          = local.environment.account_name
    environment-name = local.environment_name
    is-production    = tostring(terraform.workspace == "production" ? true : false)
    owner            = "OPG WebOps: opgteam@digital.justice.gov.uk"
  }

  optional_tags = {
    source-code            = "https://github.com/ministryofjustice/opg-reports"
    infrastructure-support = "OPG Webops: opgteam@digital.justice.gov.uk"
    terraform-managed      = "Managed by Terraform"
  }

  default_tags = merge(local.mandatory_moj_tags, local.optional_tags)

}

variable "default_role" {
  type    = string
  default = "docs-and-metadata-ci"
}

variable "management_role" {
  type    = string
  default = "docs-and-metadata-ci"
}

variable "environments" {
  type = map(
    object({
      account_id     = string
      account_name   = string
    })
  )
}
