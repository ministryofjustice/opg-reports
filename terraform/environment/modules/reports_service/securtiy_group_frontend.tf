resource "aws_security_group" "reports_frontend" {
  name                   = "opg-reports-frontend-${var.environment_name}"
  revoke_rules_on_delete = true
  vpc_id                 = data.aws_vpc.reports.id
  description            = "OPG Reports Frontend ECS Service"
  tags                   = { "Name" : "opg-reports-frontend-${var.environment_name}" }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_security_group_rule" "reports_frontend_egress" {
  type                     = "egress"
  protocol                 = "tcp"
  from_port                = 8080
  to_port                  = 8080
  source_security_group_id = aws_security_group.reports_api.id
  security_group_id        = aws_security_group.reports_frontend.id
  description              = "Allow Reports Frontend to connect to Reports API"
}

resource "aws_security_group_rule" "reports_frontend_alb_ingress" {
  type                     = "ingress"
  protocol                 = "tcp"
  from_port                = 80
  to_port                  = 80
  source_security_group_id = aws_security_group.reports_loadbalancer.id
  security_group_id        = aws_security_group.reports_frontend.id
  description              = "Ingress rule for Reports Frontend ECS Task from Load Balancer"
}
