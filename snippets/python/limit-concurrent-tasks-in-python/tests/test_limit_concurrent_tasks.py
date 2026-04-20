import sys
import threading
import time
import unittest
from pathlib import Path


SNIPPET_ROOT = Path(__file__).resolve().parent.parent
if str(SNIPPET_ROOT) not in sys.path:
    sys.path.insert(0, str(SNIPPET_ROOT))

from src.limit_concurrent_tasks import limit_concurrent_tasks


class LimitConcurrentTasksTests(unittest.TestCase):
    def test_runs_tasks_and_returns_results(self):
        result = limit_concurrent_tasks([1, 2, 3], lambda item: item * 2, max_workers=2)

        self.assertEqual([2, 4, 6], result)

    def test_preserves_input_order_even_when_completion_order_differs(self):
        def worker(item):
            time.sleep({1: 0.03, 2: 0.01, 3: 0.02}[item])
            return f"done-{item}"

        result = limit_concurrent_tasks([1, 2, 3], worker, max_workers=3)

        self.assertEqual(["done-1", "done-2", "done-3"], result)

    def test_propagates_worker_errors(self):
        def worker(item):
            if item == 2:
                raise RuntimeError("bad item")
            return item + 10

        with self.assertRaisesRegex(RuntimeError, "bad item"):
            limit_concurrent_tasks([1, 2, 3], worker, max_workers=2)

    def test_rejects_non_positive_max_workers(self):
        with self.assertRaises(ValueError):
            limit_concurrent_tasks([1], lambda item: item, max_workers=0)

        with self.assertRaises(ValueError):
            limit_concurrent_tasks([1], lambda item: item, max_workers=-1)

        with self.assertRaises(ValueError):
            limit_concurrent_tasks([1], lambda item: item, max_workers="2")

        with self.assertRaises(ValueError):
            limit_concurrent_tasks([1], lambda item: item, max_workers=None)

    def test_returns_empty_list_for_empty_input(self):
        result = limit_concurrent_tasks([], lambda item: item, max_workers=3)

        self.assertEqual([], result)

    def test_never_exceeds_configured_concurrency_limit(self):
        lock = threading.Lock()
        running_count = 0
        max_seen = 0

        def worker(item):
            nonlocal running_count, max_seen
            with lock:
                running_count += 1
                max_seen = max(max_seen, running_count)
            time.sleep(0.03)
            with lock:
                running_count -= 1
            return item

        result = limit_concurrent_tasks(range(6), worker, max_workers=2)

        self.assertEqual([0, 1, 2, 3, 4, 5], result)
        self.assertLessEqual(max_seen, 2)


if __name__ == "__main__":
    unittest.main()