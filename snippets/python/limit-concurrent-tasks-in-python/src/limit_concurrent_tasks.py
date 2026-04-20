from concurrent.futures import ThreadPoolExecutor
from typing import Any, Callable, Iterable


def limit_concurrent_tasks(
    items: Iterable[Any],
    worker: Callable[[Any], Any],
    max_workers: int = 5,
) -> list[Any]:
    """Run I/O-bound tasks with a fixed thread limit."""
    if isinstance(max_workers, bool) or not isinstance(max_workers, int) or max_workers <= 0:
        raise ValueError("max_workers must be a positive integer")

    with ThreadPoolExecutor(max_workers=max_workers) as executor:
        return list(executor.map(worker, items))


if __name__ == "__main__":
    def worker(item: int) -> int:
        return item * 2


    result = limit_concurrent_tasks([1, 2, 3], worker, max_workers=2)
    print(result)