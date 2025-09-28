#!/bin/sh
set -e

echo "Waiting for ClickHouse..."
echo "$CLICKHOUSE_DSN"

until goose -dir ./migrations clickhouse "$CLICKHOUSE_DSN" status >/dev/null 2>&1; do
  sleep 2
done

echo "Running migrations..."
goose -dir ./migrations clickhouse "$CLICKHOUSE_DSN" up

echo "Starting service..."
exec ./analytics-service
