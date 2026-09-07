# MyDrive

A personal Google Drive clone built with Go, PostgreSQL, Redis, MinIO, and Next.js.

> **⚠️ Learning Project**: This is an educational project created to understand distributed storage systems, cloud architecture, and full-stack development patterns. It is **not production-ready** and should be used for learning purposes only.
>
> Use this project to explore:
> - Go backend architecture without frameworks (net/http)
> - Repository → Service → Handler layer separation
> - PostgreSQL with pgx/v5
> - Docker containerization
> - Next.js with TypeScript
> - S3-compatible storage (MinIO)

## Architecture

```
mydrive/
├── apps/
│   ├── api/          # Go backend (net/http)
│   └── web/          # Next.js frontend
├── infrastructure/   # Docker configurations
├── migrations/       # Database migrations
└── packages/         # Shared packages (future)
```

### Technology Stack

**Backend:**
- Go 1.22+ with `net/http` (no frameworks)
- PostgreSQL 16 (database)
- Redis 7 (cache)
- MinIO (S3-compatible object storage)

**Frontend:**
- Next.js 14+ with App Router
- TypeScript
- TailwindCSS

**Infrastructure:**
- Docker & Docker Compose
- Environment-based configuration

## Prerequisites

- Docker & Docker Compose
- Make (optional, for convenience commands)
- Go 1.22+ (for local development without Docker)
- Node.js 18+ (for local frontend development)

## Quick Start

### 1. Clone and Setup

```bash
cd mydrive
cp .env.example .env
```

### 2. Start Services

Using Docker Compose:

```bash
make up
```

Or manually:

```bash
docker compose up -d
```

### 3. Access Services

- **API**: http://localhost:8080
- **Frontend**: http://localhost:3000
- **MinIO Console**: http://localhost:9001 (user: `minioadmin` / password: `minioadmin`)

### 4. Check Health

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{
  "status": "ok",
  "services": {
    "database": true,
    "redis": true,
    "storage": true
  }
}
```

## Development Commands

```bash
make help          # Show all available commands
make dev           # Start containers with logs
make up            # Start containers in background
make down          # Stop containers
make logs          # View logs
make ps            # Show running containers
make api-logs      # API logs only
make web-logs      # Frontend logs only
make clean         # Remove containers and volumes
```

## Useful Shortcuts

```bash
make api-shell     # SSH into API container
make db-shell      # Connect to PostgreSQL
make redis-cli     # Connect to Redis CLI
```

## Environment Variables

See `.env.example` for all available configuration options. Key variables:

- `API_PORT` - API server port (default: 8080)
- `DATABASE_URL` - PostgreSQL connection string
- `REDIS_URL` - Redis connection string
- `MINIO_*` - MinIO configuration
- `NEXT_PUBLIC_API_URL` - Frontend API endpoint

## Project Structure

### Backend (`apps/api/`)

```
cmd/server/main.go       # Entry point
internal/
  ├── config/            # Configuration loading
  ├── storage/           # Storage abstraction (MinIO)
  ├── auth/              # Authentication
  ├── users/             # User management
  ├── files/             # File operations
  ├── folders/           # Folder management
  ├── uploads/           # Upload handling
  ├── shares/            # Sharing functionality
  └── versions/          # File versioning
```

### Backend Architecture (per module)

Each module follows the **Repository → Service → Handler** pattern:

```
internal/folders/
├── models.go          # Data structures
├── repository.go      # Data access (SQL queries)
├── service.go         # Business logic & validation
├── handlers.go        # HTTP request/response handling
└── routes.go          # Route registration
```

**Flow**: `HTTP Request → Handler → Service (logic) → Repository (DB) → Response`

This separation ensures:
- 🗄️ **Repository**: All database queries in one place
- ⚙️ **Service**: Business logic is testable without HTTP
- 📡 **Handler**: Clean HTTP concerns (parsing, encoding)
- 🧩 **Models**: Pure data structures

### Frontend (`apps/web/`)

```
app/
  ├── layout.tsx         # Root layout
  ├── page.tsx           # Home page (health check demo)
  └── globals.css        # Global styles
```

## Database Migrations

Migrations are stored in `/migrations` as numbered SQL files:

```
migrations/
└── 001_init.sql        # Initial schema
```

To run migrations locally (future):

```bash
make migrate
```

## Storage Interface

The backend uses an abstraction for storage operations. Currently implemented with MinIO, but designed to support multiple backends:

```go
type Storage interface {
    Put(ctx context.Context, key string, r io.Reader, size int64, contentType string) error
    Get(ctx context.Context, key string) (io.ReadCloser, error)
    Delete(ctx context.Context, key string) error
    Exists(ctx context.Context, key string) (bool, error)
}
```

Additional implementations (local filesystem, S3, etc.) can be added without changing caller code.

## Roadmap

This is the initial scaffold. Planned features:

- [ ] User authentication & authorization
- [ ] File upload/download
- [ ] Folder management
- [ ] File sharing
- [ ] Version history
- [ ] Search functionality
- [ ] Activity logs
- [ ] Trash/recycle bin
- [ ] Sync client (future)
- [ ] Distributed storage (future)

## Development Notes

- No external frameworks (Go uses `net/http`)
- Configuration via environment variables
- Context.Context used throughout
- Database driver: `pgx/v5`
- Redis client: `go-redis/v9`
- Storage client: `minio-go/v7`

## Troubleshooting

### Containers fail to start

Check logs:
```bash
docker compose logs
```

### API can't connect to database

Ensure PostgreSQL is healthy:
```bash
docker compose ps
```

### Frontend can't reach API

Check `NEXT_PUBLIC_API_URL` environment variable and CORS settings in API.

### MinIO bucket not created

Check MinIO init container:
```bash
docker compose logs minio-init
```

## License

Personal learning project - use for educational purposes.

## References

- [Go `net/http` documentation](https://pkg.go.dev/net/http)
- [PostgreSQL pgx driver](https://github.com/jackc/pgx)
- [MinIO Go SDK](https://github.com/minio/minio-go)
- [Next.js documentation](https://nextjs.org/docs)
