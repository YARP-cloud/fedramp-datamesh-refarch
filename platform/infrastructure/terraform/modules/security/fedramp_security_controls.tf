# FedRAMP High security controls for AWS infrastructure

# KMS Customer Managed Keys for encryption
resource "aws_kms_key" "s3_key" {
  description             = "KMS key for S3 bucket encryption"
  deletion_window_in_days = 30
  enable_key_rotation     = true
  policy                  = data.aws_iam_policy_document.s3_key_policy.json
  
  tags = {
    Name        = "fedramp-data-mesh-s3-key"
    Environment = var.environment
    Project     = "FedRAMP-Data-Mesh"
  }
}

resource "aws_kms_alias" "s3_key_alias" {
  name          = "alias/fedramp-data-mesh-s3-${var.environment}"
  target_key_id = aws_kms_key.s3_key.key_id
}

# KMS Key for Kafka encryption
resource "aws_kms_key" "kafka_key" {
  description             = "KMS key for Kafka encryption"
  deletion_window_in_days = 30
  enable_key_rotation     = true
  policy                  = data.aws_iam_policy_document.kafka_key_policy.json
  
  tags = {
    Name        = "fedramp-data-mesh-kafka-key"
    Environment = var.environment
    Project     = "FedRAMP-Data-Mesh"
  }
}

resource "aws_kms_alias" "kafka_key_alias" {
  name          = "alias/fedramp-data-mesh-kafka-${var.environment}"
  target_key_id = aws_kms_key.kafka_key.key_id
}

# KMS Key for Databricks encryption
resource "aws_kms_key" "databricks_key" {
  description             = "KMS key for Databricks encryption"
  deletion_window_in_days = 30
  enable_key_rotation     = true
  policy                  = data.aws_iam_policy_document.databricks_key_policy.json
  
  tags = {
    Name        = "fedramp-data-mesh-databricks-key"
    Environment = var.environment
    Project     = "FedRAMP-Data-Mesh"
  }
}

resource "aws_kms_alias" "databricks_key_alias" {
  name          = "alias/fedramp-data-mesh-databricks-${var.environment}"
  target_key_id = aws_kms_key.databricks_key.key_id
}

# KMS Key Policy for S3
data "aws_iam_policy_document" "s3_key_policy" {
  statement {
    sid       = "Enable IAM User Permissions"
    effect    = "Allow"
    principals {
      type        = "AWS"
      identifiers = ["arn:aws:iam::${data.aws_caller_identity.current.account_id}:root"]
    }
    actions   = ["kms:*"]
    resources = ["*"]
  }
  
  statement {
    sid       = "Allow S3 Service to use the key"
    effect    = "Allow"
    principals {
      type        = "Service"
      identifiers = ["s3.amazonaws.com"]
    }
    actions   = [
      "kms:Encrypt",
      "kms:Decrypt",
      "kms:ReEncrypt*",
      "kms:GenerateDataKey*",
      "kms:DescribeKey"
    ]
    resources = ["*"]
  }
  
  statement {
    sid       = "Allow authorized roles to use the key"
    effect    = "Allow"
    principals {
      type        = "AWS"
      identifiers = var.allowed_role_arns
    }
    actions   = [
      "kms:Encrypt",
      "kms:Decrypt",
      "kms:ReEncrypt*",
      "kms:GenerateDataKey*",
      "kms:DescribeKey"
    ]
    resources = ["*"]
  }
}

# KMS Key Policy for Kafka
data "aws_iam_policy_document" "kafka_key_policy" {
  statement {
    sid       = "Enable IAM User Permissions"
    effect    = "Allow"
    principals {
      type        = "AWS"
      identifiers = ["arn:aws:iam::${data.aws_caller_identity.current.account_id}:root"]
    }
    actions   = ["kms:*"]
    resources = ["*"]
  }
  
  statement {
    sid       = "Allow Kafka Service to use the key"
    effect    = "Allow"
    principals {
      type        = "Service"
      identifiers = ["kafka.amazonaws.com"]
    }
    actions   = [
      "kms:Encrypt",
      "kms:Decrypt",
      "kms:ReEncrypt*",
      "kms:GenerateDataKey*",
      "kms:DescribeKey"
    ]
    resources = ["*"]
  }
  
  statement {
    sid       = "Allow authorized roles to use the key"
    effect    = "Allow"
    principals {
      type        = "AWS"
      identifiers = var.allowed_role_arns
    }
    actions   = [
      "kms:Encrypt",
      "kms:Decrypt",
      "kms:ReEncrypt*",
      "kms:GenerateDataKey*",
      "kms:DescribeKey"
    ]
    resources = ["*"]
  }
}

