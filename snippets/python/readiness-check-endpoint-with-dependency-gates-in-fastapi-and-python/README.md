# Readiness Check Endpoint with Dependency Gates in FastAPI

Gate a FastAPI readiness endpoint on required dependencies and return `503` until the service is ready.

This snippet is useful when a service should expose a `/ready` endpoint that blocks traffic until required startup or dependency checks pass.

- Separates required and optional gates
- Returns `200` only when required gates are ready
- Reports optional failures without blocking readiness

## Example

```python
from src.readiness_check import app

# python src/readiness_check.py
```

Then visit `http://localhost:8000/ready`.

## Notes

- Required gates control whether the endpoint returns `200` or `503`.
- Optional gates are still reported in the response payload.
- This snippet focuses on readiness, not liveness or aggregated health reporting.

## Verification

Run the tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

Verified behavior covers:

- all required gates ready
- database gate failure
- migration gate failure
- optional cache failure without blocking readiness
- multiple required failures
- required and optional gate markers in the response

## Files

- `src/readiness_check.py`
- `tests/test_readiness_check.py`