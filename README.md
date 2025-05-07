# FedRAMP High Event-Driven Data Mesh

A FedRAMP Medium ==> High compliance, event-driven data mesh foundational architecture for a mock Construction Management Software serving U.S. Federal and State Government contracts.

## Overview

This repository contains the implementation of a decentralized, scalable, and secure data infrastructure that supports both batch (Spark) and real-time streaming (Kafka/Flink) workloads. It is designed to meet rigorous FedRAMP High compliance requirements while enabling domain-driven data ownership and real-time data sharing.

## Core Features

- **Decentralized Data Ownership:** Empowering domain teams to own, build, and serve their data as products
- **Event-Driven Communication:** Event streams as the primary mechanism for sharing data between domains
- **Self-Service Data Platform:** Tools and infrastructure for domains to publish and consume data products
- **Federated Governance:** Global standards for schemas, quality, security, and interoperability
- **FedRAMP High Compliance:** Stringent security controls to protect sensitive government data
- **LakeHouse Foundation:** Databricks and Apache Iceberg for unified analytics on AWS S3
- **Developer Enablement:** Go-based CLI tool with DuckDB integration for local analytics

## Technology Stack

- **AWS Infrastructure:** VPC, S3, MSK, IAM, KMS, etc.
- **Event Streaming:** Apache Kafka (AWS MSK)
- **Schema Management:** Confluent Schema Registry
- **Change Data Capture:** Debezium
- **Data Processing:** Databricks (Spark)
- **Table Format:** Apache Iceberg
- **Developer CLI:** Go with Charm.sh and DuckDB

## Getting Started

1. **Prerequisites:**
   - AWS Account with appropriate permissions
   - Terraform installed
   - Go development environment (for CLI)
   - Docker and Kubernetes tools

2. **Setup Infrastructure:**
   ```bash
   cd platform/infrastructure/terraform
   terraform init
   terraform plan -var-file=environments/dev/terraform.tfvars
   terraform apply -var-file=environments/dev/terraform.tfvars
Deploy Core Services:

Copycd platform/infrastructure/kubernetes
./deploy.sh
Build and Test CLI:

Copycd cli
make build
./bin/dmesh --help
License
This project is licensed under the Apache License 2.0 - see the LICENSE file for details.

Contributing
See CONTRIBUTING.md for details on how to contribute to this project.


