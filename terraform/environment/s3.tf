
module "data_bucket" {
  source                = "./modules/s3"
  bucket_name           = "report-data-${local.environment_name}"
  kms_key_id            = null
  custom_bucket_policy  = data.aws_iam_policy_document.allow_data_role_access
  access_logging_bucket = "s3-access-logs-opg-shared-development-eu-west-1"
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
