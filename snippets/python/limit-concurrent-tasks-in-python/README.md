# Limit I/O Concurrency in Python

Run I/O-bound tasks with ThreadPoolExecutor and a fixed worker limit.

This snippet is useful when you need to process many I/O-bound tasks without overwhelming an API, a file system, or a remote service.

## Highlights

- Caps thread concurrency
- Keeps result order
- Uses Python stdlib only

## Use Cases

- Limit concurrent API calls during bulk imports
- Process many files without spawning unbounded threads
- Run batch jobs where partial failures should be reported instead of stopping all work

## Code

```python
from src.limit_concurrent_tasks import limit_concurrent_tasks


def fetch_user(user_id):
    return {"user_id": user_id, "status": "ok"}


result = limit_concurrent_tasks([101, 102, 103], fetch_user, max_workers=3)
print(result)
```

## Notes

- `ThreadPoolExecutor.map` keeps the output in the same order as the input items.
- Worker exceptions are raised back to the caller.
- This snippet targets I/O-bound work; use a process-based approach for CPU-bound tasks.

## Verification

Run the unit tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- successful batch execution
- stable result ordering despite out-of-order completion
- worker error propagation
- invalid `max_workers` values
- empty input handling
- concurrency limit enforcement

## Files

- `src/limit_concurrent_tasks.py`
- `tests/test_limit_concurrent_tasks.py`
- `snippet.json`