# Acquire Redis Locks with Go

Acquire and release a Redis lock in Go with `SET NX` and a compare-and-delete release step.

This snippet is useful when one worker should claim a short-lived lock and release it only if it still owns the lock token.

## Highlights

- Uses `SET NX` for lock acquire
- Releases only matching owners
- Keeps Redis calls testable

## What It Does

- Acquires one lock key with `SET NX` and a TTL
- Uses one owner token string to identify the lock holder
- Releases the lock with a compare-and-delete Lua script
- Wraps the real `go-redis` client behind a tiny store interface for tests
- Returns `false` instead of an error when another worker already holds the lock

## Usage

```go
// Run directly:
// go run redis_lock.go
// The example expects a Redis server on localhost:6379.
```

## Notes

- This snippet intentionally targets one Redis instance, not Redlock or multi-node coordination.
- The caller should generate a unique owner token per worker or per critical section.
- The unit tests verify the lock flow with a fake store, not a live Redis server.
- If the lock holder crashes, the lock stays held until the TTL expires; this snippet does not renew leases automatically.

## Verification

Run the tests from the snippet root:

```bash
go mod tidy
go test -race ./...
```

The verified test suite covers:

- successful lock acquisition
- lock contention
- wrapped store errors
- releasing with a matching owner token
- refusing release for a different owner token
- invalid input handling

## Files

- `redis_lock.go`
- `redis_lock_test.go`
- `snippet.json`