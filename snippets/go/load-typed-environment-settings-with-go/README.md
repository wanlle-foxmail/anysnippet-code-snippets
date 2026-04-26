# Load Typed Environment Settings with Go

Load `APP_ENV`, `PORT`, and `DEBUG` from environment variables with string, integer, and boolean parsing in Go.

This snippet is useful when a small Go service should fail fast on missing or invalid environment settings without pulling in a larger config package.

- Parses text, int, and bool values
- Uses defaults for optional settings
- Reads from a provided map or the OS environment

## Example

```go
// go run load_typed_env_settings.go
```

The example `main` function logs one parsed `AppSettings` value.

## Notes

- `APP_ENV` is required.
- `PORT` defaults to `8000` and must stay between `1` and `65535`.
- `DEBUG` defaults to `false` and accepts `true`, `false`, `1`, or `0`.

## Verification

Run the tests from the snippet root:

```bash
go test -race ./...
```

Verified behavior covers:

- loading required and optional values
- using defaults for optional values
- rejecting a missing required value
- rejecting an invalid integer value
- rejecting out-of-range port values
- rejecting an invalid boolean value
- reading from the OS environment when the input map is nil

## Files

- `load_typed_env_settings.go`
- `load_typed_env_settings_test.go`