locals {
  environment_name  = lower(replace(terraform.workspace, "_", "-"))
  environment       = contains(keys(var.environments), local.environment_name) ? var.environments[local.environment_name] : var.environments["default"]

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

variable "ci_role" {
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
