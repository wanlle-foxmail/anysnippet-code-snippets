# Return Problem JSON Errors with Echo

Return `application/problem+json` errors from Echo with one small error handler and one reusable writer helper.

This snippet is useful when an API should return one consistent RFC 7807 style error shape instead of ad hoc JSON error bodies.

## Highlights

- Writes problem+json responses
- Handles Echo and generic errors
- Preserves committed responses

## What It Does

- Builds one Echo `HTTPErrorHandler` for `application/problem+json`
- Renders Echo `HTTPError` values with status, title, detail, and instance
- Renders generic errors as `500 Internal Server Error`
- Avoids leaking raw generic error text to clients
- Leaves already committed responses untouched

## Usage

```go
// Run directly:
// go run problem_json_error_handler.go
// Then call GET /orders/123 or GET /boom.
```

Example response:

```json
{
  "type": "about:blank",
  "title": "Not Found",
  "status": 404,
  "detail": "order not found",
  "instance": "/orders/123"
}
```

## Notes

- This snippet keeps the core RFC 7807 fields only.
- Generic server errors return a safe public detail instead of the raw internal error string.
- `instance` uses the request path so callers can match the error to the request URL.

## Verification

Run the tests from the snippet root:

```bash
go test -race ./...
```

The verified test suite covers:

- formatting an Echo `HTTPError`
- formatting a generic `500` error
- falling back to status text for non-string error messages
- preserving already committed responses
- filling default `type` and `title` values
- rejecting a nil Echo context in the writer helper

## Files

- `problem_json_error_handler.go`
- `problem_json_error_handler_test.go`
- `snippet.json`