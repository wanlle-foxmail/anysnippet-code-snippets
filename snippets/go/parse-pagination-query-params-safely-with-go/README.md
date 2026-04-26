# Parse Pagination Query Params Safely in Go

Parse `page` and `page_size` query parameters with defaults, validation, and a fixed maximum in Go.

This snippet is useful when one HTTP handler or service helper should turn raw pagination query params into a ready-to-use page, page size, and offset.

## Highlights

- Uses safe pagination defaults
- Caps large page sizes
- Returns a ready offset value

## What It Does

- reads `page` and `page_size` from `url.Values`
- applies defaults when values are missing or blank
- rejects invalid integers and non-positive numbers
- caps oversized page sizes at a fixed maximum
- returns a computed offset alongside the parsed values

## Usage

```go
// Run directly:
// go run parse_pagination_query_params.go
// The example parses one page and page_size pair and logs the result.
```

## Notes

- The helper trims surrounding whitespace before parsing integers.
- Offsets use the final capped page size.
- Invalid defaults or max-size settings are rejected early.

## Verification

Run the tests from the snippet root:

```bash
go test -race ./...
```

The verified test suite covers:

- valid page and page_size parsing
- missing values with defaults
- page size capping
- invalid integers
- non-positive values
- invalid default pages
- invalid default page sizes
- invalid max page sizes
- surrounding whitespace

## Files

- `parse_pagination_query_params.go`
- `parse_pagination_query_params_test.go`
- `snippet.json`