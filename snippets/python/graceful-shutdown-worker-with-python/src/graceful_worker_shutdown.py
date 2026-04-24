import logging
import queue
import signal
import threading
from types import FrameType
from typing import Any, Callable, Iterable, Optional


LOGGER = logging.getLogger(__name__)


class GracefulWorker:
    """Run queued work until shutdown is requested and accepted work is drained."""

    def __init__(
        self,
        work_handler: Callable[[Any], None],
        cleanup_handler: Optional[Callable[[], None]] = None,
        poll_interval: float = 0.1,
    ) -> None:
        if not callable(work_handler):
            raise TypeError("work_handler must be callable")
        if cleanup_handler is not None and not callable(cleanup_handler):
            raise TypeError("cleanup_handler must be callable")
        if isinstance(poll_interval, bool) or not isinstance(poll_interval, (int, float)) or poll_interval <= 0:
            raise ValueError("poll_interval must be a positive number")

        self._work_handler = work_handler
        self._cleanup_handler = cleanup_handler
        self._poll_interval = float(poll_interval)
        self._tasks = queue.Queue()
        self._shutdown_requested = threading.Event()
        self._cleanup_completed = threading.Event()
        self._stopped = threading.Event()
        self._thread = None
        self._cleanup_lock = threading.Lock()

    def start(self) -> None:
        """Start the background worker thread."""
        if self._thread is not None:
            raise RuntimeError("worker has already been started")
        if self._shutdown_requested.is_set():
            raise RuntimeError("worker cannot start after shutdown has been requested")

        thread = threading.Thread(target=self._run, name="graceful-worker", daemon=True)
        self._thread = thread
        thread.start()

    def submit(self, item: Any) -> None:
        """Queue one item for background processing."""
        if self._shutdown_requested.is_set():
            raise RuntimeError("cannot submit work after shutdown has been requested")
        self._tasks.put(item)

    def request_shutdown(self) -> None:
        """Reject new work and let the worker finish accepted items."""
        self._shutdown_requested.set()

    def shutdown(self, timeout: Optional[float] = None) -> bool:
        """Request shutdown and wait for the worker thread to stop."""
        self.request_shutdown()

        if self._thread is None:
            self._run_cleanup_once()
            self._stopped.set()
            return True

        self._thread.join(timeout=timeout)
        return not self._thread.is_alive()

    def _run(self) -> None:
        # Flow:
        #   worker loop -> stop only after shutdown is requested and the queue is drained
        #                  |
        #                  +-> no item yet -> poll again
        #                  |
        #                  +-> item received -> process it, mark done, and continue
        #   on exit -----> run cleanup once and mark the worker as stopped
        try:
            while True:
                if self._shutdown_requested.is_set() and self._tasks.empty():
                    break

                try:
                    item = self._tasks.get(timeout=self._poll_interval)
                except queue.Empty:
                    continue

                try:
                    self._work_handler(item)
                except Exception:
                    LOGGER.exception("GracefulWorker failed to process item: %r", item)
                finally:
                    self._tasks.task_done()
        finally:
            self._run_cleanup_once()
            self._stopped.set()

    def _run_cleanup_once(self) -> None:
        if self._cleanup_handler is None or self._cleanup_completed.is_set():
            return

        with self._cleanup_lock:
            if self._cleanup_completed.is_set():
                return
            self._cleanup_handler()
            self._cleanup_completed.set()


def install_shutdown_signal_handlers(
    worker: GracefulWorker,
    signals: Optional[Iterable[int]] = None,
) -> Callable[[], None]:
    """Install signal handlers that request worker shutdown and return a restore callback."""
    signal_numbers = tuple(signals) if signals is not None else _default_shutdown_signals()
    previous_handlers = {}

    def handle_shutdown(signum: int, frame: Optional[FrameType]) -> None:
        del signum, frame
        worker.request_shutdown()

    for signum in signal_numbers:
        previous_handlers[signum] = signal.getsignal(signum)
        signal.signal(signum, handle_shutdown)

    def restore_handlers() -> None:
        for signum, previous_handler in previous_handlers.items():
            signal.signal(signum, previous_handler)

    return restore_handlers


def _default_shutdown_signals() -> tuple[int, ...]:
    signal_numbers = [signal.SIGINT]
    sigterm = getattr(signal, "SIGTERM", None)
    if sigterm is not None:
        signal_numbers.append(sigterm)
    return tuple(signal_numbers)


if __name__ == "__main__":
    def handle_item(item: str) -> None:
        print(f"Processing {item}")


    def cleanup() -> None:
        print("Cleanup completed")


    worker = GracefulWorker(handle_item, cleanup_handler=cleanup)
    restore_handlers = install_shutdown_signal_handlers(worker)
    worker.start()
    worker.submit("task-1")
    worker.submit("task-2")

    try:
        worker.shutdown(timeout=5.0)
    finally:
        restore_handlers()