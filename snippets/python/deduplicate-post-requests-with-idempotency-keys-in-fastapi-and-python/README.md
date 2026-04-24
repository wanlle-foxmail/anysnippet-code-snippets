# Deduplicate POST Requests with Idempotency Keys in FastAPI

Deduplicate one FastAPI `POST` endpoint with an `Idempotency-Key` header and a single-process in-memory store.

This snippet is useful when one FastAPI write endpoint should replay successful responses for repeated keys and reject duplicates that are still in progress.

- Replays successful `POST` responses
- Returns `409` while the same key is still in progress
- Scopes keys by path, caller, and idempotency key

## Example

```python
from src.idempotency_middleware import app

# python src/idempotency_middleware.py
```

Then call `POST /orders` with an `Idempotency-Key` header.

## Notes

- Successful `2xx` responses are cached and replayed.
- Failed responses are not cached, so the same key can be retried.
- Caller scope prefers `X-User-ID`, then a normalized `Authorization` header, then anonymous scope.

## Verification

Run the tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

Verified behavior covers:

- requiring `Idempotency-Key` on `POST`
- replaying successful responses
- separating different keys
- returning `409` while work is in progress
- not caching failed responses
- bypassing non-`POST` requests
- separating user scopes
- normalizing authorization scope

## Files

- `src/idempotency_middleware.py`
- `tests/test_idempotency_middleware.py`