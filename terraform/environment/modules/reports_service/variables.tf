variable "dns_prefix" {
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

variable "environment_name" {
  type = string
}

variable "cloudwatch_log_group" {
  type = object({
    name = string
  })
}