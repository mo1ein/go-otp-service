#!/bin/bash
# start.sh

set -e  # Exit on any error

# Set default values for environment variables
POSTGRES_HOST=${DATABASE_POSTGRES_HOST:-db}
POSTGRES_USER=${DATABASE_POSTGRES_USER:-go-otp-service}
POSTGRES_PASSWORD=${DATABASE_POSTGRES_PASSWORD:-go-otp-service_password}
POSTGRES_DB=${DATABASE_POSTGRES_NAME:-go-otp-service}
POSTGRES_SSL_MODE=${POSTGRES_SSL_MODE:-disable}

# Debug: Print the values
echo "DEBUG: POSTGRES_HOST=$POSTGRES_HOST"
echo "DEBUG: POSTGRES_USER=$POSTGRES_USER"
echo "DEBUG: POSTGRES_DB=$POSTGRES_DB"
echo "DEBUG: POSTGRES_PASSWORD=$POSTGRES_PASSWORD"

echo "Waiting for PostgreSQL to be ready..."
# Wait for PostgreSQL to be available
until PGPASSWORD=$POSTGRES_PASSWORD psql -h "$POSTGRES_HOST" -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c '\q' > /dev/null 2>&1; do
  echo "PostgreSQL is unavailable - sleeping"
  sleep 1
done

echo "PostgreSQL is up - running migrations..."

# Run database migrations with Goose
/usr/local/bin/goose -dir /app/migrations postgres "host=$POSTGRES_HOST user=$POSTGRES_USER dbname=$POSTGRES_DB password=$POSTGRES_PASSWORD sslmode=$POSTGRES_SSL_MODE" up

if [ $? -eq 0 ]; then
    echo "Migrations completed successfully"
else
    echo "Migrations failed" >&2
    exit 1
fi

echo "Starting application..."
# Start the main application
exec ./main