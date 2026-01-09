# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install git for go modules
RUN apk add --no-cache git

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the server
RUN CGO_ENABLED=0 GOOS=linux go build -o server -mod vendor ./cmd/server

# Build the bot
RUN CGO_ENABLED=0 GOOS=linux go build -o bot -mod vendor ./cmd/bot

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates bash

WORKDIR /app

# Copy binaries
COPY --from=builder /app/server .
COPY --from=builder /app/bot .

# Copy web static files
COPY web/ ./web/

# Copy migrations
COPY sql/migrations/ ./sql/migrations/

# Create data directory for SQLite
RUN mkdir -p /app/data

# Expose port
EXPOSE 8080

# Healthcheck
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Default command
CMD ["./server"]
