variable "cloudwatch_log_group" {
  type = object({
    name = string
  })
}

variable "dns_prefix" {
  type = string
}

variable "environment_name" {
  type = string
}

variable "reports_api_tag" {
  type    = string
  default = "latest"
}

variable "reports_frontend_tag" {
  type    = string
  default = "latest"
}

variable "s3_data_bucket" {
  type = object({
    arn = string
    id  = string
  })
  description = "The S3 bucket where all of the reports data is stored"
}

variable "tags" {
  type = map(string)
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


locals {
  name_prefix = "${var.tags.application}-${var.tags.environment-name}"
}
