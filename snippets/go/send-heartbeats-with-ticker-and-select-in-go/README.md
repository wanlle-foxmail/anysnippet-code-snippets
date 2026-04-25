# Send Heartbeats with Ticker and select in Go

Send periodic heartbeat timestamps from Go with one `time.Ticker`, one output channel, and `select`.

This snippet is useful when background work should publish liveness signals without blocking on slow heartbeat consumers.

## Highlights

- Emits periodic heartbeat timestamps
- Drops extra ticks when receivers lag
- Closes output on context cancel

## What It Does

- starts one ticker for the heartbeat interval
- forwards tick timestamps into a buffered output channel
- drops extra heartbeats when the receiver falls behind
- stops the ticker when the context is canceled
- closes the output channel on shutdown

## Usage

```go
// Run directly:
// go run heartbeat_ticker.go
// The example logs three heartbeat timestamps.
```

## Notes

- The output channel is buffered with size 1 so slow receivers do not block the heartbeat loop.
- Dropping old heartbeats is often acceptable because the signal is only about liveness.
- A nil context falls back to `context.Background()`.

## Verification

Run the tests from the snippet root:

```bash
go test -race ./...
```

The verified test suite covers:

- rejecting an invalid interval
- emitting heartbeat ticks
- not emitting before the first interval ends
- closing the output channel on cancellation
- dropping extra heartbeats when the receiver lags
- emitting multiple ticks when the receiver drains the channel

## Files

- `heartbeat_ticker.go`
- `heartbeat_ticker_test.go`
- `snippet.json`