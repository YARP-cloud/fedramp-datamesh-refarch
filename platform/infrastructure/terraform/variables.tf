# Variables for the platform-level Terraform configuration

variable "aws_region" {
  description = "AWS region to deploy resources"
  type        = string
  default     = "us-east-1"
}

variable "environment" {
  description = "Deployment environment (dev, test, prod)"
  type        = string
  default     = "dev"
}

variable "vpc_cidr_block" {
  description = "CIDR block for the VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "public_subnets" {
  description = "CIDR blocks for public subnets"
  type        = list(string)
  default     = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]
}

variable "private_subnets" {
  description = "CIDR blocks for private subnets"
  type        = list(string)
  default     = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
}

variable "database_subnets" {
  description = "CIDR blocks for database subnets"
  type        = list(string)
  default     = ["10.0.201.0/24", "10.0.202.0/24", "10.0.203.0/24"]
}

variable "enable_vpn" {
  description = "Enable VPN Gateway"
  type        = bool
  default     = false
}

variable "domains" {
  description = "List of business domains"
  type        = list(string)
  default     = ["project_management", "financials", "safety", "documents"]
}

variable "initial_data_products" {
  description = "Initial data products to create"
  type        = list(object({
    domain = string
    name   = string
  }))
  default     = [
    {
      domain = "project_management"
      name   = "project_state_events"
    },
    {
      domain = "financials"
      name   = "cost_item_events"
    }
  ]
}

variable "kafka_broker_count" {
  description = "Number of Kafka broker nodes"
  type        = number
  default     = 3
}

variable "kafka_broker_instance_type" {
  description = "Instance type for Kafka broker nodes"
  type        = string
  default     = "kafka.m5.large"
}

variable "kafka_broker_volume_size" {
  description = "Volume size for Kafka broker nodes (GB)"
  type        = number
  default     = 1000
}

variable "allowed_role_arns" {
  description = "List of IAM role ARNs allowed to access the data mesh resources"
  type        = list(string)
  default     = []
}

variable "databricks_account_id" {
  description = "Databricks account ID"
  type        = string
  sensitive   = true
}

variable "databricks_account_username" {
  description = "Databricks account username"
  type        = string
  sensitive   = true
}

variable "databricks_account_password" {
  description = "Databricks account password"
  type        = string
  sensitive   = true
}

variable "alarm_sns_topic_arn" {
  description = "SNS topic ARN for CloudWatch alarms"
  type        = string
  default     = ""
}
