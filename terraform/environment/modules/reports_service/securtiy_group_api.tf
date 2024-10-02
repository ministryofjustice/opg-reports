resource "aws_security_group" "reports_api" {
  name                   = "opg-reports-api-${var.environment_name}"
  revoke_rules_on_delete = true
  vpc_id                 = data.aws_vpc.reports.id
  description            = "OPG Reports API ECS Service"
  tags                   = { "Name" : "opg-reports-api-${var.environment_name}" }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_security_group_rule" "reports_frontend_ingress" {
  type                     = "ingress"
  protocol                 = "tcp"
  from_port                = 8080
  to_port                  = 8080
  source_security_group_id = aws_security_group.reports_frontend.id
  security_group_id        = aws_security_group.reports_api.id
  description              = "Allow Reports Frontend to connect to Reports API"
}
