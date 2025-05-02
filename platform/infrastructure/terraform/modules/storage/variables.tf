variable "aws_region" {
  description = "AWS region to deploy resources"
  type        = string
}

variable "environment" {
  description = "Deployment environment (dev, test, prod)"
  type        = string
}

variable "kms_key_arn" {
  description = "ARN of KMS key for S3 encryption"
  type        = string
}

variable "domains" {
  description = "List of business domains"
  type        = list(string)
}

variable "initial_data_products" {
  description = "Initial data products to create"
  type        = list(object({
    domain = string
    name   = string
  }))
}
