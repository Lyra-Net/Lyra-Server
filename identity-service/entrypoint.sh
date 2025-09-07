#!/bin/sh
set -e

echo "Waiting for postgres..."
until pg_isready -h "$POSTGRES_HOST" -p "$POSTGRES_PORT" -U "$POSTGRES_USER"; do
  sleep 2
done

echo "Running migrations..."
migrate -path ./migration -database "$DB_URL" up

echo "Starting service..."
exec ./identity-service
