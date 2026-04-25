# Graceful Background Worker Shutdown with Go

Run a buffered background worker in Go that drains accepted jobs, rejects new submissions during shutdown, and runs cleanup once.

This snippet is useful when one in-process worker should finish accepted work during shutdown without pretending it can still accept more jobs.

## Highlights

- Drains accepted jobs before exit
- Rejects new jobs during shutdown
- Runs cleanup exactly once

## What It Does

- starts one buffered worker goroutine
- accepts jobs with a non-blocking `Submit`
- rejects new jobs after shutdown starts
- closes the queue and waits for accepted jobs to finish
- respects the shutdown context deadline while waiting

## Usage

```go
// Run directly:
// go run background_worker.go
// The example submits two jobs and then shuts the worker down.
```

## Notes

- `Submit` returns `ErrWorkerQueueFull` instead of blocking when the queue is full.
- `Shutdown` is idempotent and only closes the queue once.
- Cleanup runs after the worker loop finishes draining the accepted jobs.

## Verification

Run the tests from the snippet root:

```bash
go test -race ./...
```

The verified test suite covers:

- draining accepted jobs during shutdown
- rejecting new submissions after shutdown starts
- returning `ErrWorkerQueueFull` for a full queue
- running cleanup once
- respecting a shutdown deadline
- rejecting invalid constructor input

## Files

- `background_worker.go`
- `background_worker_test.go`
- `snippet.json`