# Request ID Middleware with Echo

Ensure every Echo request has an `X-Request-ID` header and a matching context value.

This snippet is useful when API handlers, logs, or support workflows need a stable request identifier that survives across middleware and handlers.

## Highlights

- Reuses incoming request IDs
- Generates IDs when missing
- Stores ID in Echo context

## What It Does

- Builds an Echo app with one request ID middleware
- Reuses the incoming `X-Request-ID` header when present
- Generates a random hex ID when the header is missing
- Writes the same ID back to the response header
- Stores the ID in `echo.Context` under `RequestIDContextKey`

## Usage

```go
// Run directly:
// go run request_id_middleware.go
// Then visit http://localhost:8080/hello
```

Response example:

```json
{
  "message": "ok",
  "request_id": "6f7ce56f63d24f4c8dcad70ed6b5a5c4"
}
```

## Notes

- Empty or whitespace-only `X-Request-ID` headers are treated as missing.
- Later handlers can read the same value with `c.Get(RequestIDContextKey)`.
- The snippet keeps one fixed header name and one fixed context key on purpose.

## Verification

Run the tests from the snippet root:

```bash
go mod tidy
go test ./...
```

The verified test suite covers:

- preserving an incoming request ID
- generating an ID when the header is missing
- replacing an empty request ID header
- replacing a whitespace-only request ID header
- always setting the response header
- exposing the same ID to the handler
- generating different IDs per request
- validating generated ID format
- returning HTTP 500 when ID generation fails

## Files

- `request_id_middleware.go`
- `request_id_middleware_test.go`
- `snippet.json`