# Handle Signals with signal.NotifyContext in Go

Bridge OS signals into a Go context with `signal.NotifyContext` so shutdown code can wait on `ctx.Done()`.

This snippet is useful when a command or service should respond to one interrupt signal without wiring a separate signal channel by hand.

## Highlights

- Turns signals into context cancellation
- Supports parent context propagation
- Uses one explicit stop function

## What It Does

- validates that at least one signal is provided
- falls back to `context.Background()` when the parent is nil
- derives one context from `signal.NotifyContext`
- lets the caller stop signal delivery with the returned cancel function
- unifies signal handling with other `ctx.Done()` shutdown paths

## Usage

```go
// Run directly:
// go run signal_notify_context.go
// The example waits for one interrupt signal and then exits.
```

## Notes

- Always call the returned stop function so signal notifications are released.
- Parent cancellation and signal delivery both close the derived context.
- The local test suite verifies real interrupt handling in a helper process on macOS and Linux.

## Verification

Run the tests from the snippet root:

```bash
go test -race ./...
```

The verified test suite covers:

- rejecting an empty signal list
- creating a context from a nil parent
- honoring parent cancellation
- canceling the context when `stop` is called
- canceling the context when a helper process receives `os.Interrupt`

## Files

- `signal_notify_context.go`
- `signal_notify_context_test.go`
- `snippet.json`