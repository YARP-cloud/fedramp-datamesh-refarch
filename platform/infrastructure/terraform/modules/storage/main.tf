# Storage infrastructure module for FedRAMP High Data Mesh

# Main S3 bucket for the data lake
resource "aws_s3_bucket" "data_lake" {
  bucket = "fedramp-data-mesh-lake-${data.aws_caller_identity.current.account_id}-${var.environment}"
  
  tags = {
    Name        = "fedramp-data-mesh-data-lake"
    Environment = var.environment
    Project     = "FedRAMP-Data-Mesh"
  }
}

# Enable versioning to track changes to objects
resource "aws_s3_bucket_versioning" "data_lake_versioning" {
  bucket = aws_s3_bucket.data_lake.id
  
  versioning_configuration {
    status = "Enabled"
  }
}

# Server-side encryption with KMS
resource "aws_s3_bucket_server_side_encryption_configuration" "data_lake_encryption" {
  bucket = aws_s3_bucket.data_lake.id
  
  rule {
    apply_server_side_encryption_by_default {
      kms_master_key_id = var.kms_key_arn
      sse_algorithm     = "aws:kms"
    }
  }
}

# Block public access
resource "aws_s3_bucket_public_access_block" "data_lake_public_access_block" {
  bucket = aws_s3_bucket.data_lake.id
  
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# S3 bucket for access logging
resource "aws_s3_bucket" "access_logs" {
  bucket = "fedramp-data-mesh-access-logs-${data.aws_caller_identity.current.account_id}-${var.environment}"
  
  tags = {
    Name        = "fedramp-data-mesh-access-logs"
    Environment = var.environment
    Project     = "FedRAMP-Data-Mesh"
  }
}

# Server-side encryption for access logs bucket
resource "aws_s3_bucket_server_side_encryption_configuration" "access_logs_encryption" {
  bucket = aws_s3_bucket.access_logs.id
  
  rule {
    apply_server_side_encryption_by_default {
      kms_master_key_id = var.kms_key_arn
      sse_algorithm     = "aws:kms"
    }
  }
}

# Block public access for access logs bucket
resource "aws_s3_bucket_public_access_block" "access_logs_public_access_block" {
  bucket = aws_s3_bucket.access_logs.id
  
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# Enable access logging for data lake bucket
resource "aws_s3_bucket_logging" "data_lake_logging" {
  bucket = aws_s3_bucket.data_lake.id
  
  target_bucket = aws_s3_bucket.access_logs.id
  target_prefix = "data-lake-access-logs/"
}

# S3 bucket for Databricks workspace
resource "aws_s3_bucket" "databricks_root" {
  bucket = "fedramp-data-mesh-databricks-root-${data.aws_caller_identity.current.account_id}-${var.environment}"
  
  tags = {
    Name        = "fedramp-data-mesh-databricks-root"
    Environment = var.environment
    Project     = "FedRAMP-Data-Mesh"
  }
}

# Server-side encryption for Databricks root bucket
resource "aws_s3_bucket_server_side_encryption_configuration" "databricks_root_encryption" {
  bucket = aws_s3_bucket.databricks_root.id
  
  rule {
    apply_server_side_encryption_by_default {
      kms_master_key_id = var.kms_key_arn
      sse_algorithm     = "aws:kms"
    }
  }
}

# Block public access for Databricks root bucket
resource "aws_s3_bucket_public_access_block" "databricks_root_public_access_block" {
  bucket = aws_s3_bucket.databricks_root.id
  
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# S3 bucket for Unity Catalog
resource "aws_s3_bucket" "unity_catalog" {
  bucket = "fedramp-data-mesh-unity-catalog-${data.aws_caller_identity.current.account_id}-${var.environment}"
  
  tags = {
    Name        = "fedramp-data-mesh-unity-catalog"
    Environment = var.environment
    Project     = "FedRAMP-Data-Mesh"
  }
}

# Server-side encryption for Unity Catalog bucket
resource "aws_s3_bucket_server_side_encryption_configuration" "unity_catalog_encryption" {
  bucket = aws_s3_bucket.unity_catalog.id
  
  rule {
    apply_server_side_encryption_by_default {
      kms_master_key_id = var.kms_key_arn
      sse_algorithm     = "aws:kms"
    }
  }
}

# Block public access for Unity Catalog bucket
resource "aws_s3_bucket_public_access_block" "unity_catalog_public_access_block" {
  bucket = aws_s3_bucket.unity_catalog.id
  
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# Create domain-specific prefixes in the data lake
resource "aws_s3_object" "domain_prefixes" {
  for_each = toset(var.domains)
  
  bucket  = aws_s3_bucket.data_lake.id
  key     = "${each.key}/"
  content = ""
  
  # Ensure the domain prefixes are encrypted
  server_side_encryption = "aws:kms"
  kms_key_id             = var.kms_key_arn
}

# Create data product hierarchy within each domain
resource "aws_s3_object" "data_product_prefixes" {
  for_each = { for entry in var.initial_data_products : "${entry.domain}/${entry.name}" => entry }
  
  bucket  = aws_s3_bucket.data_lake.id
  key     = each.key
  content = ""
  
  # Ensure the data product prefixes are encrypted
  server_side_encryption = "aws:kms"
  kms_key_id             = var.kms_key_arn
}

data "aws_caller_identity" "current" {}
