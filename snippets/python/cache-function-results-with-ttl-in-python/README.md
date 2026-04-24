# Cache Function Results with TTL in Python

Cache successful function results in memory until a small TTL expires.

This snippet is useful when a function is cheap to cache in-process and you want to avoid repeated work for the same arguments during a short time window.

## Highlights

- Caches results in memory
- Reuses results within the TTL
- Skips caching failed calls

## Use Cases

- Avoid repeating short-lived API lookups
- Reuse small derived values in one process
- Add a small cache without extra dependencies

## Code

```python
from src.cache_function_results_with_ttl import cache_function_results_with_ttl


@cache_function_results_with_ttl(30.0)
def load_status(name: str) -> str:
    return f"ready:{name}"


print(load_status("worker-a"))
```

## Notes

- Cache entries are stored only in the current process.
- The cache grows with unique argument combinations and does not evict old keys automatically.
- The cache uses a plain dictionary and is not thread-safe.
- Function arguments must be hashable.
- Keyword arguments are normalized, so different keyword order still hits the same cache entry.
- Exceptions are never cached.

## Verification

Run the unit tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- cache reuse within the TTL
- cache refresh after expiration
- separate cache entries by arguments
- keyword argument order normalization
- exception paths that do not cache failures
- invalid TTL input

## Files

- `src/cache_function_results_with_ttl.py`
- `tests/test_cache_function_results_with_ttl.py`
- `snippet.json`