output "data_lake_bucket" {
  description = "The name of the data lake bucket"
  value       = aws_s3_bucket.data_lake.id
}

output "data_lake_arn" {
  description = "The ARN of the data lake bucket"
  value       = aws_s3_bucket.data_lake.arn
}

output "logs_bucket_name" {
  description = "The name of the access logs bucket"
  value       = aws_s3_bucket.access_logs.id
}

output "databricks_root_bucket" {
  description = "The name of the Databricks root bucket"
  value       = aws_s3_bucket.databricks_root.id
}

output "unity_catalog_bucket" {
  description = "The name of the Unity Catalog bucket"
  value       = aws_s3_bucket.unity_catalog.id
}
