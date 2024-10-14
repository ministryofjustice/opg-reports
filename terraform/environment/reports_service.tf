module "reports_service" {
  source = "./modules/reports_service"

  cloudwatch_log_group = aws_cloudwatch_log_group.reports
  dns_prefix           = local.config[local.environment_name]["dns_prefix"]
  environment_name     = local.environment_name
  tags                 = local.default_tags

  providers = {
    aws            = aws
    aws.management = aws.management
  }
}