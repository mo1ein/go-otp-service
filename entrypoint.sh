#!/bin/bash
# entrypoint.sh

set -e

# Set default values if environment variables are not set
POSTGRES_HOST=${DATABASE_POSTGRES_HOST:-${POSTGRES_HOST:-db}}
POSTGRES_USER=${DATABASE_POSTGRES_USER:-${POSTGRES_USER:-go-otp-service}}
POSTGRES_PASSWORD=${DATABASE_POSTGRES_PASSWORD:-${POSTGRES_PASSWORD:-go-otp-service_password}}
POSTGRES_DB=${DATABASE_POSTGRES_NAME:-${POSTGRES_DB:-go-otp-service}}
POSTGRES_SSL_MODE=${POSTGRES_SSL_MODE:-disable}

# Wait for database to be ready
echo "Waiting for database to be ready..."
echo "Connecting to database with: host=$POSTGRES_HOST user=$POSTGRES_USER db=$POSTGRES_DB"
until PGPASSWORD=$POSTGRES_PASSWORD psql -h "$POSTGRES_HOST" -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c '\q'; do
  echo "Postgres is unavailable - sleeping"
  sleep 1
done

echo "Postgres is up - running migrations"

# Run migrations with goose
goose -allow-missing -dir ./migrations postgres "host=$POSTGRES_HOST user=$POSTGRES_USER dbname=$POSTGRES_DB password=$POSTGRES_PASSWORD sslmode=$POSTGRES_SSL_MODE" up

if [ $? -eq 0 ]; then
    echo "Migrations completed successfully"
else
    echo "Migrations failed" >&2
    exit 1
fi

# Start the application
echo "Starting application..."
exec ./main