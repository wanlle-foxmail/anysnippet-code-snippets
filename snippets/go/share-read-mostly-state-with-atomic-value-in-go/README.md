# Share Read-Mostly State with atomic.Value in Go

Share read-mostly state in Go by swapping whole immutable snapshots with `atomic.Value`.

This snippet is useful when many goroutines need lock-free reads of the latest shared configuration or routing table.

## Highlights

- Reads without read locks
- Swaps whole state snapshots
- Works cleanly under race tests

## What It Does

- stores one initial snapshot in an `atomic.Value`
- loads the current snapshot without taking a mutex
- swaps in a complete replacement snapshot on writes
- returns explicit errors for nil receivers
- stays focused on immutable read-mostly data

## Usage

```go
// Run directly:
// go run atomic_value_state.go
// The example swaps one routing snapshot and loads the latest version.
```

## Notes

- Treat stored snapshots as immutable after `Store` returns.
- `atomic.Value` is a good fit when writes are infrequent and readers only need the latest whole snapshot.
- This snippet uses concrete struct snapshots to avoid nil stores and partial mutation.

## Verification

Run the tests from the snippet root:

```bash
go test -race ./...
```

The verified test suite covers:

- initial snapshot reads
- updated snapshot writes
- concurrent loads and stores
- nil receivers
- zero-value snapshots

## Files

- `atomic_value_state.go`
- `atomic_value_state_test.go`
- `snippet.json`