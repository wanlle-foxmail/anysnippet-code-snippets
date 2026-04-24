import signal
import sys
import threading
import time
import unittest
from pathlib import Path


SNIPPET_ROOT = Path(__file__).resolve().parent.parent
if str(SNIPPET_ROOT) not in sys.path:
    sys.path.insert(0, str(SNIPPET_ROOT))

from src.graceful_worker_shutdown import GracefulWorker, install_shutdown_signal_handlers


class GracefulWorkerTests(unittest.TestCase):
    def test_worker_processes_submitted_items(self):
        processed = []

        def handle_item(item):
            processed.append(item)

        worker = GracefulWorker(handle_item)
        worker.start()
        worker.submit("task-1")
        worker.submit("task-2")
        worker.submit("task-3")

        completed = worker.shutdown(timeout=1.0)

        self.assertTrue(completed)
        self.assertEqual(["task-1", "task-2", "task-3"], processed)

    def test_submit_raises_after_shutdown_is_requested(self):
        worker = GracefulWorker(lambda item: item)
        worker.start()
        worker.request_shutdown()

        with self.assertRaisesRegex(RuntimeError, "shutdown"):
            worker.submit("late-task")

        self.assertTrue(worker.shutdown(timeout=1.0))

    def test_shutdown_waits_for_queued_items_to_complete(self):
        processed = []

        def handle_item(item):
            time.sleep(0.02)
            processed.append(item)

        worker = GracefulWorker(handle_item)
        worker.start()
        for item in range(5):
            worker.submit(item)

        completed = worker.shutdown(timeout=1.0)

        self.assertTrue(completed)
        self.assertEqual([0, 1, 2, 3, 4], processed)

    def test_cleanup_runs_once_after_queue_is_empty(self):
        processed = []
        cleanup_snapshots = []

        def handle_item(item):
            processed.append(item)

        def cleanup():
            cleanup_snapshots.append(list(processed))

        worker = GracefulWorker(handle_item, cleanup_handler=cleanup)
        worker.start()
        worker.submit("first")
        worker.submit("second")

        completed = worker.shutdown(timeout=1.0)

        self.assertTrue(completed)
        self.assertEqual([["first", "second"]], cleanup_snapshots)

    def test_shutdown_timeout_returns_false_when_work_does_not_finish(self):
        started = threading.Event()
        release = threading.Event()
        cleanup_calls = []

        def handle_item(item):
            started.set()
            release.wait(timeout=1.0)

        def cleanup():
            cleanup_calls.append("done")

        worker = GracefulWorker(handle_item, cleanup_handler=cleanup)
        worker.start()
        worker.submit("blocking-task")
        self.assertTrue(started.wait(timeout=0.5))

        completed = worker.shutdown(timeout=0.05)

        self.assertFalse(completed)
        self.assertEqual([], cleanup_calls)

        release.set()
        self.assertTrue(worker.shutdown(timeout=1.0))
        self.assertEqual(["done"], cleanup_calls)

    def test_worker_continues_after_handler_exception(self):
        processed = []

        def handle_item(item):
            if item == "bad":
                raise ValueError("bad task")
            processed.append(item)

        worker = GracefulWorker(handle_item)
        worker.start()
        worker.submit("good-1")
        worker.submit("bad")
        worker.submit("good-2")

        completed = worker.shutdown(timeout=1.0)

        self.assertTrue(completed)
        self.assertEqual(["good-1", "good-2"], processed)

    def test_signal_handler_requests_shutdown_and_can_be_restored(self):
        previous_handler = signal.getsignal(signal.SIGINT)
        worker = GracefulWorker(lambda item: item)
        worker.start()

        restore_handlers = install_shutdown_signal_handlers(worker, signals=[signal.SIGINT])
        current_handler = signal.getsignal(signal.SIGINT)

        current_handler(signal.SIGINT, None)

        try:
            with self.assertRaisesRegex(RuntimeError, "shutdown"):
                worker.submit("late-task")

            self.assertTrue(worker.shutdown(timeout=1.0))
        finally:
            restore_handlers()

        self.assertEqual(previous_handler, signal.getsignal(signal.SIGINT))


if __name__ == "__main__":
    unittest.main()