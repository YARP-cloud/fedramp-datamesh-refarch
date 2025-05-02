package com.frocore.datamesh.project_management.processors

import org.apache.spark.sql.{DataFrame, SparkSession}
import org.apache.spark.sql.functions._
import org.apache.spark.sql.streaming.Trigger
import org.apache.spark.sql.avro.functions._
import org.apache.spark.sql.types._
import java.util.UUID

object ProjectStateProcessor {
  def main(args: Array[String]): Unit = {
    val spark = SparkSession.builder()
      .appName("Project State Processor")
      .config("spark.sql.extensions", "org.apache.iceberg.spark.extensions.IcebergSparkSessionExtensions")
      .config("spark.sql.catalog.ice", "org.apache.iceberg.spark.SparkCatalog")
      .config("spark.sql.catalog.ice.type", "hadoop")
      .config("spark.sql.catalog.ice.warehouse", "s3a://fedramp-data-mesh-warehouse/")
      .config("spark.sql.catalog.ice.io-impl", "org.apache.iceberg.aws.s3.S3FileIO")
      .config("spark.databricks.delta.schema.autoMerge.enabled", "true")
      .getOrCreate()

    import spark.implicits._

    // Read from Kafka topic
    val kafkaStreamDF = spark
      .readStream
      .format("kafka")
      .option("kafka.bootstrap.servers", "${KAFKA_BOOTSTRAP_SERVERS}")
      .option("subscribe", "projects.project_state_events")
      .option("startingOffsets", "latest")
      .option("kafka.security.protocol", "SSL")
      .option("kafka.ssl.truststore.location", "/mnt/certs/kafka.truststore.jks")
      .option("kafka.ssl.truststore.password", "${TRUSTSTORE_PASSWORD}")
      .option("kafka.ssl.keystore.location", "/mnt/certs/kafka.keystore.jks")
      .option("kafka.ssl.keystore.password", "${KEYSTORE_PASSWORD}")
      .option("kafka.ssl.key.password", "${KEY_PASSWORD}")
      .load()

    // Schema registry URL
    val schemaRegistryUrl = "${SCHEMA_REGISTRY_URL}"
    
    // Get Avro schema from Schema Registry
    val projectStateEventSchema = spark.read
      .format("avro")
      .load("file:///tmp/schema/project_state_event.avsc")
      .schema

    // Parse the Avro payload
    val parsedDF = kafkaStreamDF
      .select(
        col("key").cast("string").as("kafka_key"),
        col("topic").as("kafka_topic"),
        col("partition").as("kafka_partition"),
        col("offset").as("kafka_offset"),
        col("timestamp").as("kafka_timestamp"),
        from_avro(col("value"), projectStateEventSchema).as("event")
      )
      .select(
        col("kafka_key"),
        col("kafka_topic"),
        col("kafka_partition"),
        col("kafka_offset"),
        col("kafka_timestamp"),
        col("event.*")
      )

    // Add processing metadata
    val enrichedDF = parsedDF
      .withColumn("processing_time", current_timestamp())
      .withColumn("processing_id", lit(UUID.randomUUID().toString))

    // Write to Iceberg table
    val query = enrichedDF
      .writeStream
      .format("iceberg")
      .outputMode("append")
      .option("path", "ice.project_management.project_state_history")
      .option("checkpointLocation", "s3a://fedramp-data-mesh-checkpoints/project_state_processor/")
      .trigger(Trigger.ProcessingTime("1 minute"))
      .start()

    // Also write the latest state to a separate table (using foreachBatch)
    val latestStateQuery = enrichedDF
      .writeStream
      .foreachBatch { (batchDF: DataFrame, batchId: Long) =>
        // Upsert (merge) the latest state
        batchDF.createOrReplaceTempView("updates")
        
        spark.sql("""
          MERGE INTO ice.project_management.project_state_latest target
          USING (
            SELECT project_id, 
                   name, 
                   description, 
                   status, 
                   start_date, 
                   end_date, 
                   owner_id,
                   budget,
                   location,
                   tags,
                   created_at,
                   modified_at,
                   event_timestamp,
                   event_type,
                   source_system,
                   security_classification,
                   metadata
            FROM updates
          ) source
          ON target.project_id = source.project_id
          WHEN MATCHED THEN
            UPDATE SET *
          WHEN NOT MATCHED THEN
            INSERT *
        """)
      }
      .option("checkpointLocation", "s3a://fedramp-data-mesh-checkpoints/project_state_latest/")
      .trigger(Trigger.ProcessingTime("1 minute"))
      .start()

    query.awaitTermination()
    latestStateQuery.awaitTermination()
  }
}
