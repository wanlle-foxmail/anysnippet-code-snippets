# Retry HTTP GET Requests with Go

Retry an HTTP GET request on transport errors, `429`, and `5xx` responses with a fixed delay.

This snippet is useful when one idempotent GET endpoint can fail transiently and you want a small retry helper without pulling in a larger HTTP client wrapper.

## Highlights

- Retries transport failures
- Retries `429` and `5xx`
- Closes retryable response bodies

## What It Does

- Sends one HTTP GET request with a provided `http.Client`
- Retries transport errors until attempts run out
- Retries `429 Too Many Requests` and `5xx` responses
- Returns `4xx` responses other than `429` without retrying
- Closes retryable response bodies before the next attempt

## Usage

```go
// Run directly:
// go run retry_http_get.go
// The example sends one GET request to https://example.com/health.
```

## Notes

- This snippet is intentionally limited to GET requests.
- The caller owns the final response body and should close it.
- Retryable status responses are closed before the next attempt starts.
- The helper uses a fixed delay on purpose and does not implement backoff or `Retry-After` parsing.

## Verification

Run the tests from the snippet root:

```bash
go test -race ./...
```

The verified test suite covers:

- immediate success on the first response
- retry after a `500` response
- retry after a `429` response
- no retry for `400 Bad Request`
- retry after a transport error
- closing retryable response bodies before retrying
- returning the final `5xx` response when retries run out
- rejecting invalid input

## Files

- `retry_http_get.go`
- `retry_http_get_test.go`
- `snippet.json`