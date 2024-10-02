resource "aws_lb" "reports" {
  name               = "opg-reports-${var.environment_name}"
  internal           = false
  load_balancer_type = "application"
  subnets            = data.aws_subnets.public.ids

  security_groups = [
    aws_security_group.reports_loadbalancer.id,
  ]

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_alb_target_group" "reports_frontend" {
  name                 = "reports-lb-tg"
  port                 = 8000
  protocol             = "HTTP"
  target_type          = "ip"
  vpc_id               = data.aws_vpc.reports.id
  deregistration_delay = 0
  depends_on           = [aws_lb.reports]

  health_check {
    protocol            = "HTTP"
    path                = "/overview/"
    interval            = 15
    timeout             = 10
    healthy_threshold   = 2
    unhealthy_threshold = 5
    matcher             = "200"
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_lb_listener" "https" {
  load_balancer_arn = aws_lb.reports.arn
  port              = "443"
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-TLS-1-2-Ext-2018-06"
  certificate_arn   = aws_acm_certificate_validation.response.certificate_arn

  default_action {
    type             = "forward"
    target_group_arn = aws_alb_target_group.reports_frontend.arn
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
  name        = "opg-reports-lb-${var.environment_name}"
  description = "Allow inbound traffic"
  vpc_id      = data.aws_vpc.reports.id

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_security_group_rule" "loadbalancer_ingress" {
  type              = "ingress"
  from_port         = 80
  to_port           = 80
  protocol          = "tcp"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.reports_loadbalancer.id
}

resource "aws_security_group_rule" "loadbalancer_egress" {
  type              = "egress"
  from_port         = 0
  to_port           = 0
  protocol          = "-1"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.reports_loadbalancer.id
}