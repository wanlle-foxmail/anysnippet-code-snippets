# Limit Request Rate with Token Bucket in Echo

Limit one Echo route per client IP with an in-memory token bucket and `429` responses.

This snippet is useful when one small API should allow short bursts but still reject repeated requests from the same client after the bucket is empty.

## Highlights

- Limits requests per client IP
- Returns 429 with Retry-After
- Refills tokens over time

## What It Does

- Builds one Echo middleware with a token bucket per client IP
- Allows a short burst up to the configured bucket capacity
- Refills tokens at a fixed interval
- Returns `429 Too Many Requests` when a bucket is empty
- Sets a `Retry-After` header for blocked requests

## Usage

```go
// Run directly:
// go run token_bucket_rate_limit.go
// Then call GET /hello more than twice from the same client.
```

## Notes

- This snippet keeps rate limiting state in one process only.
- Client buckets are not evicted, so production services should add idle-bucket cleanup, a maximum bucket count, or a shared limiter such as Redis.
- Client buckets are keyed by `c.RealIP()` and fall back to one shared `unknown` bucket when the client IP is missing.
- When the app runs behind a proxy, configure Echo's trusted proxy handling before relying on `c.RealIP()` for security-sensitive rate limits.
- Use a distributed store instead of in-memory state when multiple instances must share the same limit.

## Verification

Run the tests from the snippet root:

```bash
go mod tidy
go test -race ./...
```

The verified test suite covers:

- allowing requests within capacity
- returning `429` when a bucket is empty
- refilling tokens after the interval
- separating buckets per client IP
- handling missing client IPs
- panicking on invalid configuration at setup time

## Files

- `token_bucket_rate_limit.go`
- `token_bucket_rate_limit_test.go`
- `snippet.json`