# KMS Key Policy for Databricks
data "aws_iam_policy_document" "databricks_key_policy" {
  statement {
    sid       = "Enable IAM User Permissions"
    effect    = "Allow"
    principals {
      type        = "AWS"
      identifiers = ["arn:aws:iam::${data.aws_caller_identity.current.account_id}:root"]
    }
    actions   = ["kms:*"]
    resources = ["*"]
  }
  
  statement {
    sid       = "Allow authorized roles to use the key"
    effect    = "Allow"
    principals {
      type        = "AWS"
      identifiers = var.allowed_role_arns
    }
    actions   = [
      "kms:Encrypt",
      "kms:Decrypt",
      "kms:ReEncrypt*",
      "kms:GenerateDataKey*",
      "kms:DescribeKey"
    ]
    resources = ["*"]
  }
  
  statement {
    sid       = "Allow Databricks service principals to use the key"
    effect    = "Allow"
    principals {
      type        = "AWS"
      identifiers = [aws_iam_role.databricks_role.arn]
    }
    actions   = [
      "kms:Encrypt",
      "kms:Decrypt",
      "kms:ReEncrypt*",
      "kms:GenerateDataKey*",
      "kms:DescribeKey"
    ]
    resources = ["*"]
  }
}

# IAM Role for Databricks
resource "aws_iam_role" "databricks_role" {
  name = "fedramp-data-mesh-databricks-role-${var.environment}"
  
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "databricks.amazonaws.com"
        }
      }
    ]
  })
}

# IAM Role for Unity Catalog
resource "aws_iam_role" "unity_catalog_role" {
  name = "fedramp-data-mesh-unity-catalog-role-${var.environment}"
  
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "databricks.amazonaws.com"
        }
      }
    ]
  })
}

# AWS Config for FedRAMP compliance monitoring
resource "aws_config_configuration_recorder" "fedramp" {
  name     = "fedramp-data-mesh-recorder-${var.environment}"
  role_arn = aws_iam_role.config_role.arn
  
  recording_group {
    all_supported                 = true
    include_global_resource_types = true
  }
}

resource "aws_config_delivery_channel" "fedramp" {
  name           = "fedramp-data-mesh-delivery-channel-${var.environment}"
  s3_bucket_name = aws_s3_bucket.config_bucket.bucket
  
  snapshot_delivery_properties {
    delivery_frequency = "Six_Hours"
  }
  
  depends_on = [aws_config_configuration_recorder.fedramp]
}

resource "aws_config_configuration_recorder_status" "fedramp" {
  name       = aws_config_configuration_recorder.fedramp.name
  is_enabled = true
  
  depends_on = [aws_config_delivery_channel.fedramp]
}

# S3 bucket for AWS Config logs
resource "aws_s3_bucket" "config_bucket" {
  bucket = "fedramp-data-mesh-config-${data.aws_caller_identity.current.account_id}-${var.environment}"
  
  tags = {
    Name        = "fedramp-data-mesh-config-bucket"
    Environment = var.environment
    Project     = "FedRAMP-Data-Mesh"
  }
}

# S3 bucket server-side encryption 
resource "aws_s3_bucket_server_side_encryption_configuration" "config_bucket_encryption" {
  bucket = aws_s3_bucket.config_bucket.id
  
  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm     = "aws:kms"
      kms_master_key_id = aws_kms_key.s3_key.arn
    }
  }
}

# S3 bucket versioning
resource "aws_s3_bucket_versioning" "config_bucket_versioning" {
  bucket = aws_s3_bucket.config_bucket.id
  
  versioning_configuration {
    status = "Enabled"
  }
}

# S3 bucket for CloudTrail logs
resource "aws_s3_bucket" "cloudtrail_bucket" {
  bucket = "fedramp-data-mesh-cloudtrail-${data.aws_caller_identity.current.account_id}-${var.environment}"
  
  tags = {
    Name        = "fedramp-data-mesh-cloudtrail-bucket"
    Environment = var.environment
    Project     = "FedRAMP-Data-Mesh"
  }
}

