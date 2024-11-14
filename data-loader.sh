#!/bin/bash
# Start ClickHouse server in the background
clickhouse-server --daemon

# Wait for the server to start
sleep 5

# Load the schema from schema.sql
clickhouse-client --multiquery < /docker-entrypoint-initdb.d/schema.sql

# Load the data from data.tsv
# Replace `your_table_name` with the actual table name in ClickHouse
clickhouse-client --query="INSERT INTO weather_scores FORMAT TSV" < /docker-entrypoint-initdb.d/data.tsv

# Keep the container running
tail -f /dev/null
