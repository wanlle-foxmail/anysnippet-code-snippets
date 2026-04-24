# Validate Multipart Uploads with Echo

Validate one multipart upload in Echo with filename, content type, and file size checks.

This snippet is useful when an Echo route should accept one uploaded file only after validating its name, MIME type, and body size.

- Checks filename and MIME type
- Rejects empty or oversized files
- Returns accepted file metadata

## Example

```go
// go run upload_validation_handler.go
```

Then send one multipart `POST` request to `/upload` with a `file` part.

## Notes

- Allowed content types stay explicit: `text/plain` and `text/csv`.
- The handler sanitizes the returned filename down to its basename.
- File size is measured by reading the uploaded body, not by trusting metadata alone.

## Verification

Run the tests from the snippet root:

```bash
go test ./...
```

Verified behavior covers:

- valid uploads
- blank filenames
- unsupported content types
- empty files
- oversized files
- exact-limit files
- sanitized basenames
- the explicit content-type allowlist

## Files

- `upload_validation_handler.go`
- `upload_validation_handler_test.go`