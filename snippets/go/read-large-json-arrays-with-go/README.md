# Read Large JSON Arrays with Go

Stream a top-level JSON array in Go and decode one item at a time.

This snippet is useful when a large file or request body should be processed incrementally instead of loaded fully into memory.

## Highlights

- Streams one item at a time
- Avoids full-array unmarshalling
- Keeps nested JSON values intact

## What It Does

- opens a top-level JSON array with `json.Decoder`
- decodes one item at a time into the target type
- passes each decoded item to a handler callback
- rejects malformed arrays and non-array top-level values
- enforces exactly one top-level JSON array per input

## Usage

```go
// Run directly:
// go run read_large_json_array.go
// The example streams two Event values from one JSON array.
```

## Notes

- This pattern helps when arrays are too large for one `json.Unmarshal` call.
- The handler can stop processing early by returning an error.
- The snippet intentionally focuses on top-level arrays, not newline-delimited JSON.

## Verification

Run the tests from the snippet root:

```bash
go test -race ./...
```

The verified test suite covers:

- ordered item reads
- nested objects and lists
- empty arrays
- malformed arrays
- non-array inputs
- trailing content after one array
- handler errors
- nil readers
- nil handlers

## Files

- `read_large_json_array.go`
- `read_large_json_array_test.go`
- `snippet.json`