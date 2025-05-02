# Platform-level Terraform configuration

provider "aws" {
  region = var.aws_region

  # Use AWS GovCloud for FedRAMP High requirements
  # When using GovCloud, uncomment these lines and modify as needed
  # alias      = "govcloud"
  # region     = "us-gov-west-1"
  # profile    = "govcloud"
  
  default_tags {
    tags = {
      Project     = "FedRAMP-Data-Mesh"
      Environment = var.environment
      Terraform   = "true"
    }
  }
}

# Remote backend for state management
terraform {
  backend "s3" {
    # Configuration provided in backend.tfvars file
  }
  
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
    databricks = {
      source  = "databricks/databricks"
      version = "~> 1.0"
    }
  }
}

# Create core infrastructure
module "networking" {
  source = "./modules/network"
  
  environment       = var.environment
  aws_region        = var.aws_region
  vpc_cidr_block    = var.vpc_cidr_block
  public_subnets    = var.public_subnets
  private_subnets   = var.private_subnets
  database_subnets  = var.database_subnets
  enable_vpn        = var.enable_vpn
  logs_bucket_name  = module.storage.logs_bucket_name
}

module "security" {
  source = "./modules/security"
  
  environment        = var.environment
  aws_region         = var.aws_region
  vpc_id             = module.networking.vpc_id
  logs_bucket_name   = module.storage.logs_bucket_name
  allowed_role_arns  = var.allowed_role_arns
}

module "storage" {
  source = "./modules/storage"
  
  environment       = var.environment
  aws_region        = var.aws_region
  kms_key_arn       = module.security.s3_kms_key_arn
  domains           = var.domains
  initial_data_products = var.initial_data_products
}

module "kafka" {
  source = "./modules/kafka"
  
  environment           = var.environment
  aws_region            = var.aws_region
  vpc_id                = module.networking.vpc_id
  vpc_cidr_block        = var.vpc_cidr_block
  subnet_ids            = module.networking.private_subnet_ids
  broker_count          = var.kafka_broker_count
  broker_instance_type  = var.kafka_broker_instance_type
  broker_volume_size    = var.kafka_broker_volume_size
  kafka_kms_key_arn     = module.security.kafka_kms_key_arn
  logs_bucket_name      = module.storage.logs_bucket_name
}

module "databricks" {
  source = "./modules/databricks"
  
  environment                   = var.environment
  aws_region                    = var.aws_region
  vpc_id                        = module.networking.vpc_id
  subnet_ids                    = module.networking.private_subnet_ids
  security_group_id             = module.networking.databricks_security_group_id
  databricks_account_id         = var.databricks_account_id
  databricks_account_username   = var.databricks_account_username
  databricks_account_password   = var.databricks_account_password
  databricks_root_bucket        = module.storage.databricks_root_bucket
  databricks_cross_account_role_arn = module.security.databricks_role_arn
  databricks_kms_key_arn        = module.security.databricks_kms_key_arn
  databricks_kms_key_alias      = module.security.databricks_kms_key_alias
  unity_catalog_bucket          = module.storage.unity_catalog_bucket
  unity_catalog_role_arn        = module.security.unity_catalog_role_arn
  domains                       = var.domains
}

module "monitoring" {
  source = "./modules/monitoring"
  
  environment     = var.environment
  aws_region      = var.aws_region
  vpc_id          = module.networking.vpc_id
  logs_bucket_name = module.storage.logs_bucket_name
  alarm_sns_topic_arn = var.alarm_sns_topic_arn
}

# Outputs
output "vpc_id" {
  description = "The ID of the VPC"
  value       = module.networking.vpc_id
}

output "public_subnet_ids" {
  description = "List of public subnet IDs"
  value       = module.networking.public_subnet_ids
}

output "private_subnet_ids" {
  description = "List of private subnet IDs"
  value       = module.networking.private_subnet_ids
}

output "database_subnet_ids" {
  description = "List of database subnet IDs"
  value       = module.networking.database_subnet_ids
}

output "kafka_bootstrap_brokers" {
  description = "Kafka bootstrap brokers (TLS)"
  value       = module.kafka.bootstrap_brokers_tls
}
