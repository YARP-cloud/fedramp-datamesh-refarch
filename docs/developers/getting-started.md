# Developer Getting Started Guide

This guide helps developers get started with the FedRAMP High Event-Driven Data Mesh, including setting up the development environment, accessing data products, and creating new data products.

## Prerequisites

Before you begin, ensure you have the following:

1. **Access Credentials**:
   - AWS IAM credentials with appropriate permissions
   - Databricks workspace access
   - Kafka access credentials

2. **Required Tools**:
   - Go 1.18+ (for CLI tool)
   - AWS CLI configured with your credentials
   - Git client
   - Docker and Docker Compose (for local development)
   - Terraform (for infrastructure changes)

## Installing the CLI Tool

The Data Mesh CLI (`dmesh`) provides a convenient way to interact with the data mesh from your local machine.

### From Binary Release

1. Download the latest release for your platform from the releases page
2. Extract the archive and move the binary to a location in your PATH:
```bash
tar -xzf dmesh_v1.0.0_linux_amd64.tar.gz
chmod +x dmesh
sudo mv dmesh /usr/local/bin/
3. Verify the installation:
dmesh --version

### From Source

1. Clone the repository:
```bash
git clone https://github.com/frocore/fedramp-data-mesh.git
cd fedramp-data-mesh
2. Build the CLI tool:
```bash
make build
3. Install the binary:
```bash
sudo cp cli/bin/dmesh /usr/local/bin/

### Configuring the CLI

1. Run the CLI once to create the default configuration:
```bash
dmesh info
2. Edit the configuration file at ~/.fedramp-data-mesh/config.yaml:
```bash
aws_region: us-east-1
aws_profile: fedramp-data-mesh
aws_account_id: "123456789012"
default_role: DataMeshDeveloper
catalog_url: "https://catalog.fedramp-data-mesh.example.com"
s3_data_lake: "s3://fedramp-data-mesh-lake-123456789012-dev"
schema_registry_url: "https://schema-registry.fedramp-data-mesh.example.com"

### Discovering Data Products

### Using the CLI

1. List all available data products:
```bash
dmesh discover
2. Filter by domain:
```bash
dmesh discover --domain project_management
3. View details for a specific data product:
```bash
dmesh info project_management.project_state_events
4. View the schema for a data product:
```bash
dmesh schema project_management.project_state_events

### Using the Databricks Catalog

Log in to the Databrick

1. List all available data products:
```bash
dmesh discover
2. Filter by domain:
```bash
dmesh discover --domain project_management
3. View details for a specific data product:
```bash
dmesh info project_management.project_state_events
4. View the schema for a data product:
```bash
dmesh schema project_management.project_state_events

### Using the Databricks Catalog

1. Log in to the Databricks workspace
2. Navigate to the Data tab
3. Browse the Catalog for available data products
4. View table details, schema, and sample data

### Querying Data Products

### Using the CLI

1. Run a simple query against a data product:
```bash
dmesh query "SELECT * FROM project_management.project_state_latest LIMIT 10"
2. Use the interactive query UI:
```bash
dmesh query --interactive
3. Output results in different formats:
```bash
dmesh query "SELECT * FROM project_management.project_state_latest LIMIT 10" --output csv
dmesh query "SELECT * FROM project_management.project_state_latest LIMIT 10" --output json

### Using Databricks

1. Log in to the Databricks workspace
2. Create a new notebook
3. Use Spark SQL to query data products:
```bash
SELECT * FROM project_management.project_state_latest LIMIT 10
4. Use Spark DataFrame API:
```python
df = spark.table("project_management.project_state_latest")
display(df.limit(10))

### Creating a New Data Product

### 1. Define the Schema
Create an Avro schema file for your data product. Example:
```bash
{
  "type": "record",
  "name": "TaskStateEvent",
  "namespace": "com.frocore.projectmanagement.events",
  "doc": "Represents the current state of a task after a change",
  "fields": [
    {
      "name": "event_id",
      "type": "string",
      "doc": "Unique identifier for this event"
    },
    {
      "name": "event_timestamp",
      "type": {
        "type": "long",
        "logicalType": "timestamp-millis"
      },
      "doc": "Timestamp when this event was created"
    },
    {
      "name": "event_type",
      "type": {
        "type": "enum",
        "name": "TaskEventType",
        "symbols": ["CREATED", "UPDATED", "DELETED"]
      },
      "doc": "Type of event that occurred"
    },
    {
      "name": "task_id",
      "type": "string",
      "doc": "Unique identifier for the task"
    },
    {
      "name": "project_id",
      "type": "string",
      "doc": "ID of the project this task belongs to"
    },
    {
      "name": "title",
      "type": "string",
      "doc": "Task title"
    },
    {
      "name": "description",
      "type": ["string", "null"],
      "doc": "Task description"
    },
    {
      "name": "status",
      "type": {
        "type": "enum",
        "name": "TaskStatus",
        "symbols": ["TODO", "IN_PROGRESS", "DONE", "BLOCKED"]
      },
      "doc": "Current status of the task"
    },
    {
      "name": "assignee_id",
      "type": ["string", "null"],
      "doc": "ID of the person assigned to the task"
    },
    {
      "name": "created_at",
      "type": {
        "type": "long",
        "logicalType": "timestamp-millis"
      },
      "doc": "Timestamp when the task was initially created"
    },
    {
      "name": "modified_at",
      "type": {
        "type": "long",
        "logicalType": "timestamp-millis"
      },
      "doc": "Timestamp when the task was last modified"
    }
  ]
}

