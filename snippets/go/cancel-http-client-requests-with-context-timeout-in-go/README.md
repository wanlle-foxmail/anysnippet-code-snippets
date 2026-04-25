# Cancel HTTP Client Requests with Context Timeout in Go

Send an outbound HTTP GET request with a per-request context timeout so slow servers are canceled cleanly.

This snippet is useful when a Go service calls another HTTP endpoint and should stop waiting after one fixed deadline.

## Highlights

- Adds one timeout per outbound request
- Honors parent context cancellation
- Uses request-scoped context propagation

## What It Does

- validates the client, URL, and timeout inputs
- derives one request-scoped timeout from the parent context
- builds the request with `http.NewRequestWithContext`
- cancels the request when the timeout or parent cancellation fires
- returns the final response to the caller on success

## Usage

```go
// Run directly:
// go run http_client_timeout.go
// The example sends one GET request to https://example.com/health.
```

## Notes

- The caller owns the final response body and should close it.
- Parent context cancellation wins over the per-request timeout.
- This snippet is intentionally limited to one GET helper so the cancellation pattern stays obvious.

## Verification

Run the tests from the snippet root:

```bash
go test -race ./...
```

The verified test suite covers:

- returning a successful response before the timeout
- canceling a slow request when the deadline expires
- honoring parent context cancellation
- rejecting a nil client
- rejecting an empty URL
- rejecting an invalid timeout

## Files

- `http_client_timeout.go`
- `http_client_timeout_test.go`
- `snippet.json`