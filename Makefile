# Makefile
.PHONY: build up down logs clean migration-create migration-up migration-down

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
	goose -dir migrations postgres "host=0.0.0.0 port=5433 user=go-otp-service dbname=go-otp-service password=go-otp-service_password sslmode=disable" up

migration-down:
	goose -dir migrations postgres "host=localhost user=go-otp-service dbname=go-otp-service password=go-otp-service_password sslmode=disable" down

migration-status:
	goose -dir migrations postgres "host=localhost user=go-otp-service dbname=go-otp-service password=go-otp-service_password sslmode=disable" status

# For local development without Docker
dev:
	go run ./cmd/server

test:
	go test ./... -v

# Build goose binary locally
build-goose:
	go build -o goose ./vendor/github.com/pressly/goose/v3/cmd/goose