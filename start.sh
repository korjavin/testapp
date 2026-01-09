#!/bin/bash
set -e

echo "Starting testapp..."

# Ensure database directory exists
mkdir -p data

# Run database migrations
echo "Running database migrations..."
go run -tags 'sqlite_json1' ./cmd/server -migrate || true

# Start the web server in background
echo "Starting web server..."
go run ./cmd/server &

# Store PID
SERVER_PID=$!

# Function to cleanup on exit
cleanup() {
    echo "Shutting down..."
    kill $SERVER_PID 2>/dev/null || true
    exit 0
}

trap cleanup SIGINT SIGTERM

# Wait for server to be ready
sleep 2
echo "Server started on http://localhost:8080"

# Keep script running
wait
