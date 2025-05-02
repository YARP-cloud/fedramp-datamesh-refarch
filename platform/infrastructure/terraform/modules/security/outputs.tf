output "s3_kms_key_arn" {
  description = "ARN of KMS key for S3 encryption"
  value       = aws_kms_key.s3_key.arn
}

output "kafka_kms_key_arn" {
  description = "ARN of KMS key for Kafka encryption"
  value       = aws_kms_key.kafka_key.arn
}

output "databricks_kms_key_arn" {
  description = "ARN of KMS key for Databricks encryption"
  value       = aws_kms_key.databricks_key.arn
}

output "databricks_kms_key_alias" {
  description = "Alias of KMS key for Databricks encryption"
  value       = aws_kms_alias.databricks_key_alias.name
}

output "cloudtrail_bucket_name" {
  description = "Name of S3 bucket for CloudTrail logs"
  value       = aws_s3_bucket.cloudtrail_bucket.bucket
}

output "config_bucket_name" {
  description = "Name of S3 bucket for AWS Config logs"
  value       = aws_s3_bucket.config_bucket.bucket
}

output "databricks_role_arn" {
  description = "ARN of IAM role for Databricks"
  value       = aws_iam_role.databricks_role.arn
}

output "unity_catalog_role_arn" {
  description = "ARN of IAM role for Unity Catalog"
  value       = aws_iam_role.unity_catalog_role.arn
}
