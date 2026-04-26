# Retry Operations with Exponential Backoff and Jitter in Go

Retry one operation in Go with exponential backoff, jitter, and retryable error filtering.

This snippet is useful when a transient failure should be retried with increasing delays while still respecting context cancellation.

## Highlights

- Retries only selected errors
- Adds jitter to each delay
- Stops on context cancellation

## What It Does

- validates retry callbacks and backoff settings
- runs one context-aware operation up to a fixed attempt limit
- retries only when the error matches a retryable predicate
- waits with exponential backoff plus bounded jitter
- stops early when the context is canceled during a delay

## Usage

```go
// Run directly:
// go run retry_backoff.go
// The example retries one temporary failure until it succeeds.
```

## Notes

- Set `JitterFraction` to `0` when you need deterministic backoff spacing.
- The snippet returns the last retryable error once it reaches `MaxAttempts`.
- The wait helper is separated so tests can stay fast and deterministic.
- Large attempt counts clamp to `MaxDelay` before converting the computed delay.

## Verification

Run the tests from the snippet root:

```bash
go test -race ./...
```

The verified test suite covers:

- first-attempt success
- retry before success
- non-retryable failures
- retry exhaustion
- context cancellation during backoff
- invalid input
- large-attempt delay clamping

## Files

- `retry_backoff.go`
- `retry_backoff_test.go`
- `snippet.json`