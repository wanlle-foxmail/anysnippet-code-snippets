# Stream JSONL HTTP Responses with Go

Stream JSONL HTTP responses from Go with one record channel and per-record flushes.

This snippet is useful when a handler should send many JSON records incrementally instead of waiting for one large array response to finish.

## Highlights

- Streams one JSONL record per line
- Flushes after each record
- Uses stdlib only

## What It Does

- Validates the writer, request, and record channel
- Sets `application/x-ndjson` response headers
- Encodes each record as one JSON line
- Flushes after every record so clients can process data incrementally
- Stops cleanly when the channel closes or the client disconnects

## Usage

```go
// Run directly:
// go run jsonl_stream.go
// Then call GET /records.
```

## Notes

- This snippet expects each record to be JSON serializable.
- `application/x-ndjson` is a practical content type for JSONL over HTTP.
- The handler streams one process-local channel and does not add buffering or replay behavior.

## Verification

Run the tests from the snippet root:

```bash
go test -race ./...
```

The verified test suite covers:

- writing JSONL headers and records
- flushing after each record
- stopping on client disconnect
- stopping when the channel closes
- rejecting writers without streaming support
- propagating marshal errors
- rejecting nil inputs and propagating write errors

## Files

- `jsonl_stream.go`
- `jsonl_stream_test.go`
- `snippet.json`