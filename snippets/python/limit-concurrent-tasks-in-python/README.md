# Limit Concurrent Tasks with Python

Run a batch of tasks with a fixed concurrency limit and return ordered success and failure results.

This snippet is useful when you need to process many I/O-bound tasks without overwhelming an API, a file system, or a remote service.

## Highlights

- Uses `ThreadPoolExecutor` to cap concurrent work
- Preserves input order in the returned results
- Captures task failures without stopping the whole batch
- Returns summary counts together with per-task outcomes

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
print(result["succeeded_count"])
print(result["results"])
```

## Notes

- Results keep the same order as the input items even if tasks finish out of order.
- Individual task exceptions are captured per result instead of aborting the full batch.
- This snippet targets I/O-bound work; use a process-based approach for CPU-bound tasks.

## Verification

Run the unit tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- successful batch execution and summary counts
- stable result ordering despite out-of-order completion
- per-task error capture without aborting the whole batch
- invalid `max_workers` values
- empty input handling
- concurrency limit enforcement
- generator input support

## Files

- `src/limit_concurrent_tasks.py`
- `tests/test_limit_concurrent_tasks.py`
- `snippet.json`