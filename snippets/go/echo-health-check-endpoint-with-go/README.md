# Echo Health Check Endpoint with Go

A `/health` endpoint that checks database, Redis, and disk space, then returns aggregated service status with uptime.

## What It Does

- Defines an Echo app with a `GET /health` endpoint
- Calls `CheckDatabase()`, `CheckRedis()`, and `CheckDisk()` to probe each dependency
- `CheckDisk` uses `syscall.Statfs` to report available disk space in MB
- Aggregates results: if all checks return `"ok"`, the service is `"healthy"`; otherwise `"unhealthy"`
- Includes server uptime as a top-level field
- Stub check functions for DB and Redis are provided — replace the body with real connection logic

## Usage

```go
// Run directly:
// go run health_check.go
// Then visit http://localhost:8080/health
```

Response when all dependencies are healthy:

```json
{
  "status": "healthy",
  "uptime": "2h30m15s",
  "checks": {
    "database": {"status": "ok"},
    "redis": {"status": "ok"},
    "disk": {"status": "ok", "available": "51200 MB"}
  }
}
```

## How to Extend

Replace the stub check functions with real logic:

```go
var CheckDatabase = func() map[string]interface{} {
    if err := db.Ping(); err != nil {
        return map[string]interface{}{"status": "error", "detail": err.Error()}
    }
    return map[string]interface{}{"status": "ok"}
}
```

## Verification

```bash
cd snippets/go/echo-health-check-endpoint-with-go
go mod tidy
go test -v ./...
```
