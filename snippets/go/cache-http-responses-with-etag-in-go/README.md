# Cache HTTP Responses with ETag in Go

Cache one JSON HTTP response in Go with an `ETag` header and `304 Not Modified` handling.

This snippet is useful when a small read endpoint should let clients skip downloading the same JSON body again when the response has not changed.

## Highlights

- Writes JSON with ETag
- Returns 304 on cache hits
- Uses stdlib only

## What It Does

- Marshals one JSON response body
- Hashes that body into one strong `ETag` value
- Returns `304 Not Modified` when `If-None-Match` matches
- Writes the JSON body normally when the tag does not match
- Supports exact and wildcard `If-None-Match` checks

## Usage

```go
// Run directly:
// go run etag_response.go
// Then call GET /products/42 with or without If-None-Match.
```

## Notes

- This snippet hashes the rendered JSON body, so response changes automatically produce a new `ETag`.
- The helper focuses on one JSON response at a time and does not add a separate cache store.
- `If-None-Match` matching stays intentionally small: exact tag matches, wildcard matches, and comma-separated validator lists.

## Verification

Run the tests from the snippet root:

```bash
go test ./...
```

The verified test suite covers:

- writing JSON bodies with an `ETag`
- returning `304` when the tag matches
- ignoring non-matching validators
- matching wildcard validators
- matching one tag in a validator list
- rejecting nil inputs
- propagating write errors

## Files

- `etag_response.go`
- `etag_response_test.go`
- `snippet.json`