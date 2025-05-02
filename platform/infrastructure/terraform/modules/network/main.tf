# Network infrastructure module for FedRAMP High Data Mesh

provider "aws" {
  region = var.aws_region
}

module "vpc" {
  source = "terraform-aws-modules/vpc/aws"
  version = "3.0.0"

  name = "fedramp-data-mesh-vpc-${var.environment}"
  cidr = var.vpc_cidr_block

  azs             = ["${var.aws_region}a", "${var.aws_region}b", "${var.aws_region}c"]
  private_subnets = var.private_subnets
  public_subnets  = var.public_subnets
  database_subnets = var.database_subnets

  enable_nat_gateway = true
  single_nat_gateway = var.environment != "prod"
  one_nat_gateway_per_az = var.environment == "prod"

  enable_vpn_gateway = var.enable_vpn

  # FedRAMP-specific VPC Flow Logs
  enable_flow_log = true
  flow_log_destination_type = "s3"
  flow_log_destination_arn = aws_s3_bucket.vpc_flow_logs.arn
  flow_log_traffic_type = "ALL"
  
  # Additional security configurations
  manage_default_security_group = true
  default_security_group_ingress = []
  default_security_group_egress = []
  
  # DNS settings
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Environment = var.environment
    Project     = "FedRAMP-Data-Mesh"
  }
}

resource "aws_s3_bucket" "vpc_flow_logs" {
  bucket = "fedramp-data-mesh-vpc-flow-logs-${var.environment}"
  
  tags = {
    Name        = "fedramp-data-mesh-vpc-flow-logs"
    Environment = var.environment
    Project     = "FedRAMP-Data-Mesh"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "vpc_logs_encryption" {
  bucket = aws_s3_bucket.vpc_flow_logs.id
  
  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_lifecycle_configuration" "vpc_logs_lifecycle" {
  bucket = aws_s3_bucket.vpc_flow_logs.id
  
  rule {
    id = "log-expiration"
    status = "Enabled"
    
    expiration {
      days = 731  # FedRAMP requires 2 years of logs
    }
  }
}

# Security group for Kafka access
resource "aws_security_group" "kafka_security_group" {
  name        = "fedramp-data-mesh-kafka-sg-${var.environment}"
  description = "Security group for Kafka access"
  vpc_id      = module.vpc.vpc_id

  # Internal traffic only - no public access to Kafka
  ingress {
    from_port   = 9094
    to_port     = 9094
    protocol    = "tcp"
    cidr_blocks = var.private_subnets
    description = "Kafka TLS"
  }
  
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
    description = "Allow all outbound traffic"
  }

  tags = {
    Name        = "fedramp-data-mesh-kafka-sg"
    Environment = var.environment
    Project     = "FedRAMP-Data-Mesh"
  }
}

# Security group for Databricks
resource "aws_security_group" "databricks_security_group" {
  name        = "fedramp-data-mesh-databricks-sg-${var.environment}"
  description = "Security group for Databricks workspace"
  vpc_id      = module.vpc.vpc_id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
    description = "Allow all outbound traffic"
  }

  tags = {
    Name        = "fedramp-data-mesh-databricks-sg"
    Environment = var.environment
    Project     = "FedRAMP-Data-Mesh"
  }
}

# VPC Endpoints for AWS services to enhance security
resource "aws_vpc_endpoint" "s3" {
  vpc_id            = module.vpc.vpc_id
  service_name      = "com.amazonaws.${var.aws_region}.s3"
  vpc_endpoint_type = "Gateway"
  route_table_ids   = concat(module.vpc.private_route_table_ids, module.vpc.public_route_table_ids)
  
  tags = {
    Name        = "fedramp-data-mesh-s3-endpoint"
    Environment = var.environment
    Project     = "FedRAMP-Data-Mesh"
  }
}

resource "aws_vpc_endpoint" "dynamodb" {
  vpc_id            = module.vpc.vpc_id
  service_name      = "com.amazonaws.${var.aws_region}.dynamodb"
  vpc_endpoint_type = "Gateway"
  route_table_ids   = concat(module.vpc.private_route_table_ids, module.vpc.public_route_table_ids)
  
  tags = {
    Name        = "fedramp-data-mesh-dynamodb-endpoint"
    Environment = var.environment
    Project     = "FedRAMP-Data-Mesh"
  }
}

# Create Interface VPC Endpoints for other AWS services
locals {
  interface_endpoints = [
    "ecr.api",
    "ecr.dkr",
    "kms",
    "logs",
    "monitoring",
    "sqs",
    "sns",
    "glue"
  ]
}

resource "aws_security_group" "vpc_endpoints_sg" {
  name        = "fedramp-data-mesh-vpc-endpoints-sg-${var.environment}"
  description = "Security group for VPC endpoints"
  vpc_id      = module.vpc.vpc_id

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = [var.vpc_cidr_block]
    description = "HTTPS from VPC"
  }
  
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
    description = "Allow all outbound traffic"
  }

  tags = {
    Name        = "fedramp-data-mesh-vpc-endpoints-sg"
    Environment = var.environment
    Project     = "FedRAMP-Data-Mesh"
  }
}

resource "aws_vpc_endpoint" "interface_endpoints" {
  for_each = toset(local.interface_endpoints)
  
  vpc_id              = module.vpc.vpc_id
  service_name        = "com.amazonaws.${var.aws_region}.${each.value}"
  vpc_endpoint_type   = "Interface"
  subnet_ids          = module.vpc.private_subnets
  security_group_ids  = [aws_security_group.vpc_endpoints_sg.id]
  private_dns_enabled = true
  
  tags = {
    Name        = "fedramp-data-mesh-${each.value}-endpoint"
    Environment = var.environment
    Project     = "FedRAMP-Data-Mesh"
  }
}
