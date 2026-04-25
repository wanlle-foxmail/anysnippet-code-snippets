# Coalesce Bursty Events with select in Go

Merge rapid Go events into fewer downstream batches with `select`, one timer, and a quiet window.

This snippet is useful when repeated updates such as save signals, invalidation events, or refresh requests should be collapsed before expensive downstream work runs.

## Highlights

- Merges rapid events into one batch
- Flushes pending events on input close
- Stops cleanly on context cancel

## What It Does

- collects incoming values into one in-memory batch
- starts or resets a timer for the quiet window
- emits one batch after the input stays quiet long enough
- flushes the remaining batch immediately when the input closes
- stops early when the context is canceled

## Usage

```go
// Run directly:
// go run coalesce_events.go
// The example turns three rapid save signals into one logged batch.
```

## Notes

- The helper preserves the arrival order inside each emitted batch.
- A nil input channel is rejected because the function would otherwise wait forever.
- The batch is copied before emission so later appends do not mutate previously emitted slices.

## Verification

Run the tests from the snippet root:

```bash
go test -race ./...
```

The verified test suite covers:

- rejecting a nil input channel
- rejecting an invalid coalescing window
- merging a quick burst into one batch
- splitting separated bursts into separate batches
- flushing a pending batch on input close
- closing cleanly for empty input and context cancellation

## Files

- `coalesce_events.go`
- `coalesce_events_test.go`
- `snippet.json`