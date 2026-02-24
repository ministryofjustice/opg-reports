resource "aws_s3_bucket" "bucket" {
  bucket        = var.bucket_name
  force_destroy = var.force_destroy
}

resource "aws_s3_bucket_ownership_controls" "bucket_object_ownership" {
  bucket = aws_s3_bucket.bucket.id
  rule {
    object_ownership = "BucketOwnerEnforced"
  }
}

resource "aws_s3_bucket_versioning" "bucket_versioning" {
  bucket = aws_s3_bucket.bucket.id

  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "bucket_encryption_configuration" {
  bucket = aws_s3_bucket.bucket.bucket

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm     = "aws:kms"
      kms_master_key_id = var.kms_key_id
    }
  }
}

resource "aws_s3_bucket_lifecycle_configuration" "bucket" {
  count  = var.enable_lifecycle ? 1 : 0
  bucket = aws_s3_bucket.bucket.id

  dynamic "rule" {
    for_each = var.enable_lifecycle ? [1] : []

    content {
      id     = "ExpireObjectsAfter${var.expiration_days}Days"
      status = var.enable_lifecycle ? "Enabled" : "Disabled"

      expiration {
        days = var.expiration_days
      }

      noncurrent_version_expiration {
        noncurrent_days = var.non_current_expiration_days
      }
    }
  }
}

resource "aws_s3_bucket_public_access_block" "public_access_policy" {
  bucket = aws_s3_bucket.bucket.id

  block_public_acls       = var.block_public_acls
  block_public_policy     = var.block_public_policy
  ignore_public_acls      = var.ignore_public_acls
  restrict_public_buckets = var.restrict_public_buckets
}

resource "aws_s3_bucket_policy" "bucket" {
  depends_on = [aws_s3_bucket_public_access_block.public_access_policy]
  bucket     = aws_s3_bucket.bucket.id
  policy     = local.custom_bucket_policy ? data.aws_iam_policy_document.custom_bucket_policy[0].json : data.aws_iam_policy_document.bucket.json
}

resource "aws_s3_bucket_logging" "bucket" {
  count  = var.access_logging_bucket == "" ? 0 : 1
  bucket = aws_s3_bucket.bucket.id

  target_bucket = var.access_logging_bucket
  target_prefix = "log/${aws_s3_bucket.bucket.id}/"
}

data "aws_iam_policy_document" "custom_bucket_policy" {
  count = local.custom_bucket_policy ? 1 : 0
  source_policy_documents = [
    data.aws_iam_policy_document.bucket.json,
    var.custom_bucket_policy.json
  ]
}

data "aws_iam_policy_document" "bucket" {
  policy_id = "PutObjPolicy"

  statement {
    sid       = "DenyUnEncryptedObjectUploads"
    effect    = "Deny"
    actions   = ["s3:PutObject"]
    resources = ["${aws_s3_bucket.bucket.arn}/*"]

    condition {
      test     = "StringNotEquals"
      variable = "s3:x-amz-server-side-encryption"
      values   = ["AES256"]
    }

    principals {
      type        = "AWS"
      identifiers = ["*"]
    }
  }

  statement {
    sid     = "DenyNoneSSLRequests"
    effect  = "Deny"
    actions = ["s3:*"]
    resources = [
      aws_s3_bucket.bucket.arn,
      "${aws_s3_bucket.bucket.arn}/*"
    ]

    condition {
      test     = "Bool"
      variable = "aws:SecureTransport"
      values   = [false]
    }

    principals {
      type        = "AWS"
      identifiers = ["*"]
    }
  }
}