### 2. Register the Schema

Register the schema with the Schema Registry:
```bash
curl -X POST -H "Content-Type: application/vnd.schemaregistry.v1+json" \
  --data @task_state_event.avsc \
  https://schema-registry.fedramp-data-mesh.example.com/subjects/projects.task_state_events/versions

### 3. Create a Data Product Definition

Create a YAML definition for your data product:
```yaml
kind: DataProduct
apiVersion: datamesh.frocore.io/v1
metadata:
  name: task_state_events
  domain: project_management
  owner: project-management-team@frocore.io
  description: State events for task entities
  documentation: |
    This data product captures the state of each task after changes.
    It is related to the project_state_events data product.
spec:
  schemaRef:
    type: avro
    path: /domains/project-management/schemas/task_state_event.avsc
  eventStream:
    topicName: projects.task_state_events
    partitionKey: task_id
    retention: 
      time: 30d
    replication: 3
  tables:
    - name: task_state_history
      catalog: project_management
      format: iceberg
      location: s3://fedramp-data-mesh-lake/project_management/task_state_history
      partitioning:
        - name: event_date
          transform: "day(event_timestamp)"
    - name: task_state_latest
      catalog: project_management
      format: iceberg
      location: s3://fedramp-data-mesh-lake/project_management/task_state_latest
      retention:
        snapshots: 5
  sla:
    latency: 1m
    availability: 99.9%
  securityClassification: CONTROLLED_UNCLASSIFIED
  lineage:
    upstream:
      - source: projects-db.public.tasks
        type: database-table
  access:
    roles:
      - name: project_admin
        permissions: [read, write]
      - name: project_analyst
        permissions: [read]
      - name: data_engineer
        permissions: [read]

### 4. Configure the Source Connector

Create a configuration for your CDC connector:
```bash
{
  "name": "tasks-source-connector",
  "config": {
    "connector.class": "io.debezium.connector.postgresql.PostgresConnector",
    "database.hostname": "${DB_HOST}",
    "database.port": "${DB_PORT}",
    "database.user": "${DB_USER}",
    "database.password": "${DB_PASSWORD}",
    "database.dbname": "frocore",
    "database.server.name": "frocore-projects",
    "table.include.list": "public.tasks",
    "schema.include.list": "public",
    "database.history.kafka.bootstrap.servers": "${KAFKA_BOOTSTRAP_SERVERS}",
    "database.history.kafka.topic": "schema-changes.frocore.tasks",
    "transforms": "unwrap,AddSourceMetadata",
    "transforms.unwrap.type": "io.debezium.transforms.ExtractNewRecordState",
    "transforms.unwrap.drop.tombstones": "false",
    "transforms.unwrap.delete.handling.mode": "rewrite",
    "transforms.AddSourceMetadata.type": "org.apache.kafka.connect.transforms.InsertField$Value",
    "transforms.AddSourceMetadata.static.field": "source_system",
    "transforms.AddSourceMetadata.static.value": "projects-db",
    "key.converter": "io.confluent.connect.avro.AvroConverter",
    "key.converter.schema.registry.url": "${SCHEMA_REGISTRY_URL}",
    "value.converter": "io.confluent.connect.avro.AvroConverter",
    "value.converter.schema.registry.url": "${SCHEMA_REGISTRY_URL}"
  }
}

### 5. Implement the Data Processor

Create a Spark job to process the events:
```java
package com.frocore.datamesh.project_management.processors

