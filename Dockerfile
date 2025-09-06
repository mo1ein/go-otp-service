# Build stage
FROM golang:1.23-bookworm AS builder

# Install dependencies
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    ca-certificates \
    tzdata \
    postgresql-client && \
    rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies with retry and timeout settings
RUN go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w GOSUMDB=off \
    && go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/go-otp-server ./cmd/server

# Install goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Final stage
FROM debian:bookworm-slim

# Install security updates and CA certificates
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    ca-certificates \
    tzdata \
    postgresql-client \
    curl && \
    rm -rf /var/lib/apt/lists/*


# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/go-otp-server /app/main

# Copy migrations files
COPY  migrations /app/migrations
COPY  scripts /app/scripts

# Make scripts executable
RUN chmod +x /app/scripts/*.sh

# Copy goose binary
COPY --from=builder  /go/bin/goose /usr/local/bin/goose

# Copy environment file template (if it exists)
COPY .env ./

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1


# Run the application
CMD ["/app/scripts/start.sh"]