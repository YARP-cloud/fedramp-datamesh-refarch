---
title: Deployment Guide
description: How to deploy and manage the FedRAMP High Event-Driven Data Mesh infrastructure.
---

This guide provides instructions for deploying and managing the FedRAMP High Event-Driven Data Mesh infrastructure.

## Prerequisites

- AWS Account with appropriate permissions
- Terraform 1.0+
- AWS CLI configured with appropriate credentials
- kubectl configured for Kubernetes access (if using EKS)
- Databricks CLI configured with workspace credentials

## Infrastructure Deployment

### 1. Initialize Terraform

1. Clone the repository:
   ```bash
   git clone https://github.com/frocore/fedramp-data-mesh.git
   cd fedramp-data-mesh
2. Initialize Terraform:
```bash
cd platform/infrastructure/terraform
terraform init -backend-config=environments/dev/backend.tfvars

### 2. Configure Env Variables
Create a .env file with the necessary environment variables:
```bash
# AWS Configuration
export AWS_REGION=us-east-1
export AWS_PROFILE=fedramp-data-mesh

# Databricks Configuration
export DATABRICKS_ACCOUNT_ID=your-account-id
export DATABRICKS_ACCOUNT_USERNAME=your-username
export DATABRICKS_ACCOUNT_PASSWORD=your-password
```
Source the envars:
```bash
source .env
```

### 3. Deploy infrastructure

1. Plan the deployment:
```bash
terraform plan -var-file=environments/dev/terraform.tfvars
```
2. Apply the changes:
```bash
terraform apply -var-file=environments/dev/terraform.tfvars
```
3. Take note of the outputs, which include important information about the deployed resources.

### 4. Deploy Kubernetes Components

1. Configure kubectl to connect to the newly created EKS cluster:
```bash
aws eks update-kubeconfig --name fedramp-data-mesh-eks-dev --region us-east-1
```
2. Deploy Kubernetes components:
```bash
cd ../kubernetes
kubectl apply -f namespace.yaml
kubectl apply -k schema-registry
kubectl apply -k kafka-connect
kubectl apply -k monitoring
```

### Data Product Deployment

### 1. Create Domain Catalogs

1. Log in to Databricks:
```bash
databricks configure --token
```
2. Create catalogs for each domain:
```bash
# Create Project Management catalog
databricks unity-catalog catalogs create \
  --name project_management \
  --comment "Project Management domain catalog"

# Create Financials catalog
databricks unity-catalog catalogs create \
  --name financials \
  --comment "Financials domain catalog"
```

### 2. Deploy Kafka Connectors

1. Configure Kafka Connect for source databases:
```bash
# Create Projects source connector
curl -X POST -H "Content-Type: application/json" \
  --data @domains/project-management/producers/project-state/connector-config.json \
  http://kafka-connect.fedramp-data-mesh.example.com:8083/connectors
```

### 3. Deploy Spark Jobs

1. Create a Databricks job for each processor:
```bash
# Create Project State Processor job
databricks jobs create --json @domains/project-management/processors/spark/job-config.json
```

### Monitoring and Operations

### 1. Set Up Monitoring

1. Configure CloudWatch Dashboards:
```bash
aws cloudwatch create-dashboard \
  --dashboard-name FedRAMP-DataMesh-Overview \
  --dashboard-body file://monitoring/cloudwatch-dashboards/overview.json```
2. Set up alerts:
```bash
aws cloudwatch put-metric-alarm \
  --alarm-name DataMesh-Kafka-HighLag \
  --alarm-description "Alert when Kafka consumer lag is too high" \
  --metric-name "kafka-consumer-lag" \
  --namespace "AWS/MSK" \
  --statistic Average \
  --period 300 \
  --threshold 1000 \
  --comparison-operator GreaterThanThreshold \
  --dimensions "Name=ClusterName,Value=fedramp-data-mesh-kafka-dev" \
  --evaluation-periods 2 \
  --alarm-actions ${SNS_TOPIC_ARN}
```

### 2. Regular Maintenance

1. Rotate encryption keys:
```bash
# Update KMS key for S3
aws kms enable-key-rotation --key-id ${S3_KMS_KEY_ID}
```
2. Update Kafka configurations:
```bash
aws kafka update-cluster-configuration \
  --cluster-arn ${KAFKA_CLUSTER_ARN} \
  --current-version ${CURRENT_CLUSTER_VERSION} \
  --configuration-info file://kafka-config-updates.json
```
3. Patch Kubernetes components:
```bash
kubectl apply -k schema-registry
```

### Backup and Disaster Recovery

### 1. Backup Strategy

1. S3 data is automatically versioned and cross-region replicated
2. Kafka topics should be configured with appropriate replication factor (3)
3. Critical configurations are stored in version control
4. Database backups are automated through AWS Backup

### 2. Disaster Recovery
1. In case of region failure, follow these steps:

- Activate standby infrastructure in secondary region
- Update DNS to point to secondary region
- Ensure all credentials and configurations are available

2. Test DR procedures regularly:
```bash
# Run DR test script
./scripts/dr-test.sh
```

### Security Operations

### 1. Access Management

1. Rotate credentials regularly:
```bash
# Rotate service account credentials
./scripts/rotate-credentials.sh
```
2. Review access:
```bash
# Generate access report
./scripts/access-review.sh > access-review-$(date +%Y-%m-%d).txt
```

### 2. Security Monitoring

1. Check GuardDuty findings:
```bash
aws guardduty list-findings
```
2. Run security scans:
```bash
# Run infrastructure security scan
./scripts/security-scan.sh
```

### Troubleshooting

### Common Issues

### 1. Kafka Connection Issues

- Check security groups
- Verify credentials
- Check network connectivity

### 2. Databricks Job Failures

- Check job logs
- Verify access to S3
- Check schema compatibility issues

### 3. Schema Evolution Errors

- Verify schema compatibility
- Check for breaking changes

For more detailed troubleshooting, refer to the Troubleshooting Guide.
