# Reset Idle Timeouts with time.Timer in Go

Reset a `time.Timer` on activity so Go code can detect idle timeouts without recreating timers.

This snippet is useful when a connection, worker, or stream should stay alive while activity continues and fail only after a quiet period.

## Highlights

- Resets one timer after each activity
- Drains the timer channel safely
- Returns a clear idle timeout error

## What It Does

- starts one `time.Timer` for the idle window
- watches for activity, timeout, or context cancellation
- stops and drains the timer before each reset
- returns `ErrIdleTimeout` when the idle window expires
- exits cleanly when the activity channel closes

## Usage

```go
// Run directly:
// go run idle_timeout_timer.go
// The example sends one activity signal, resets the timer, and then exits cleanly.
```

## Notes

- The `stopAndDrainTimer` helper avoids stale timer events after a reset.
- A nil activity channel behaves like "no activity" and will eventually time out.
- Closing the activity channel signals that the monitored work finished before the idle timeout.

## Verification

Run the tests from the snippet root:

```bash
go test -race ./...
```

The verified test suite covers:

- timing out without activity
- resetting after one activity signal
- resetting after multiple activity signals
- returning nil when the activity channel closes
- honoring context cancellation
- rejecting invalid timeout values

## Files

- `idle_timeout_timer.go`
- `idle_timeout_timer_test.go`
- `snippet.json`