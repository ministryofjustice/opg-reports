locals {
  key_admin = var.environment_name == "development" ? "operator" : "breakglass"
}

module "secrets_manager_kms" {
  source         = "git@github.com:ministryofjustice/opg-terraform-aws-kms-key.git?ref=v1.1.0"
  alias          = "/opg-reports/${var.environment_name}/secrets-manager-kms-key"
  description    = "Encryption keys for OPG Report ${var.environment_name} Secrets"
  primary_region = "eu-west-1"

  replicas_to_create = []

  administrator_roles = [
    "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/${local.key_admin}",
  ]
  decryption_roles = [
    "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/${local.key_admin}",
  ]
  encryption_roles = [
    "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/${local.key_admin}",
  ]
  usage_services = ["secretsmanager.*.amazonaws.com"]
}

resource "aws_secretsmanager_secret" "entra_tenant_id" {
  name                           = "opg-reports/${var.environment_name}/entra/tenant-id"
  description                    = "Entra Tenant ID for OIDC Integration in OPG Reports ${var.environment_name}"
  kms_key_id                     = module.secrets_manager_kms.primary_key.arn
  policy                         = data.aws_iam_policy_document.secrets_manager_access.json
  force_overwrite_replica_secret = true
  region                         = "eu-west-1"
}

data "aws_secretsmanager_secret_version" "entra_tenant_id" {
  secret_id = aws_secretsmanager_secret.entra_tenant_id.id
  region    = "eu-west-1"
}

removed {
  from = aws_secretsmanager_secret_version.entra_tenant_id
  lifecycle {
    destroy = false
  }
}

resource "aws_secretsmanager_secret" "entra_client_id" {
  name                           = "opg-reports/${var.environment_name}/entra/client-id"
  description                    = "Entra Client ID for OIDC Integration in OPG Reports ${var.environment_name}"
  kms_key_id                     = module.secrets_manager_kms.primary_key.arn
  policy                         = data.aws_iam_policy_document.secrets_manager_access.json
  force_overwrite_replica_secret = true
  region                         = "eu-west-1"
}

data "aws_secretsmanager_secret_version" "entra_client_id" {
  secret_id = aws_secretsmanager_secret.entra_client_id.id
  region    = "eu-west-1"
}

removed {
  from = aws_secretsmanager_secret_version.entra_client_id
  lifecycle {
    destroy = false
  }
}

resource "aws_secretsmanager_secret" "entra_client_secret" {
  name                           = "opg-reports/${var.environment_name}/entra/client-secret"
  description                    = "Entra Client Secret for OIDC Integration in OPG Reports ${var.environment_name}"
  kms_key_id                     = module.secrets_manager_kms.primary_key.arn
  policy                         = data.aws_iam_policy_document.secrets_manager_access.json
  force_overwrite_replica_secret = true
  region                         = "eu-west-1"
}

data "aws_secretsmanager_secret_version" "entra_client_secret" {
  secret_id = aws_secretsmanager_secret.entra_client_secret.id
  region    = "eu-west-1"
}

removed {
  from = aws_secretsmanager_secret_version.entra_client_secret
  lifecycle {
    destroy = false
  }
}

data "aws_iam_policy_document" "secrets_manager_access" {
  statement {
    effect  = "Allow"
    actions = ["secretsmanager:*"]

    resources = ["*"]
    principals {
      identifiers = var.environment_name == "development" ? [
        "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/operator",
        "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/breakglass",
        "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/docs-and-metadata-ci"
        ] : [
        "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/breakglass",
        "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/docs-and-metadata-ci"
      ]
      type = "AWS"
    }
  }
}
