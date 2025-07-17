
module "data_bucket" {
  source                = "./modules/s3"
  bucket_name           = "opg-reports-${local.environment_name}"
  kms_key_id            = null
  custom_bucket_policy  = data.aws_iam_policy_document.allow_data_role_access
  access_logging_bucket = "s3-access-logs-opg-shared-${local.environment_name}-eu-west-1"
  force_destroy         = true
  enable_lifecycle      = false
}

data "aws_iam_policy_document" "allow_data_role_access" {
  statement {
    sid = "AllowReportDataCIRoleBucketWrite"

    actions = [
      "s3:GetObject",
      "s3:PutObject",
      "s3:DeleteObject",
      "s3:ListBucket",
    ]

    effect = "Allow"

    principals {
      type = "AWS"
      identifiers = [
        "arn:aws:iam::${local.environment.account_id}:role/docs-and-metadata-ci",
        "arn:aws:iam::${local.environment.account_id}:role/operator",
      ]
    }

    resources = [
      module.data_bucket.bucket.arn,
      "${module.data_bucket.bucket.arn}/*",
    ]
  }
}

# OIDC access to the s3 bucket
module "reports_oidc_role_gha_ecr_push" {
  source      = "git@github.com:ministryofjustice/opg-terraform-aws-account//modules/github_oidc_roles?ref=v7.4.4"
  name        = "opg-reports-github-actions-s3"
  description = "A role for S3 permissions for GitHub Actions"

  permissions = [
    "repo:ministryofjustice/opg-reports:pull_request",
    "repo:ministryofjustice/opg-reports:ref:refs/heads/main",
    "repo:ministryofjustice/opg-reports:ref:refs/heads/*",
  ]

  custom_policy_documents = [data.aws_iam_policy_document.oidc_s3_access.json]

}

data "aws_iam_policy_document" "oidc_s3_access" {
  version = "2012-10-17"

  statement {
    sid    = "OIDCS3Access"
    effect = "Allow"
    actions = [
      "s3:GetObject",
      "s3:PutObject",
      "s3:DeleteObject",
      "s3:ListBucket",
    ]
    resources = [
      module.data_bucket.bucket.arn,
      "${module.data_bucket.bucket.arn}/*",
    ]
  }
  statement {
    effect    = "Allow"
    actions   = ["ecr:GetAuthorizationToken"]
    resources = ["*"]
  }
}
