import time
from typing import Any, Callable, Tuple, Type


def retry_failed_operation(
    operation: Callable[[], Any],
    max_attempts: int = 3,
    delay_seconds: float = 0.0,
    retry_exceptions: Tuple[Type[BaseException], ...] = (Exception,),
) -> Any:
    """Retry one callable on retryable errors until it succeeds or attempts run out."""
    # Flow:
    #   run operation
    #      |
    #      +-> success -> return result
    #      +-> retryable error -> sleep if needed -> retry
    #      `-> last or non-retryable error -> raise
    if not callable(operation):
        raise TypeError("operation must be callable")
    if isinstance(max_attempts, bool):
        raise TypeError("max_attempts must be an integer")
    if not isinstance(max_attempts, int):
        raise TypeError("max_attempts must be an integer")
    if max_attempts <= 0:
        raise ValueError("max_attempts must be greater than 0")
    if isinstance(delay_seconds, bool):
        raise TypeError("delay_seconds must be a number")
    if not isinstance(delay_seconds, (int, float)):
        raise TypeError("delay_seconds must be a number")
    if delay_seconds < 0:
        raise ValueError("delay_seconds must be greater than or equal to 0")

    for attempt_number in range(1, max_attempts + 1):
        try:
            return operation()
        except retry_exceptions:
            if attempt_number == max_attempts:
                raise
            if delay_seconds > 0:
                time.sleep(float(delay_seconds))

    raise RuntimeError("retry loop ended unexpectedly")


if __name__ == "__main__":
    attempts = {"count": 0}

    def sometimes_ready() -> str:
        attempts["count"] += 1
        if attempts["count"] < 3:
            raise TimeoutError("not ready yet")
        return "ready"

    print(retry_failed_operation(sometimes_ready, max_attempts=3, delay_seconds=0.1, retry_exceptions=(TimeoutError,)))