# Reject Unknown JSON Fields with Go

Decode one JSON body in Go while rejecting fields that are not declared on the target struct.

This snippet is useful when an API should fail fast on unexpected input keys instead of silently ignoring them.

- Rejects unknown JSON fields
- Rejects malformed JSON input
- Rejects multiple JSON values in one body

## Example

```go
// go run reject_unknown_json_fields.go
```

The example `main` function decodes one valid request body into `CreateUserRequest`.

## Notes

- The helper uses `json.Decoder.DisallowUnknownFields()`.
- It expects exactly one JSON value in the request body.
- This snippet handles strict decoding only; field validation should stay in a separate step.

## Verification

Run the tests from the snippet root:

```bash
go test ./...
```

Verified behavior covers:

- valid known fields
- unknown fields
- malformed JSON
- multiple JSON values
- nil readers
- nil targets

## Files

- `reject_unknown_json_fields.go`
- `reject_unknown_json_fields_test.go`