import org.apache.spark.sql.{DataFrame, SparkSession}
import org.apache.spark.sql.functions._
import org.apache.spark.sql.streaming.Trigger
import org.apache.spark.sql.avro.functions._
import java.util.UUID

object TaskStateProcessor {
  def main(args: Array[String]): Unit = {
    val spark = SparkSession.builder()
      .appName("Task State Processor")
      .config("spark.sql.extensions", "org.apache.iceberg.spark.extensions.IcebergSparkSessionExtensions")
      .config("spark.sql.catalog.ice", "org.apache.iceberg.spark.SparkCatalog")
      .config("spark.sql.catalog.ice.type", "hadoop")
      .config("spark.sql.catalog.ice.warehouse", "s3a://fedramp-data-mesh-warehouse/")
      .getOrCreate()

    import spark.implicits._

    // Read from Kafka topic
    val kafkaStreamDF = spark
      .readStream
      .format("kafka")
      .option("kafka.bootstrap.servers", "${KAFKA_BOOTSTRAP_SERVERS}")
      .option("subscribe", "projects.task_state_events")
      .option("startingOffsets", "latest")
      .option("kafka.security.protocol", "SSL")
      .load()

    // Parse the Avro payload
    val parsedDF = kafkaStreamDF
      .select(
        col("key").cast("string").as("kafka_key"),
        col("timestamp").as("kafka_timestamp"),
        from_avro(col("value"), "task_state_event_schema").as("event")
      )
      .select(
        col("kafka_key"),
        col("kafka_timestamp"),
        col("event.*")
      )

    // Add processing metadata
    val enrichedDF = parsedDF
      .withColumn("processing_time", current_timestamp())
      .withColumn("processing_id", lit(UUID.randomUUID().toString))

    // Write to Iceberg table (history)
    val historyQuery = enrichedDF
      .writeStream
      .format("iceberg")
      .outputMode("append")
      .option("path", "ice.project_management.task_state_history")
      .option("checkpointLocation", "s3a://fedramp-data-mesh-checkpoints/task_state_processor/history/")
      .trigger(Trigger.ProcessingTime("1 minute"))
      .start()

    // Write to Iceberg table (latest state)
    val latestQuery = enrichedDF
      .writeStream
      .foreachBatch { (batchDF: DataFrame, batchId: Long) =>
        // Upsert (merge) the latest state
        batchDF.createOrReplaceTempView("updates")
        
        spark.sql("""
          MERGE INTO ice.project_management.task_state_latest target
          USING (
            SELECT task_id, 
                   project_id,
                   title, 
                   description, 
                   status, 
                   assignee_id,
                   created_at,
                   modified_at,
                   event_timestamp
            FROM updates
          ) source
          ON target.task_id = source.task_id
          WHEN MATCHED THEN
            UPDATE SET *
          WHEN NOT MATCHED THEN
            INSERT *
        """)
      }
      .option("checkpointLocation", "s3a://fedramp-data-mesh-checkpoints/task_state_processor/latest/")
      .trigger(Trigger.ProcessingTime("1 minute"))
      .start()

    spark.streams.awaitAnyTermination()
  }
}

### 6. Deploy the Data Product

1. Submit your Kafka connector configuration:
```bash
curl -X POST -H "Content-Type: application/json" \
  --data @tasks-source-connector.json \
  http://kafka-connect.fedramp-data-mesh.example.com:8083/connectors
2. Deploy the Spark job to Databricks:
```bash
databricks jobs create --json @job-config.json
3. Register the data product tables in Unity Catalog:
```bash
CREATE TABLE project_management.task_state_history
USING iceberg
LOCATION 's3://fedramp-data-mesh-lake/project_management/task_state_history';

CREATE TABLE project_management.task_state_latest
USING iceberg
LOCATION 's3://fedramp-data-mesh-lake/project_management/task_state_latest';

### Next Steps

- Learn more about Data Product Design Patterns
- Explore Advanced Querying Techniques
- Understand Event Schema Evolution
- Set up Data Quality Monitoring
