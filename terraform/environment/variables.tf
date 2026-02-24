locals {
  environment_name = lower(replace(terraform.workspace, "_", "-"))
  environment      = contains(keys(var.environments), local.environment_name) ? var.environments[local.environment_name] : var.environments["default"]

  config = {
    development = {
      dns_prefix = "dev.reports"
    }
    production = {
      dns_prefix = "reports"
    }
  }

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

  reports_api_tag      = var.api_image_tag != "" ? var.api_image_tag : local.environment.api_tag
  reports_frontend_tag = var.front_image_tag != "" ? var.front_image_tag : local.environment.front_tag

}

variable "api_image_tag" {
  type    = string
  default = ""
}
variable "front_image_tag" {
  type    = string
  default = ""
}

variable "default_role" {
  type    = string
  default = "docs-and-metadata-ci"
}


variable "management_role" {
  type    = string
  default = "docs-and-metadata-ci"
}

variable "semver_tag" {
  type        = string
  default     = "v0.0.0"
  description = "passed along for display and version tracking"
}

variable "commit_sha" {
  type        = string
  default     = "0000"
  description = "passed along to track versions"
}

variable "environments" {
  type = map(
    object({
      account_id   = string
      account_name = string
      api_tag      = string
      front_tag    = string
    })
  )
}
