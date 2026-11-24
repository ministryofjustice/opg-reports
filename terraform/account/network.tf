module "network" {
  source = "git@github.com:ministryofjustice/opg-terraform-aws-firewalled-network.git?ref=v1.5.0"

  cidr                                = "10.1.0.0/16"
  aws_networkfirewall_firewall_policy = local.firewall_config.shared_firewall_configuration != null ? null : aws_networkfirewall_firewall_policy.main
  default_security_group_ingress      = []
  default_security_group_egress       = []
  enable_dns_hostnames                = true
  network_firewall_enabled            = local.firewall_config.network_firewall_enabled
  shared_firewall_configuration       = local.firewall_config.shared_firewall_configuration
  providers = {
    aws = aws
  }
}


locals {
  firewall_config = lookup(local.environment_firewall_config, terraform.workspace, local.environment_firewall_config["production"])
  environment_firewall_config = {
    development = {
      network_firewall_enabled      = false
      none_matching_traffic_action  = "alert"
      shared_firewall_configuration = null
    }
    production = {
      network_firewall_enabled      = false
      none_matching_traffic_action  = "alert"
      shared_firewall_configuration = null
    }
  }
}

resource "aws_networkfirewall_firewall_policy" "main" {
  name = "opg-reports"

  firewall_policy {
    stateless_default_actions          = ["aws:forward_to_sfe"]
    stateless_fragment_default_actions = ["aws:forward_to_sfe"]

    stateful_engine_options {
      rule_order              = "DEFAULT_ACTION_ORDER"
      stream_exception_policy = "DROP"
    }
    stateful_rule_group_reference {
      resource_arn = aws_networkfirewall_rule_group.rule_file.arn
    }
  }
}

resource "aws_networkfirewall_rule_group" "rule_file" {
  capacity = 100
  name     = "opg-reports-${replace(filebase64sha256("${path.module}/network_firewall_rules.rules"), "/[^[:alnum:]]/", "")}"
  type     = "STATEFUL"
  rules    = templatefile("${path.module}/network_firewall_rules.rules", { action = local.firewall_config.none_matching_traffic_action })
  lifecycle {
    create_before_destroy = true
  }
}
