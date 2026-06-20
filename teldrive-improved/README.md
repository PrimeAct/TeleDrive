# TelDrive Improved

A production-ready, improved version of TelDrive with clean architecture, better performance, and enhanced maintainability.

## Features

- Clean Architecture with internal/pkg separation
- Unified Cache Interface (Memory and Redis backends)
- Structured Logging with rotation
- Graceful Shutdown with context-based lifecycle
- JWT Authentication with middleware
- CORS Support
- Health Checks with metrics
- Database Migrations
- Docker Multi-stage builds
- CI/CD Ready with GitHub Actions

## Quick Start

```bash
# Clone and build
git clone https://github.com/tgdrive/teldrive-improved.git
cd teldrive-improved
make build

# Run with config
./bin/teldrive run --config config.toml

# Or with Docker
docker-compose -f docker/docker-compose.yml up -d
```

## Configuration

Create `config.toml`:

```toml
[server]
host = "0.0.0.0"
port = 8080

[log]
level = "info"
format = "json"

[db]
host = "localhost"
port = 5432
user = "postgres"
password = "postgres"
name = "teldrive"

[tg]
app-id = 12345
app-hash = "your-app-hash"
phone = "+1234567890"

[jwt]
secret = "your-secret-key"
session-time = "720h"
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| POST | `/auth/login` | Login |
| GET | `/files` | List files |
| POST | `/files` | Create file |
| GET | `/files/{id}` | Get file |
| DELETE | `/files/{id}` | Delete file |
| POST | `/uploads` | Start upload |
| POST | `/uploads/{id}/chunk` | Upload chunk |
| POST | `/uploads/{id}/complete` | Complete upload |

## Architecture

```
HTTP API -> Services -> Repositories -> Database
    |          |            |
Middleware  Telegram     Cache
            Client
```

## License

MIT License
