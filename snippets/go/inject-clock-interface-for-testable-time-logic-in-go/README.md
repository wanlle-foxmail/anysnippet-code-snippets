# Inject a Clock Interface for Testable Time Logic in Go

Inject a tiny clock interface into Go time checks so tests can control the current time without sleeping.

This snippet is useful when production code depends on `time.Now()` but your tests should stay deterministic and fast.

## Highlights

- Injects time through one interface
- Removes sleeps from time-based tests
- Uses the real clock in production

## What It Does

- defines a small `Clock` interface with one `Now()` method
- uses `RealClock` in production code
- accepts any clock implementation in `HasExpired`
- lets tests freeze time with a simple fake clock

## Usage

```go
// Run directly:
// go run clock_interface.go
// The example checks whether one expiration time has already passed.
```

## Notes

- `HasExpired` treats the exact expiration instant as expired.
- The snippet keeps the interface deliberately small so you can fake it in tests.
- The fake clock in the test file is one way to control time without waiting.

## Verification

Run the tests from the snippet root:

```bash
go test -race ./...
```

The verified test suite covers:

- checking a time before expiration
- checking a time after expiration
- the exact expiration boundary
- rejecting a nil clock
- rejecting a zero expiration time
- confirming that `RealClock` reads from the real clock window

## Files

- `clock_interface.go`
- `clock_interface_test.go`
- `snippet.json`