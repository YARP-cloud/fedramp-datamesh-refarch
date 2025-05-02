# Development environment configuration

environment     = "dev"
aws_region      = "us-east-1"
vpc_cidr_block  = "10.0.0.0/16"
public_subnets  = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]
private_subnets = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
database_subnets = ["10.0.201.0/24", "10.0.202.0/24", "10.0.203.0/24"]
enable_vpn      = false

domains = ["project_management", "financials", "safety", "documents"]

initial_data_products = [
  {
    domain = "project_management"
    name   = "project_state_events"
  },
  {
    domain = "project_management"
    name   = "project_activity_events"
  },
  {
    domain = "financials"
    name   = "cost_item_events"
  },
  {
    domain = "safety"
    name   = "incident_events"
  }
]

kafka_broker_count       = 3
kafka_broker_instance_type = "kafka.m5.large"
kafka_broker_volume_size = 1000

# Replace these with actual values in a real implementation
# Do not store sensitive information in version control
databricks_account_id       = "REPLACE_WITH_ACTUAL_ID"
databricks_account_username = "REPLACE_WITH_ACTUAL_USERNAME"
databricks_account_password = "REPLACE_WITH_ACTUAL_PASSWORD"

allowed_role_arns = [
  "arn:aws:iam::123456789012:role/DataMeshAdmin",
  "arn:aws:iam::123456789012:role/DataMeshDeveloper"
]

alarm_sns_topic_arn = "arn:aws:sns:us-east-1:123456789012:DataMeshAlarms"
