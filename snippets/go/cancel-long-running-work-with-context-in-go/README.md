# Cancel Long-Running Work with Context in Go

Run batch work in Go while checking `context.Context` cancellation between units of work.

This snippet is useful when one long-running loop should stop cleanly without starting more work after a shutdown or timeout signal.

## Highlights

- Checks cancellation before each item
- Passes context into the worker callback
- Wraps item-specific worker errors

## What It Does

- iterates through one batch of string items
- checks `ctx.Done()` before starting each item
- passes the same context into the worker callback
- stops immediately when the context is canceled
- wraps worker errors with the item that failed

## Usage

```go
// Run directly:
// go run process_batch_with_context.go
// The example processes three invoice IDs with one callback.
```

## Notes

- Cancellation is cooperative, so the callback should also respect the context if it blocks on I/O.
- The helper checks cancellation between items instead of spawning internal goroutines.
- A nil context falls back to `context.Background()`.

## Verification

Run the tests from the snippet root:

```bash
go test -race ./...
```

The verified test suite covers:

- processing every item when the context stays active
- stopping after cancellation
- wrapping worker errors with the failed item
- rejecting a nil worker function
- accepting a nil context
- handling an empty batch

## Files

- `process_batch_with_context.go`
- `process_batch_with_context_test.go`
- `snippet.json`