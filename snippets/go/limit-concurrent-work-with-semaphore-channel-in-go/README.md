# Limit Concurrent Work with a Semaphore Channel in Go

Limit concurrent Go work with a buffered channel semaphore and explicit acquire and release helpers.

This snippet is useful when you want to cap the number of in-flight goroutines without building a larger worker pool abstraction.

## Highlights

- Caps in-flight work with one channel
- Supports blocking and non-blocking acquire
- Returns clear release errors

## What It Does

- creates a buffered channel semaphore with a fixed capacity
- acquires permits with optional context cancellation
- supports non-blocking `TryAcquire`
- releases permits explicitly after work finishes
- returns an error if code releases more permits than it acquired

## Usage

```go
// Run directly:
// go run semaphore_channel.go
// The example processes three jobs while allowing only two to run at once.
```

## Notes

- `Acquire` blocks until a permit is available or the context is canceled.
- `TryAcquire` is useful when you want to reject or defer work instead of waiting.
- `Release` returns an error on over-release instead of panicking.

## Verification

Run the tests from the snippet root:

```bash
go test -race ./...
```

The verified test suite covers:

- rejecting invalid capacity values
- enforcing the maximum concurrency limit
- making permits available again after release
- reporting a full semaphore with `TryAcquire`
- honoring context cancellation while waiting
- rejecting release without an acquired permit

## Files

- `semaphore_channel.go`
- `semaphore_channel_test.go`
- `snippet.json`