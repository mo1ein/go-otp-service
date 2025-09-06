.PHONY: build up down logs clean migration-create migration-up migration-down migration-status dev test build-goose

# Load environment variables from .env file
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# Docker commands
build:
	docker compose build

up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f

clean:
	docker compose down -v
	docker system prune -f

# Migration commands for local development
migration-create:
	@read -p "Enter migration name: " name; \
	goose -dir migrations create $${name} sql

migration-up:
	goose -dir migrations postgres "host=${DATABASE_POSTGRES_HOST} port=${DATABASE_POSTGRES_PORT} user=${DATABASE_POSTGRES_USER} dbname=${DATABASE_POSTGRES_NAME} password=${DATABASE_POSTGRES_PASSWORD} sslmode=${POSTGRES_SSL_MODE}" up

migration-down:
	goose -dir migrations postgres "host=${DATABASE_POSTGRES_HOST} port=${DATABASE_POSTGRES_PORT} user=${DATABASE_POSTGRES_USER} dbname=${DATABASE_POSTGRES_NAME} password=${DATABASE_POSTGRES_PASSWORD} sslmode=${POSTGRES_SSL_MODE}" down

migration-status:
	goose -dir migrations postgres "host=${DATABASE_POSTGRES_HOST} port=${DATABASE_POSTGRES_PORT} user=${DATABASE_POSTGRES_USER} dbname=${DATABASE_POSTGRES_NAME} password=${DATABASE_POSTGRES_PASSWORD} sslmode=${POSTGRES_SSL_MODE}" status

# For local development without Docker
dev:
	go run ./cmd/server

test:
	go test ./... -v

# Build goose binary locally
build-goose:
	go build -o goose ./vendor/github.com/pressly/goose/v3/cmd/goose