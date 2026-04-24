# Deduplicate POST Requests with Idempotency Keys in Echo

Deduplicate one Echo `POST` endpoint with an `Idempotency-Key` header and a single-process in-memory store.

This snippet is useful when a client can retry the same create request and you want to replay the first successful response instead of running the handler twice.

## Highlights

- Replays successful POST responses
- Returns `409` while work runs
- Keeps storage in one process

## What It Does

- Requires `Idempotency-Key` on `POST` requests
- Replays a cached successful response for a repeated key
- Returns `409 Conflict` when the same key is already being processed
- Does not cache failed responses
- Uses a simple mutex-protected in-memory store for one process

## Usage

```go
// Run directly:
// go run idempotency_middleware.go
// Then call POST http://localhost:8080/orders with Idempotency-Key: your-key.
```

## Notes

- This snippet is intentionally limited to a single process and does not share state across instances.
- Only successful `2xx` responses are cached and replayed.
- `GET` and other non-`POST` requests bypass the idempotency logic.
- Successful responses stay in memory until the process exits; add TTL or cleanup before using this pattern in a long-running service.
- Idempotency keys are scoped by HTTP method, request path, and caller context when `X-User-ID` or `Authorization` is present.
- Anonymous requests share one caller scope; if you need per-client isolation, derive the scope from your auth or session layer.
- For multi-user APIs, wire the caller scope helper to your authenticated subject or tenant identity before using this pattern in production.

## Verification

Run the tests from the snippet root:

```bash
go test -race ./...
```

The verified test suite covers:

- missing idempotency keys on `POST`
- caching and replaying successful responses
- keeping different keys separate
- replaying a response that completes between store checks
- separating different caller scopes with the same key
- normalizing the `Authorization` header scope for the same caller
- returning `409` for in-progress duplicates
- retrying after failed responses
- bypassing non-`POST` requests

## Files

- `idempotency_middleware.go`
- `idempotency_middleware_test.go`
- `snippet.json`