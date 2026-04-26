# Initialize Shared Dependencies with sync.OnceValues in Go

Initialize one shared dependency in Go with `sync.OnceValues` so every caller sees the same first result.

This snippet is useful when many goroutines may need a lazily initialized client, config object, or connection string at the same time.

## Highlights

- Runs the initializer once
- Replays the first value or error
- Handles concurrent callers safely

## What It Does

- wraps one initializer with `sync.OnceValues`
- caches the first successful value for later callers
- caches the first error instead of retrying silently
- replays panics the same way on later calls
- rejects nil initializer functions with a clear error

## Usage

```go
// Run directly:
// go run sync_once_values.go
// The example lazily loads one AppSettings value.
```

## Notes

- This pattern fits dependencies that should initialize at most once per process.
- If the first result is an error, later callers see the same error until the process restarts.
- Use a different pattern when you need retryable or refreshable initialization.

## Verification

Run the tests from the snippet root:

```bash
go test -race ./...
```

The verified test suite covers:

- cached success values
- cached errors
- concurrent callers
- nil initializer functions
- repeated panics

## Files

- `sync_once_values.go`
- `sync_once_values_test.go`
- `snippet.json`