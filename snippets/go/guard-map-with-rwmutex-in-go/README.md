# Guard a Map with RWMutex in Go

Wrap a Go map with `sync.RWMutex` so concurrent reads and writes stay race-free.

This snippet is useful when a plain map is shared across goroutines and you want one small wrapper instead of scattering locks throughout the call sites.

## Highlights

- Makes the zero value usable
- Uses read locks for lookups
- Returns snapshot copies safely

## What It Does

- stores map data behind one `sync.RWMutex`
- supports `Set`, `Get`, `Delete`, and `Len`
- lazily initializes the internal map so the zero value works
- returns a copied snapshot for iteration or logging
- keeps map access race-free under concurrent reads and writes

## Usage

```go
// Run directly:
// go run safe_map.go
// The example stores one counter value and reads it back.
```

## Notes

- `Snapshot` returns a copy so callers can iterate without holding the lock.
- The zero value is ready to use, but `NewSafeMap` is available when you prefer an explicit constructor.
- Run the tests with `-race` to verify the concurrency contract.

## Verification

Run the tests from the snippet root:

```bash
go test -race ./...
```

The verified test suite covers:

- using the zero value directly
- missing-key lookups
- deleting keys
- tracking map length
- returning independent snapshots
- concurrent writes
- concurrent reads during writes

## Files

- `safe_map.go`
- `safe_map_test.go`
- `snippet.json`