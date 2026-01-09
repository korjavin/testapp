# Agent Guidelines for testapp

This document explains the project structure for AI assistants.

## Tech Stack

- **Language**: Go 1.22+
- **Database**: SQLite with sqlc for type-safe queries
- **Telegram**: telebot via telegram-bot-api/v5
- **Templating**: Go html/template

## Project Structure

```
testapp/
├── cmd/
│   ├── server/main.go    # Web server entry point
│   └── bot/main.go       # Telegram bot entry point
├── internal/
│   ├── api/              # HTTP handlers and routes
│   │   └── handlers.go   # Request handlers
│   ├── auth/             # Authentication middleware
│   │   └── telegram.go   # Telegram WebApp auth
│   ├── db/               # Database setup
│   │   └── db.go         # DB connection and migrations
│   ├── middleware/       # HTTP middleware
│   │   └── middleware.go # Logging, recovery, etc.
│   ├── repository/       # sqlc generated code
│   │   ├── db.go         # Database interface
│   │   └── queries/      # Generated query methods
│   ├── service/          # Business logic
│   │   └── services.go   # Core services
│   └── storage/          # Storage clients (future)
├── sql/
│   ├── migrations/       # Goose migrations
│   │   └── *.sql
│   └── queries/          # sqlc query files
│       └── *.sql
├── web/                  # Static files served by web server
│   ├── index.html
│   ├── css/
│   └── js/
└── Dockerfile
```

## Adding New Features

### Adding a New API Endpoint

1. Add route in `cmd/server/main.go` or create a new router file
2. Create handler in `internal/api/handlers.go`
3. Add sqlc queries in `sql/queries/` if DB access needed
4. Run `sqlc generate` to update repository layer

### Adding a New Telegram Command

1. Add handler in `cmd/bot/main.go`
2. Create service method in `internal/service/` if needed
3. Use repository layer for data access

### Database Migrations

1. Create new SQL file in `sql/migrations/` (e.g., `002_add_feature.sql`)
2. Use [Goose](https://pressly.github.io/goose/) syntax
3. Migrations run automatically on app startup

## Key Files

- `internal/db/db.go` - Database connection with WAL mode and migrations
- `internal/auth/telegram.go` - Telegram WebApp initData validation
- `cmd/server/main.go` - HTTP server with graceful shutdown
- `cmd/bot/main.go` - Telegram bot polling

## Commands

```bash
# Generate sqlc code
sqlc generate

# Run database migrations
go run -tags 'sqlite_json1' ./cmd/server -migrate

# Run locally
go run ./cmd/server
go run ./cmd/bot
```
