variable "cidr" {
  type    = string
  default = "0.0.0.0/0"
}

variable "tags" {
  type = map(string)
}
