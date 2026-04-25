# Cancel External Commands with Context in Go

Run an external command from Go with `exec.CommandContext` so timeouts and cancellations stop the process cleanly.

This snippet is useful when one command should respect the same `context.Context` lifecycle as the rest of your request or shutdown flow.

## Highlights

- Uses exec.CommandContext directly
- Maps process cancellation to context errors
- Returns combined stdout and stderr

## What It Does

- validates the command name
- falls back to `context.Background()` for a nil context
- starts the process with `exec.CommandContext`
- returns combined stdout and stderr output
- reports context cancellation or timeout explicitly when the process is stopped

## Usage

```go
// Run directly:
// go run command_context.go
// The example runs `go version` and logs the combined output.
```

## Notes

- `CombinedOutput` is convenient when you want one error path with the full process output.
- Wrapping `ctx.Err()` keeps timeout and cancellation handling clearer than a generic killed-process error.
- The local test suite uses the current test binary as a helper process so no shell dependency is required.

## Verification

Run the tests from the snippet root:

```bash
go test -race ./...
```

The verified test suite covers:

- returning combined output for a successful command
- returning combined stderr output for a failing command
- canceling a long-running command on timeout
- honoring parent cancellation
- rejecting an empty command name
- accepting a nil context

## Files

- `command_context.go`
- `command_context_test.go`
- `snippet.json`