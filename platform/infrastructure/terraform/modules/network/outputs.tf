output "vpc_id" {
  description = "The ID of the VPC"
  value       = module.vpc.vpc_id
}

output "vpc_cidr_block" {
  description = "The CIDR block of the VPC"
  value       = module.vpc.vpc_cidr_block
}

output "public_subnet_ids" {
  description = "List of public subnet IDs"
  value       = module.vpc.public_subnets
}

output "private_subnet_ids" {
  description = "List of private subnet IDs"
  value       = module.vpc.private_subnets
}

output "database_subnet_ids" {
  description = "List of database subnet IDs"
  value       = module.vpc.database_subnets
}

output "nat_gateway_ids" {
  description = "List of NAT Gateway IDs"
  value       = module.vpc.natgw_ids
}

output "kafka_security_group_id" {
  description = "ID of the security group for Kafka"
  value       = aws_security_group.kafka_security_group.id
}

output "databricks_security_group_id" {
  description = "ID of the security group for Databricks"
  value       = aws_security_group.databricks_security_group.id
}

output "vpc_endpoints_security_group_id" {
  description = "ID of the security group for VPC endpoints"
  value       = aws_security_group.vpc_endpoints_sg.id
}

output "vpc_flow_logs_bucket" {
  description = "Name of the S3 bucket for VPC flow logs"
  value       = aws_s3_bucket.vpc_flow_logs.bucket
}
