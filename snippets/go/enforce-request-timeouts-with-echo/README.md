# Enforce Request Timeouts with Echo

Cap one Echo request with a fixed deadline and return `503` when the handler runs too long.

This snippet is useful when an API route should stop waiting on slow work and tell the client the request timed out.

## Highlights

- Applies one fixed request deadline
- Returns `503` on slow handlers
- Drops late handler writes

## What It Does

- Wraps one Echo handler chain in `context.WithTimeout`
- Buffers handler output until the handler finishes
- Flushes the buffered response when the handler finishes in time
- Returns a timeout response when the deadline fires first
- Discards writes that happen after the timeout response was already sent

## Usage

```go
// Run directly:
// go run request_timeout_middleware.go
// Then visit http://localhost:8080/hello
```

Timeout response example:

```json
{
  "message": "request timed out"
}
```

## Notes

- This middleware still depends on the handler honoring `request.Context()` for timely cancellation.
- A non-positive timeout is treated as a server configuration error.
- The snippet buffers normal handler output so late writes do not leak into the client response.
- Buffered responses stay in memory until the handler finishes, so add body size limits before using this approach for large payloads.

## Verification

Run the tests from the snippet root:

```bash
go test -race ./...
```

The verified test suite covers:

- returning a normal response inside the deadline
- returning `503` when the deadline fires first
- canceling the request context on timeout
- keeping fast and slow requests independent
- rejecting a non-positive timeout
- preserving handler errors

## Files

- `request_timeout_middleware.go`
- `request_timeout_middleware_test.go`
- `snippet.json`