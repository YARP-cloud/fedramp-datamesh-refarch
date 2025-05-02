variable "aws_region" {
  description = "AWS region to deploy resources"
  type        = string
}

variable "environment" {
  description = "Deployment environment (dev, test, prod)"
  type        = string
}

variable "vpc_id" {
  description = "VPC ID where resources will be deployed"
  type        = string
}

variable "logs_bucket_name" {
  description = "Name of the S3 bucket for access logs"
  type        = string
}

variable "allowed_role_arns" {
  description = "List of IAM role ARNs allowed to access the data mesh resources"
  type        = list(string)
  default     = []
}
