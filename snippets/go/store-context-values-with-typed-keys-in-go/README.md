# Store Context Values with Typed Keys in Go

Store and read `context.Context` values in Go with typed keys instead of raw strings.

This snippet is useful when middleware and handlers should share request-scoped values without risking accidental key collisions.

## Highlights

- Uses typed context keys
- Avoids string key collisions
- Works with nil parent contexts

## What It Does

- creates unique typed keys with `NewContextKey`
- attaches one typed value with `context.WithValue`
- reads the value back with the expected type
- treats missing values as a clean `ok=false` result
- keeps distinct keys separate even when names match

## Usage

```go
// Run directly:
// go run store_context_values.go
// The example stores one request_id and reads it back from the context.
```

## Notes

- The key name is only for human-readable intent; pointer identity keeps keys distinct.
- A nil parent context falls back to `context.Background()` for convenience.
- This snippet is intended for small request-scoped values, not optional parameter bags.

## Verification

Run the tests from the snippet root:

```bash
go test -race ./...
```

The verified test suite covers:

- one typed value round-trip
- missing values
- nil keys
- nil parent contexts
- distinct keys with the same name
- blank key names

## Files

- `store_context_values.go`
- `store_context_values_test.go`
- `snippet.json`