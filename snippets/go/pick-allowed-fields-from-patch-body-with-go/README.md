# Pick Allowed Fields from a PATCH Body in Go

Keep only allowed top-level fields from a PATCH-style request body in Go.

This snippet is useful when partial updates should ignore unknown top-level keys before the payload reaches business logic or persistence code.

## Highlights

- Keeps only allowed keys
- Preserves explicit null values
- Ignores unknown patch fields

## What It Does

- validates the allowlist before filtering
- accepts a decoded `map[string]interface{}` PATCH body
- keeps only allowed top-level keys
- preserves explicit `null`, `false`, and `0` values
- returns an empty map for a nil body

## Usage

```go
// Run directly:
// go run pick_allowed_patch_fields.go
// The example keeps only display_name and nickname from one PATCH body.
```

## Notes

- The helper does not mutate the input map.
- Unknown fields are ignored instead of treated as an error.
- The snippet is intentionally limited to top-level keys.

## Verification

Run the tests from the snippet root:

```bash
go test -race ./...
```

The verified test suite covers:

- selecting allowed fields
- preserving explicit values
- ignoring disallowed fields
- nil bodies
- non-map inputs
- blank allowed field names

## Files

- `pick_allowed_patch_fields.go`
- `pick_allowed_patch_fields_test.go`
- `snippet.json`