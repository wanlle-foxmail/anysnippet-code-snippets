# Run Parallel Tasks with errgroup in Go

Run several context-aware tasks in parallel in Go and cancel siblings on the first error with `errgroup`.

This snippet is useful when one request or worker should fan out to multiple tasks and fail fast if any one of them fails.

## Highlights

- Cancels sibling tasks on error
- Waits for every goroutine cleanly
- Uses one shared group context

## What It Does

- validates the task list before starting goroutines
- starts each task in one `errgroup.Group`
- shares one derived context across all tasks
- returns the first task error
- lets sibling tasks observe cancellation through `ctx.Done()`

## Usage

```go
// Run directly:
// go run run_parallel_tasks.go
// The example starts two short tasks under one errgroup.
```

## Notes

- A nil parent context falls back to `context.Background()`.
- Tasks should return promptly when `ctx.Done()` is closed.
- An empty task list is treated as a successful no-op.

## Verification

Run the tests from the snippet root:

```bash
go test -race ./...
```

The verified test suite covers:

- all-success runs
- first-error returns
- sibling cancellation
- nil parent contexts
- nil tasks
- empty task lists

## Files

- `run_parallel_tasks.go`
- `run_parallel_tasks_test.go`
- `snippet.json`