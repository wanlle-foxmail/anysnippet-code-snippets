# Retry Failed Operations with Python

Retry one callable on retryable errors with a fixed delay between attempts.

This snippet is useful when one small operation can fail transiently and you want a short, explicit retry loop without bringing in a larger retry library.

## Highlights

- Retries only selected errors
- Uses a fixed retry delay
- Re-raises the last failure

## Use Cases

- Retry short network probes
- Retry a temporary file or lock check
- Wrap one flaky dependency call in a small script

## Code

```python
from src.retry_failed_operation import retry_failed_operation


attempts = {"count": 0}


def sometimes_ready() -> str:
    attempts["count"] += 1
    if attempts["count"] < 3:
        raise TimeoutError("not ready yet")
    return "ready"


print(retry_failed_operation(sometimes_ready, max_attempts=3, delay_seconds=0.1, retry_exceptions=(TimeoutError,)))
```

## Notes

- Only exceptions listed in `retry_exceptions` are retried.
- Non-retryable exceptions fail immediately.
- The last retryable exception is re-raised when attempts run out.
- This snippet uses a fixed delay on purpose and does not include backoff or jitter.

## Verification

Run the unit tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- first-attempt success
- retry before eventual success
- re-raising the last retryable failure
- non-retryable errors that fail immediately
- sleep behavior between retry attempts
- invalid attempt counts
- invalid delay values

## Files

- `src/retry_failed_operation.py`
- `tests/test_retry_failed_operation.py`
- `snippet.json`