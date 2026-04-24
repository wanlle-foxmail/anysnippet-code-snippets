# Stream Server-Sent Events with Go

Stream `text/event-stream` responses from Go with one message channel and per-message flushes.

This snippet is useful when one handler should push a small stream of text updates to the browser without switching to WebSockets.

## Highlights

- Writes `text/event-stream` headers
- Flushes after each message
- Stops on disconnect or close

## What It Does

- Validates the response writer, request, and message channel
- Sets the standard SSE response headers
- Writes each message as one `data: ...` event block
- Splits multi-line payloads into valid repeated `data:` lines
- Flushes after every message so the client sees updates immediately
- Stops cleanly when the message channel closes or the client disconnects

## Usage

```go
// Run directly:
// go run stream_sse.go
// Then open http://localhost:8080/events in a browser or curl client.
```

## Notes

- This snippet streams plain `data:` events and does not add custom event names or retry hints.
- Messages should already be formatted as the payload you want the browser to receive.
- The handler should keep the channel alive until it is done sending events.

## Verification

Run the tests from the snippet root:

```bash
go test -race ./...
```

The verified test suite covers:

- SSE headers and message formatting
- one flush per message
- stopping on client disconnect
- stopping on channel close
- missing flusher rejection
- write error propagation
- nil input validation

## Files

- `stream_sse.go`
- `stream_sse_test.go`
- `snippet.json`