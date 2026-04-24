# Graceful Worker Shutdown in Python

Run a queue-based worker that finishes accepted work, rejects new submissions during shutdown, and runs cleanup once.

This snippet is useful when a long-running Python process needs to stop cleanly after Ctrl+C, a SIGTERM, or an orchestrator shutdown request.

## Highlights

- Finishes queued work before exit
- Rejects new work after shutdown
- Restores signal handlers cleanly

## Use Cases

- Drain accepted work before a container stops
- Stop a background importer without losing queued tasks
- Close files, connections, or other resources after a stop request

## Code

```python
from src.graceful_worker_shutdown import GracefulWorker, install_shutdown_signal_handlers


def process_item(item):
    print(f"Processing {item}")


def cleanup():
    print("Closing resources")


worker = GracefulWorker(process_item, cleanup_handler=cleanup)
restore_handlers = install_shutdown_signal_handlers(worker)
worker.start()
worker.submit("task-1")
worker.submit("task-2")

try:
    worker.shutdown(timeout=5.0)
finally:
    restore_handlers()
```

## Notes

- Install signal handlers from the main thread before your process starts waiting for signals.
- `shutdown()` stops accepting new work and waits for accepted tasks to finish.
- Worker handler exceptions are logged and the worker continues with the next queued item.

## Verification

Run the unit tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- queued work execution
- rejecting submissions after shutdown
- draining accepted work during shutdown
- cleanup execution after queue drain
- shutdown timeout behavior
- handler exception recovery
- signal handler installation and restore

## Files

- `src/graceful_worker_shutdown.py`
- `tests/test_graceful_worker_shutdown.py`
- `snippet.json`