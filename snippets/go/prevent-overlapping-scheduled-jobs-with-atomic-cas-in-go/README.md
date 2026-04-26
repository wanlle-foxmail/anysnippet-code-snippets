# Prevent Overlapping Scheduled Jobs with Atomic CAS in Go

Prevent overlapping scheduled job runs in Go with one atomic compare-and-swap guard.

This snippet is useful when a cron-like job should skip a tick instead of running the same work twice in parallel.

## Highlights

- Skips overlapping job runs
- Clears the flag on every exit
- Uses one atomic CAS check

## What It Does

- attempts to flip a shared running flag from false to true
- skips the job when another run is already active
- clears the running flag after success, error, or panic
- returns whether the current tick actually ran the job
- rejects nil guards and nil job functions

## Usage

```go
// Run directly:
// go run atomic_cas_job_guard.go
// The example runs one short job behind the guard.
```

## Notes

- This pattern is good for best-effort scheduled work that can skip one interval.
- The deferred flag reset lets a panic release the guard before it propagates.
- If skipped work must be queued instead of dropped, use a different coordination pattern.

## Verification

Run the tests from the snippet root:

```bash
go test -race ./...
```

The verified test suite covers:

- idle runs
- overlapping runs
- flag reset after success
- flag reset after errors
- flag reset after panics
- invalid input

## Files

- `atomic_cas_job_guard.go`
- `atomic_cas_job_guard_test.go`
- `snippet.json`