# Graceful HTTP Server Shutdown with Go

Run an HTTP server until a context is canceled, then shut it down with a timeout.

This snippet is useful when a Go service should stop accepting new requests, wait for in-flight work to finish, and still fail fast if shutdown takes too long.

## Highlights

- Shuts down on context cancel
- Waits for in-flight requests
- Returns timeout errors clearly

## What It Does

- Starts an `http.Server` on a provided listener
- Keeps serving until the context is canceled
- Calls `Shutdown` with a fixed timeout window
- Ignores the expected `http.ErrServerClosed` shutdown path
- Returns real serve or shutdown errors to the caller

## Usage

```go
// Run directly:
// go run graceful_http_server_shutdown.go
// Then visit http://localhost:8080/hello
// Press Ctrl+C to stop the server gracefully.
```

## Notes

- `ServeWithGracefulShutdown` returns an error when the server or listener is missing.
- A non-positive shutdown timeout is rejected.
- If an in-flight request does not finish before the timeout, the function returns the shutdown error.
- The example `main` uses `os.Interrupt` so it stays portable across common operating systems.

## Verification

Run the tests from the snippet root:

```bash
go test -race ./...
```

The verified test suite covers:

- serving requests before shutdown
- waiting for an in-flight request to finish
- returning a timeout error when shutdown takes too long
- surfacing serve errors from a bad listener
- rejecting a nil server
- rejecting a nil listener or bad timeout

## Files

- `graceful_http_server_shutdown.go`
- `graceful_http_server_shutdown_test.go`
- `snippet.json`