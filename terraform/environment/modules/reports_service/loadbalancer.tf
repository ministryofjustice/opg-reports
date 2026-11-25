resource "aws_lb" "reports" {
  name_prefix        = "rep-"
  internal           = false
  load_balancer_type = "application"
  subnets            = data.aws_subnets.public.ids

  security_groups = [
    aws_security_group.reports_loadbalancer.id,
  ]

  tags = { Name = "opg-reports-alb-${var.environment_name}" }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_alb_target_group" "reports_frontend" {
  name_prefix          = "reptg-"
  port                 = 8080
  protocol             = "HTTP"
  target_type          = "ip"
  vpc_id               = data.aws_vpc.reports.id
  deregistration_delay = 0
  depends_on           = [aws_lb.reports]

  health_check {
    protocol            = "HTTP"
    path                = "/"
    interval            = 15
    timeout             = 10
    healthy_threshold   = 2
    unhealthy_threshold = 5
    matcher             = "200"
  }

  tags = { Name = "opg-reports-alb-tg-${var.environment_name}" }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_lb_listener" "https" {
  load_balancer_arn = aws_lb.reports.arn
  port              = "443"
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-TLS-1-2-Ext-2018-06"
  certificate_arn   = aws_acm_certificate_validation.reports.certificate_arn

  default_action {
    type = "fixed-response"

    fixed_response {
      content_type = "text/plain"
      message_body = "OPG Reports ${var.environment_name} ALB"
      status_code  = "200"
    }
  }
}

resource "aws_lb_listener_rule" "reports" {
  listener_arn = aws_lb_listener.https.arn
  priority     = 200

  action {
    type = "authenticate-oidc"

    authenticate_oidc {
      authorization_endpoint = "https://login.microsoftonline.com/${data.aws_secretsmanager_secret_version.entra_tenant_id.secret_string}/oauth2/v2.0/authorize"
      client_id              = data.aws_secretsmanager_secret_version.entra_client_id.secret_string
      client_secret          = data.aws_secretsmanager_secret_version.entra_client_secret.secret_string
      issuer                 = "https://login.microsoftonline.com/${data.aws_secretsmanager_secret_version.entra_tenant_id.secret_string}/v2.0"
      scope                  = "openid profile email"
      session_cookie_name    = "AWSELBAuthSessionCookie"
      session_timeout        = "21600"
      token_endpoint         = "https://login.microsoftonline.com/${data.aws_secretsmanager_secret_version.entra_tenant_id.secret_string}/oauth2/v2.0/token"
      user_info_endpoint     = "https://graph.microsoft.com/oidc/userinfo"
    }
  }

  action {
    type             = "forward"
    target_group_arn = aws_alb_target_group.reports_frontend.arn
  }
  condition {
    path_pattern {
      values = ["/*"]
    }
  }
}

resource "aws_lb_listener" "http" {
  load_balancer_arn = aws_lb.reports.arn
  port              = "80"
  protocol          = "HTTP"

  default_action {
    type = "redirect"

    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_302"
    }
  }
}

resource "aws_security_group" "reports_loadbalancer" {
  name_prefix = "opg-reports-lb-${var.environment_name}-"
  description = "Allow inbound traffic"
  vpc_id      = data.aws_vpc.reports.id

  lifecycle {
    create_before_destroy = true
  }
}

module "allow_list" {
  source = "git@github.com:ministryofjustice/opg-terraform-aws-moj-ip-allow-list.git?ref=v3.5.2"
}

resource "aws_security_group_rule" "loadbalancer_ingress_http" {
  type              = "ingress"
  from_port         = 80
  to_port           = 80
  protocol          = "tcp"
  cidr_blocks       = concat(module.allow_list.moj_global_protect_vpn, module.allow_list.palo_alto_prisma_access)
  security_group_id = aws_security_group.reports_loadbalancer.id
  description       = "Loadbalancer HTTP inbound from the MoJ VPN"
}

resource "aws_security_group_rule" "loadbalancer_ingress_https" {
  type              = "ingress"
  from_port         = 443
  to_port           = 443
  protocol          = "tcp"
  cidr_blocks       = concat(module.allow_list.moj_global_protect_vpn, module.allow_list.palo_alto_prisma_access)
  security_group_id = aws_security_group.reports_loadbalancer.id
  description       = "Loadbalancer HTTPS inbound from the MoJ VPN"
}

resource "aws_security_group_rule" "loadbalancer_egress" {
  type                     = "egress"
  protocol                 = "tcp"
  from_port                = 8080
  to_port                  = 8080
  security_group_id        = aws_security_group.reports_loadbalancer.id
  source_security_group_id = aws_security_group.reports_frontend.id
}

resource "aws_security_group_rule" "alb_global_egress_for_oidc" {
  description       = "Load Balancer Global Egress for OIDC"
  type              = "egress"
  protocol          = "tcp"
  from_port         = 443
  to_port           = 443
  security_group_id = aws_security_group.reports_loadbalancer.id
  cidr_blocks       = ["0.0.0.0/0"]
}
