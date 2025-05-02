variable "aws_region" {
  description = "AWS region to deploy resources"
  type        = string
}

variable "environment" {
  description = "Deployment environment (dev, test, prod)"
  type        = string
}

variable "vpc_cidr_block" {
  description = "CIDR block for the VPC"
  type        = string
}

variable "public_subnets" {
  description = "CIDR blocks for public subnets"
  type        = list(string)
}

variable "private_subnets" {
  description = "CIDR blocks for private subnets"
  type        = list(string)
}

variable "database_subnets" {
  description = "CIDR blocks for database subnets"
  type        = list(string)
}

variable "enable_vpn" {
  description = "Enable VPN Gateway"
  type        = bool
  default     = false
}

variable "logs_bucket_name" {
  description = "Name of the S3 bucket for access logs"
  type        = string
}