# S3 bucket server-side encryption for CloudTrail
resource "aws_s3_bucket_server_side_encryption_configuration" "cloudtrail_bucket_encryption" {
  bucket = aws_s3_bucket.cloudtrail_bucket.id
  
  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm     = "aws:kms"
      kms_master_key_id = aws_kms_key.s3_key.arn
    }
  }
}

# S3 bucket versioning for CloudTrail
resource "aws_s3_bucket_versioning" "cloudtrail_bucket_versioning" {
  bucket = aws_s3_bucket.cloudtrail_bucket.id
  
  versioning_configuration {
    status = "Enabled"
  }
}

# S3 bucket policy for CloudTrail
resource "aws_s3_bucket_policy" "cloudtrail_bucket_policy" {
  bucket = aws_s3_bucket.cloudtrail_bucket.id
  policy = data.aws_iam_policy_document.cloudtrail_bucket_policy.json
}

data "aws_iam_policy_document" "cloudtrail_bucket_policy" {
  statement {
    sid    = "AWSCloudTrailAclCheck"
    effect = "Allow"
    
    principals {
      type        = "Service"
      identifiers = ["cloudtrail.amazonaws.com"]
    }
    
    actions   = ["s3:GetBucketAcl"]
    resources = [aws_s3_bucket.cloudtrail_bucket.arn]
  }
  
  statement {
    sid    = "AWSCloudTrailWrite"
    effect = "Allow"
    
    principals {
      type        = "Service"
      identifiers = ["cloudtrail.amazonaws.com"]
    }
    
    actions   = ["s3:PutObject"]
    resources = ["${aws_s3_bucket.cloudtrail_bucket.arn}/AWSLogs/${data.aws_caller_identity.current.account_id}/*"]
    
    condition {
      test     = "StringEquals"
      variable = "s3:x-amz-acl"
      values   = ["bucket-owner-full-control"]
    }
  }
}

# CloudTrail for audit logging (required for FedRAMP)
resource "aws_cloudtrail" "fedramp" {
  name                          = "fedramp-data-mesh-trail-${var.environment}"
  s3_bucket_name                = aws_s3_bucket.cloudtrail_bucket.id
  include_global_service_events = true
  is_multi_region_trail         = true
  enable_log_file_validation    = true
  kms_key_id                    = aws_kms_key.s3_key.arn
  
  # FedRAMP requires logging of data events
  event_selector {
    read_write_type           = "All"
    include_management_events = true
    
    data_resource {
      type   = "AWS::S3::Object"
      values = ["arn:aws:s3:::"]
    }
  }
  
  tags = {
    Name        = "fedramp-data-mesh-cloudtrail"
    Environment = var.environment
    Project     = "FedRAMP-Data-Mesh"
  }
}

# IAM Role for AWS Config
resource "aws_iam_role" "config_role" {
  name = "fedramp-data-mesh-config-role-${var.environment}"
  
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "config.amazonaws.com"
        }
      }
    ]
  })
}

# IAM Policy for AWS Config
resource "aws_iam_role_policy" "config_policy" {
  name = "fedramp-data-mesh-config-policy-${var.environment}"
  role = aws_iam_role.config_role.id
  
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "s3:PutObject",
          "s3:PutObjectAcl"
        ]
        Effect = "Allow"
        Resource = [
          "${aws_s3_bucket.config_bucket.arn}/*"
        ]
        Condition = {
          StringLike = {
            "s3:x-amz-acl" = "bucket-owner-full-control"
          }
        }
      },
      {
        Action = [
          "s3:GetBucketAcl"
        ]
        Effect = "Allow"
        Resource = [
          aws_s3_bucket.config_bucket.arn
        ]
      }
    ]
  })
}

# Security Hub for FedRAMP compliance monitoring
resource "aws_securityhub_account" "fedramp" {}

# Enable Security Hub FedRAMP standard
resource "aws_securityhub_standards_subscription" "fedramp" {
  standards_arn = "arn:aws:securityhub:${var.aws_region}::standards/nist-800-53/v/5.0.0"
  
  depends_on = [aws_securityhub_account.fedramp]
}

# GuardDuty for threat detection (required for FedRAMP)
resource "aws_guardduty_detector" "fedramp" {
  enable = true
  
  finding_publishing_frequency = "SIX_HOURS"
  
  tags = {
    Name        = "fedramp-data-mesh-guardduty"
    Environment = var.environment
    Project     = "FedRAMP-Data-Mesh"
  }
}

data "aws_caller_identity" "current" {}
