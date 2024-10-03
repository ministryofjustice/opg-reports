resource "aws_security_group" "reports_frontend" {
  name_prefix                   = "opg-reports-frontend-${var.environment_name}-"
  revoke_rules_on_delete = true
  vpc_id                 = data.aws_vpc.reports.id
  description            = "OPG Reports Frontend ECS Service"
  tags                   = { "Name" : "opg-reports-frontend-${var.environment_name}" }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_security_group_rule" "reports_frontend_to_ecr_vpc_endpoint" {
  description              = "Reports Frontend VPC Endpoint Access for ECR"
  type                     = "egress"
  protocol                 = "tcp"
  from_port                = 443
  to_port                  = 443
  security_group_id        = aws_security_group.reports_frontend.id
  source_security_group_id = data.aws_security_group.vpc_regional_endpoints.id
}

resource "aws_security_group_rule" "reports_frontend_to_s3_vpc_endpoint" {
  description       = "Reports Frontend VPC Endpoint Access for S3 (For ECR Pull)"
  type              = "egress"
  protocol          = "tcp"
  from_port         = 443
  to_port           = 443
  security_group_id = aws_security_group.reports_frontend.id
  prefix_list_ids   = [data.aws_prefix_list.s3.id]
}

resource "aws_security_group_rule" "reports_frontend_egress" {
  type                     = "egress"
  protocol                 = "tcp"
  from_port                = 8081
  to_port                  = 8081
  source_security_group_id = aws_security_group.reports_api.id
  security_group_id        = aws_security_group.reports_frontend.id
  description              = "Allow Reports Frontend to connect to Reports API"
}

resource "aws_security_group_rule" "reports_frontend_alb_ingress" {
  type                     = "ingress"
  protocol                 = "tcp"
  from_port                = 8080
  to_port                  = 8080
  source_security_group_id = aws_security_group.reports_loadbalancer.id
  security_group_id        = aws_security_group.reports_frontend.id
  description              = "Ingress rule for Reports Frontend ECS Task from Load Balancer"
}

resource "aws_security_group_rule" "reports_frontend_outbound" {
  description       = "Allow all outbound traffic from Reports Frontend"
  type              = "egress"
  from_port         = 0
  to_port           = 0
  protocol          = "-1"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.reports_frontend.id
}
