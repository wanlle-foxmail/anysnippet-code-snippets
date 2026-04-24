# Verify JWT Bearer Tokens with Echo

Protect one Echo route by validating HS256 bearer tokens and exposing claims to the handler.

This snippet is useful when one API route should require a short-lived JWT without pulling in a larger authentication framework.

## Highlights

- Validates HS256 bearer tokens
- Rejects bad or expired tokens
- Exposes claims to handlers

## What It Does

- Builds one Echo middleware for `Authorization: Bearer ...`
- Accepts only `HS256` signed tokens
- Requires an `exp` claim and rejects expired tokens
- Stores validated claims in `echo.Context` under `JWTClaimsContextKey`
- Shows one protected `GET /profile` route

## Usage

```go
// Run directly:
// go run jwt_bearer_auth.go
// The example prints one sample bearer token for GET /profile.
```

Success response example:

```json
{
  "sub": "demo-user",
  "role": "admin"
}
```

## Notes

- This snippet intentionally validates only one shared secret and one signing method.
- `Authorization` must use the bearer scheme.
- Later handlers can read the claims with `c.Get(JWTClaimsContextKey)`.
- The example keeps claims as `jwt.MapClaims` on purpose instead of adding a custom claim struct.
- Replace the demo secret with an environment variable or secret manager before using this pattern in production.

## Verification

Run the tests from the snippet root:

```bash
go mod tidy
go test -race ./...
```

The verified test suite covers:

- accepting a valid token
- rejecting missing and malformed authorization headers
- rejecting malformed JWT strings
- rejecting tokens with invalid signatures
- rejecting expired tokens
- rejecting tokens without an `exp` claim
- rejecting unexpected signing methods

## Files

- `jwt_bearer_auth.go`
- `jwt_bearer_auth_test.go`
- `snippet.json`