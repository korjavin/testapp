# testapp

test go/web/tg app to test claude skill

## Quick Start

```bash
# Copy environment file
cp .env.example .env

# Edit .env with your configuration
nano .env

# Run locally
./start.sh
```

## Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `APP_ENV` | Environment (development/production) | No |
| `APP_PORT` | HTTP server port | No |
| `LOG_LEVEL` | Log level (debug/info/warn/error) | No |
| `DB_PATH` | SQLite database path | Yes |
| `TG_BOT_TOKEN` | Telegram bot token | Yes |
| `TG_WEBHOOK_URL` | Telegram webhook URL | No |
| `LITESTREAM_R2_ACCESS_KEY_ID` | R2 access key for Litestream | No |
| `LITESTREAM_R2_SECRET_ACCESS_KEY` | R2 secret key for Litestream | No |
| `LITESTREAM_R2_ACCOUNT_ID` | R2 account ID | No |
| `LITESTREAM_R2_BUCKET_NAME` | R2 bucket name for DB backup | No |
| `LITESTREAM_REPLICATE_URL` | Litestream replication URL | No |
| `COOKIE_SECRET` | Secret for session cookies | Yes |

## Development

### Running with Docker Compose

```bash
docker-compose up -d
```

### Running locally

```bash
# Start the web server
go run ./cmd/server

# In another terminal, start the bot
go run ./cmd/bot
```

## Deployment

### Docker

```bash
docker build -t testapp .
docker run -d --name testapp -p 8080:8080 --env-file .env testapp
```

### Docker Compose with Litestream

```bash
docker-compose up -d
```

## Project Structure

```
testapp/
├── cmd/
│   ├── server/     # Web server entry point
│   └── bot/        # Telegram bot entry point
├── internal/
│   ├── api/        # HTTP handlers
│   ├── auth/       # Authentication (Telegram WebApp)
│   ├── db/         # Database connection
│   ├── middleware/ # HTTP middleware
│   ├── repository/ # sqlc generated code
│   ├── service/    # Business logic
│   └── storage/    # Storage clients
├── sql/
│   ├── migrations/ # Database migrations
│   └── queries/    # sqlc query definitions
├── web/
│   ├── index.html
│   ├── css/
│   └── js/
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── start.sh
```
