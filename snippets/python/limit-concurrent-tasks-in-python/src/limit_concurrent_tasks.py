from concurrent.futures import FIRST_COMPLETED, ThreadPoolExecutor, wait
from typing import Any, Callable, Iterable, Optional, TypedDict


class TaskOutcome(TypedDict):
    input_item: Any
    status: str
    value: Any
    error_message: Optional[str]


class LimitedConcurrencyRunResult(TypedDict):
    max_workers: int
    submitted_count: int
    succeeded_count: int
    failed_count: int
    results: list[TaskOutcome]


def limit_concurrent_tasks(
    items: Iterable[Any],
    worker: Callable[[Any], Any],
    *,
    max_workers: int = 5,
) -> LimitedConcurrencyRunResult:
    """Run tasks with a fixed concurrency limit and stable output ordering.

    Args:
        items: Iterable of input items to process.
        worker: Function that processes a single input item.
        max_workers: Maximum number of concurrent worker threads.

    Returns:
        A dictionary containing execution counts and one outcome per input item.

    Raises:
        ValueError: If max_workers is not a positive integer.
    """
    if not isinstance(max_workers, int) or isinstance(max_workers, bool) or max_workers <= 0:
        raise ValueError("max_workers must be greater than 0")

    indexed_items = enumerate(items)
    pending_futures = {}
    results_by_index: dict[int, TaskOutcome] = {}
    succeeded_count = 0
    failed_count = 0
    submitted_count = 0

    def submit_next(executor: ThreadPoolExecutor) -> bool:
        nonlocal submitted_count
        try:
            index, item = next(indexed_items)
        except StopIteration:
            return False

        pending_futures[executor.submit(worker, item)] = (index, item)
        submitted_count += 1
        return True

    with ThreadPoolExecutor(max_workers=max_workers) as executor:
        for _ in range(max_workers):
            if not submit_next(executor):
                break

        if submitted_count == 0:
            return {
                "max_workers": max_workers,
                "submitted_count": 0,
                "succeeded_count": 0,
                "failed_count": 0,
                "results": [],
            }

        while pending_futures:
            completed_futures, _ = wait(pending_futures, return_when=FIRST_COMPLETED)
            for future in completed_futures:
                index, input_item = pending_futures.pop(future)
                try:
                    value = future.result()
                except Exception as error:
                    failed_count += 1
                    results_by_index[index] = {
                        "input_item": input_item,
                        "status": "error",
                        "value": None,
                        "error_message": str(error),
                    }
                else:
                    succeeded_count += 1
                    results_by_index[index] = {
                        "input_item": input_item,
                        "status": "ok",
                        "value": value,
                        "error_message": None,
                    }

                submit_next(executor)

    ordered_indexes = sorted(results_by_index)
    if not ordered_indexes:
        return {
            "max_workers": max_workers,
            "submitted_count": 0,
            "succeeded_count": 0,
            "failed_count": 0,
            "results": [],
        }

    return {
        "max_workers": max_workers,
        "submitted_count": submitted_count,
        "succeeded_count": succeeded_count,
        "failed_count": failed_count,
        "results": [results_by_index[index] for index in ordered_indexes],
    }