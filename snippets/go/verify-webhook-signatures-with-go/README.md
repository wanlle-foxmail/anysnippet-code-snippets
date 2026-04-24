# Verify Webhook Signatures with Go

Verify a webhook body against a `sha256=` HMAC signature header in Go.

This snippet is useful when one inbound webhook should be accepted only if the sender signed the raw request body with a shared secret.

## Highlights

- Parses `sha256=` signature headers
- Uses constant-time comparison
- Accepts upper or lower hex

## What It Does

- Reads one `sha256=<hex>` signature header value
- Decodes the provided hex digest
- Computes an `HMAC-SHA256` digest for the raw request body
- Uses `hmac.Equal` for constant-time comparison
- Returns `false` for malformed headers instead of panicking

## Usage

```go
// Run directly:
// go run verify_webhook_signature.go
// The example prints true for one matching body and signature.
```

## Notes

- The body must be the exact raw bytes that the sender signed.
- This snippet intentionally validates one `sha256=` header format and one shared secret.
- Empty secrets are technically valid HMAC keys, even though production webhooks should use a real secret.

## Verification

Run the tests from the snippet root:

```bash
go test -race ./...
```

The verified test suite covers:

- matching signatures
- mismatched signatures
- empty bodies
- empty headers
- empty secrets
- uppercase hex signatures
- malformed hex values
- wrong header prefixes

## Files

- `verify_webhook_signature.go`
- `verify_webhook_signature_test.go`
- `snippet.json`