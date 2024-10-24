locals {
  dns_prefix = var.dns_prefix
  dns_suffix = "opg.service.justice.gov.uk"
  dns_name   = "${local.dns_prefix}.${local.dns_suffix}"
}
data "aws_route53_zone" "opg_service_justice_gov_uk" {
  provider = aws.management
  name     = local.dns_suffix
}

resource "aws_route53_record" "reports" {
  provider = aws.management
  zone_id  = data.aws_route53_zone.opg_service_justice_gov_uk.zone_id
  name     = local.dns_prefix
  type     = "A"

  alias {
    evaluate_target_health = false
    name                   = aws_lb.reports.dns_name
    zone_id                = aws_lb.reports.zone_id
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_acm_certificate" "reports" {
  domain_name       = local.dns_name
  validation_method = "DNS"

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_route53_record" "validation" {
  for_each = {
    for dvo in aws_acm_certificate.reports.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }
  allow_overwrite = true
  name            = each.value.name
  records         = [each.value.record]
  ttl             = 60
  type            = each.value.type
  provider        = aws.management
  zone_id         = data.aws_route53_zone.opg_service_justice_gov_uk.id
}

resource "aws_acm_certificate_validation" "reports" {
  certificate_arn         = aws_acm_certificate.reports.arn
  validation_record_fqdns = [for record in aws_route53_record.validation : record.fqdn]
  depends_on              = [aws_route53_record.validation]
}

resource "aws_service_discovery_private_dns_namespace" "reports" {
  name = "opg-reports.${var.environment_name}.ecs"
  vpc  = data.aws_vpc.reports.id
}
