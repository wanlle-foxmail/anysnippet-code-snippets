# Limit JSON Request Body Size in Go

Decode one JSON request body in Go while enforcing a fixed maximum size.

This snippet is useful when an HTTP handler should reject oversized JSON bodies before they consume more memory or CPU than expected.

## Highlights

- Enforces one fixed body limit
- Uses http.MaxBytesReader directly
- Accepts exactly one JSON value

## What It Does

- wraps the request body with `http.MaxBytesReader`
- decodes one JSON value into the target struct
- rejects oversized bodies with the standard `http.MaxBytesError`
- rejects malformed JSON and extra JSON values
- validates required inputs before decoding starts

## Usage

```go
// Run directly:
// go run limit_json_body_size.go
// The example decodes one small JSON body under a fixed byte limit.
```

## Notes

- This snippet teaches request-size limiting only and stays separate from unknown-field validation.
- The caller can inspect `*http.MaxBytesError` with `errors.As` when a body is too large.
- The helper requires an `http.ResponseWriter` because `http.MaxBytesReader` is request-aware.

## Verification

Run the tests from the snippet root:

```bash
go test -race ./...
```

The verified test suite covers:

- accepting a body under the limit
- accepting a body at the exact limit
- rejecting oversized bodies
- rejecting malformed JSON under the limit
- rejecting multiple JSON values
- rejecting invalid input

## Files

- `limit_json_body_size.go`
- `limit_json_body_size_test.go`
- `snippet.json`