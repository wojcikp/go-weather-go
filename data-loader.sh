#!/bin/bash

clickhouse-client --user="$CLICKHOUSE_USER" --password="$CLICKHOUSE_PASSWORD" --query="
CREATE TABLE IF NOT EXISTS weather_database.weather_scores (
    city String,
    time DateTime,
    temperature Float64,
    wind_speed Float64,
    weather_code Int64
) ENGINE = ReplacingMergeTree
PRIMARY KEY (city, time)
ORDER BY (city, time)
SETTINGS index_granularity = 8192;
"

clickhouse-client --user="$CLICKHOUSE_USER" --password="$CLICKHOUSE_PASSWORD" --query="INSERT INTO weather_database.weather_scores FORMAT TSV" < /docker-entrypoint-initdb.d/data.tsv

