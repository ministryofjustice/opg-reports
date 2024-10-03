resource "aws_security_group" "reports_api" {
  name_prefix            = "opg-reports-api-${var.environment_name}-"
  revoke_rules_on_delete = true
  vpc_id                 = data.aws_vpc.reports.id
  description            = "OPG Reports API ECS Service"
  tags                   = { "Name" : "opg-reports-api-${var.environment_name}" }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_security_group_rule" "reports_api_to_ecr_vpc_endpoint" {
  description              = "Reports API VPC Endpoint Access for ECR"
  type                     = "egress"
  protocol                 = "tcp"
  from_port                = 443
  to_port                  = 443
  security_group_id        = aws_security_group.reports_api.id
  source_security_group_id = data.aws_security_group.vpc_regional_endpoints.id
}

resource "aws_security_group_rule" "reports_api_to_s3_vpc_endpoint" {
  description       = "Reports API VPC Endpoint Access for S3 (For ECR Pull)"
  type              = "egress"
  protocol          = "tcp"
  from_port         = 443
  to_port           = 443
  security_group_id = aws_security_group.reports_api.id
  prefix_list_ids   = [data.aws_prefix_list.s3.id]
}


resource "aws_security_group_rule" "reports_frontend_ingress" {
  type                     = "ingress"
  protocol                 = "tcp"
  from_port                = 8081
  to_port                  = 8081
  source_security_group_id = aws_security_group.reports_frontend.id
  security_group_id        = aws_security_group.reports_api.id
  description              = "Allow Reports Frontend to connect to Reports API"
}
