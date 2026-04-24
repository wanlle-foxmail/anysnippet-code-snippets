import time
from functools import wraps
from typing import Any, Callable, Dict, Tuple


CacheKey = Tuple[Tuple[Any, ...], Tuple[Tuple[str, Any], ...]]
CacheEntry = Tuple[float, Any]


def cache_function_results_with_ttl(ttl_seconds: float) -> Callable[[Callable[..., Any]], Callable[..., Any]]:
    """Cache successful function results in memory until the TTL expires."""
    if isinstance(ttl_seconds, bool):
        raise TypeError("ttl_seconds must be a number")
    if not isinstance(ttl_seconds, (int, float)):
        raise TypeError("ttl_seconds must be a number")
    if ttl_seconds <= 0:
        raise ValueError("ttl_seconds must be greater than 0")

    def decorator(function: Callable[..., Any]) -> Callable[..., Any]:
        cache: Dict[CacheKey, CacheEntry] = {}

        @wraps(function)
        def wrapper(*args: Any, **kwargs: Any) -> Any:
            # Flow:
            #   build cache key
            #      |
            #      +-> fresh hit -> return cached result
            #      `-> expired or missing key -> call function -> store result -> return it
            key = (args, tuple(sorted(kwargs.items())))
            now = time.monotonic()
            cached_entry = cache.get(key)

            if cached_entry is not None:
                expires_at, cached_result = cached_entry
                if now < expires_at:
                    return cached_result

            result = function(*args, **kwargs)
            cache[key] = (now + float(ttl_seconds), result)
            return result

        return wrapper

    return decorator


if __name__ == "__main__":
    @cache_function_results_with_ttl(30.0)
    def load_status(name: str) -> str:
        return f"ready:{name}"

    print(load_status("worker-a"))