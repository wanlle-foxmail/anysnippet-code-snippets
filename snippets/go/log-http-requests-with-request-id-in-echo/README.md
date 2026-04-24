# Log HTTP Requests with Request ID in Echo

Log one Echo request with method, path, status, duration, client IP, and request ID.

This snippet is useful when an API already has a request ID and you want one structured access log line per request.

## Highlights

- Logs one structured access line
- Reuses request or response IDs
- Preserves handler error flow

## What It Does

- Runs the next Echo handler first
- Reads `X-Request-ID` from the response header or request header
- Logs method, path, status, duration, and client IP with `log/slog`
- Preserves the original handler error so Echo can render it normally
- Falls back to `500` in logs for generic handler errors

## Usage

```go
// Run directly:
// go run request_logger_with_id.go
// Then call GET http://localhost:8080/hello with or without X-Request-ID.
```

## Notes

- This snippet logs a request ID but does not generate one.
- If another middleware writes `X-Request-ID` into the response header, that value wins.
- Passing `nil` uses `slog.Default()`.

## Verification

Run the tests from the snippet root:

```bash
go test -race ./...
```

The verified test suite covers:

- successful requests with an incoming request ID
- response-generated request IDs
- missing request IDs
- explicit `HTTPError` statuses
- generic handler errors logged as `500`

## Files

- `request_logger_with_id.go`
- `request_logger_with_id_test.go`
- `snippet.json`