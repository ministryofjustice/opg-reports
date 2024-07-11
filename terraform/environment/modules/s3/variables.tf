variable "bucket_name" {
  description = "Name of the bucket."
  type        = string
}

variable "force_destroy" {
  description = "A boolean that indicates all objects should be deleted from the bucket so that the bucket can be destroyed without error. These objects are not recoverable."
  type        = bool
  default     = false
}

variable "block_public_acls" {
  description = "Whether Amazon S3 should block public ACLs for this bucket."
  type        = bool
  default     = true
}

variable "block_public_policy" {
  description = "Whether Amazon S3 should block public bucket policies for this bucket."
  type        = bool
  default     = true
}

variable "ignore_public_acls" {
  description = "Whether Amazon S3 should ignore public ACLs for this bucket."
  type        = bool
  default     = true
}

variable "restrict_public_buckets" {
  description = "Whether Amazon S3 should restrict public bucket policies for this bucket."
  type        = bool
  default     = true
}

variable "kms_key_id" {
  type        = string
  description = "KMS key to encrypt s3 bucket with"
}

variable "enable_lifecycle" {
  description = "Delete items in the bucket after 1 year if enabled."
  type        = bool
  default     = false
}

variable "expiration_days" {
  description = "Number of days to expire the items in the bucket. Only takes effect when enable_lifecycle is set to true."
  type        = string
  default     = "365"
}

variable "non_current_expiration_days" {
  description = "Lifecycle expiration days for non current version"
  type        = string
  default     = "14"
}

variable "access_logging_bucket" {
  description = "Name of bucket used for logging access to the bucket"
  type        = string
  default     = ""
}

variable "custom_bucket_policy" {
  description = "Bucket policy iam_object that gets merged with our base policy before being applied."
  default     = null
  type = object({
    json = string
  })
}

locals {
  custom_bucket_policy = var.custom_bucket_policy != null ? true : false
}
