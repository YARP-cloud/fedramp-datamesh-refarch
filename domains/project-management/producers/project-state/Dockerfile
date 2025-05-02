FROM confluentinc/cp-kafka-connect:7.4.0

# Install Debezium connector
RUN confluent-hub install --no-prompt debezium/debezium-connector-postgresql:2.3.0

# Install additional transforms
RUN confluent-hub install --no-prompt confluentinc/connect-transforms:latest

# Copy custom configuration and startup scripts
COPY config/connect-distributed.properties /etc/kafka/connect-distributed.properties
COPY scripts/start-connect.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/start-connect.sh

# Security enhancements
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    ca-certificates \
    curl && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Default directory for certificates
RUN mkdir -p /etc/kafka/secrets

# Remove unnecessary permissions
RUN chmod -R 500 /usr/share/java/kafka/bin
RUN chmod -R 400 /usr/share/java/kafka/lib

# Drop capabilities
USER 1001

ENTRYPOINT ["/usr/local/bin/start-connect.sh"]
