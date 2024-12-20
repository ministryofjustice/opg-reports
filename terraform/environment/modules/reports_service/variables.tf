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

variable "tags" {
  type = map(string)
}

locals {
  name_prefix = "${var.tags.application}-${var.tags.environment-name}"
}