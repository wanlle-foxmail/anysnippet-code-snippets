# Verify JWT Bearer Tokens with FastAPI

Protect one FastAPI route by validating HS256 bearer tokens and exposing claims to the handler.

This snippet is useful when one API route should require a short-lived JWT without pulling in a larger authentication framework.

## Highlights

- Validates HS256 bearer tokens
- Rejects bad or expired tokens
- Exposes claims to handlers

## What It Does

- Builds one FastAPI dependency for `Authorization: Bearer ...`
- Accepts only `HS256` signed tokens
- Requires an `exp` claim and rejects expired tokens
- Shows one protected `GET /profile` route
- Returns a small profile payload from validated claims

## Usage

```python
from src.jwt_bearer_auth import app

# Run directly:
# python src/jwt_bearer_auth.py
# Then send GET /profile with a Bearer token.
```

Success response example:

```json
{
  "sub": "demo-user",
  "role": "admin"
}
```

## Notes

- This snippet intentionally validates one shared secret and one signing method.
- `Authorization` must use the bearer scheme.
- Later handlers can reuse the validated claims dependency or call the helper directly.
- Replace the demo secret with an environment variable or secret manager before using this pattern in production.

## Verification

Run the unit tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- accepting a valid token
- rejecting a missing authorization header
- rejecting malformed authorization headers
- rejecting blank bearer tokens
- rejecting malformed JWT strings
- rejecting invalid signatures
- rejecting expired tokens
- rejecting missing `exp` claims
- rejecting unexpected signing methods

## Files

- `src/jwt_bearer_auth.py`
- `tests/test_jwt_bearer_auth.py`
- `snippet.json`