# # Key for encrypting data buckets
# resource "aws_kms_key" "s3_bucket_key" {
#   description             = "KMS Key for Encrypting S3 Buckets"
#   deletion_window_in_days = 14
#   enable_key_rotation     = true
# }

# resource "aws_kms_alias" "s3" {
#   name          = "alias/s3-bucket-report-data"
#   target_key_id = aws_kms_key.s3_bucket_key.key_id
# }

module "data_bucket" {
  source                = "./modules/s3"
  bucket_name           = "report-data-${local.environment_name}"
  kms_key_id            = null #aws_kms_key.s3_bucket_key.key_id
  custom_bucket_policy  = data.aws_iam_policy_document.allow_data_role_access
  access_logging_bucket = "s3-access-logs-report-data-${local.environment_name}-eu-west-1"
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
      ]
    }

    resources = [
      module.data_bucket.bucket.arn,
      "${module.data_bucket.bucket.arn}/*",
    ]
  }
}
