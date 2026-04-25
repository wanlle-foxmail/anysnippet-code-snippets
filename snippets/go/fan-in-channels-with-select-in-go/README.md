# Fan-In Channels with select in Go

Combine two Go channels into one output stream with `select` and optional context cancellation.

This snippet is useful when two producers feed one consumer and you want one small helper that forwards values until both inputs are done.

## Highlights

- Combines two input channels
- Stops cleanly on context cancel
- Closes output after both inputs

## What It Does

- reads from two input channels with one `select`
- forwards values into one output channel
- stops when both inputs are closed
- stops early when the provided context is canceled
- accepts nil input channels and ignores them

## Usage

```go
// Run directly:
// go run fan_in_channels.go
// The example forwards values from two small string channels.
```

## Notes

- Output order is intentionally not guaranteed across the two inputs.
- The helper keeps the fan-in surface focused on exactly two channels.
- Passing `nil` for one side lets you reuse the helper with one active input.

## Verification

Run the tests from the snippet root:

```bash
go test -race ./...
```

The verified test suite covers:

- merging values from both inputs
- closing the output after both inputs close
- draining only the left input
- draining only the right input
- returning a closed output for two nil inputs
- stopping when the context is canceled

## Files

- `fan_in_channels.go`
- `fan_in_channels_test.go`
- `snippet.json`