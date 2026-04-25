# Deduplicate In-Flight Work with singleflight in Go

Deduplicate concurrent Go loads for the same key with `singleflight.Group` while still rerunning work after the first request finishes.

This snippet is useful when multiple goroutines ask for the same expensive value at the same time and you want only one backend call in flight.

## Highlights

- Collapses duplicate concurrent loads
- Keeps different keys independent
- Does not replay completed results

## What It Does

- wraps one string-loading function with `singleflight.Group`
- shares one in-flight result across concurrent callers for the same key
- reports whether the result was shared
- reruns the load after the first in-flight call completes
- propagates one shared error to all waiting callers

## Usage

```go
// Run directly:
// go run singleflight_loader.go
// The example starts two concurrent loads for the same key.
```

## Notes

- `singleflight` only suppresses duplicate in-flight calls. It does not cache finished results.
- The returned `shared` flag tells you whether another caller received the same in-flight result.
- This snippet keeps the value type concrete as `string` to keep the pattern easy to scan.

## Verification

Run the tests from the snippet root:

```bash
go test -race ./...
```

The verified test suite covers:

- sharing one concurrent load for the same key
- keeping different keys independent
- rerunning the load after the first call completes
- sharing concurrent errors
- rejecting an empty key
- rejecting a nil load function

## Files

- `singleflight_loader.go`
- `singleflight_loader_test.go`
- `snippet.json`