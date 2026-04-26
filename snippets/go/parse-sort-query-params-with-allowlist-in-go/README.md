# Parse Sort Query Params with an Allowlist in Go

Parse a comma-separated sort query parameter against a fixed allowlist of field names in Go.

This snippet is useful when an API should accept one ordered sort expression while rejecting unknown or repeated fields before the query reaches storage code.

## Highlights

- Parses asc and desc terms
- Rejects unknown sort fields
- Prevents duplicate sort keys

## What It Does

- splits a comma-separated sort string into ordered tokens
- supports descending fields with a leading `-`
- validates each field against a fixed allowlist
- rejects blank terms and repeated fields
- preserves the caller's sort order in the parsed result

## Usage

```go
// Run directly:
// go run parse_sort_query_params.go
// The example parses one descending and one ascending sort term.
```

## Notes

- A missing or blank sort string returns `nil` so callers can apply their own default sort.
- The allowlist is trimmed before validation, and blank allowlist entries are rejected.
- Duplicate sort fields are treated as an error even if the directions differ.

## Verification

Run the tests from the snippet root:

```bash
go test -race ./...
```

The verified test suite covers:

- one ascending field
- multiple mixed-direction fields
- a missing sort param
- unknown fields
- empty terms
- repeated fields

## Files

- `parse_sort_query_params.go`
- `parse_sort_query_params_test.go`
- `snippet.json`