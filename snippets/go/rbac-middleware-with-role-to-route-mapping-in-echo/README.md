# RBAC Middleware with Role-to-Route Mapping in Echo

Protect Echo routes with a fixed role-to-route mapping.

This snippet is useful when authentication already happened upstream and an Echo app should allow or block routes based on one small role mapping.

- Maps routes to allowed roles
- Rejects missing roles with `401`
- Rejects forbidden roles or unmapped routes with `403`

## Example

```go
// go run rbac_middleware.go
```

Then call `GET /admin` or `GET /reports` with an `X-User-Role` header.

## Notes

- This snippet keeps role extraction simple by reading `X-User-Role`.
- Use it after authentication if another middleware or proxy already resolved the caller role.
- The middleware denies routes that are not listed in the mapping.

## Verification

Run the tests from the snippet root:

```bash
go test ./...
```

Verified behavior covers:

- an allowed role on a mapped route
- a missing role header
- a forbidden role
- a shared route with multiple allowed roles
- an unmapped route
- whitespace and case normalization for the role header

## Files

- `rbac_middleware.go`
- `rbac_middleware_test.go`