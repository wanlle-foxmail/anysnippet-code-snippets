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
    def test_runs_tasks_and_returns_summary_counts(self):
        result = limit_concurrent_tasks([1, 2, 3], lambda item: item * 2, max_workers=2)

        self.assertEqual(2, result["max_workers"])
        self.assertEqual(3, result["submitted_count"])
        self.assertEqual(3, result["succeeded_count"])
        self.assertEqual(0, result["failed_count"])
        self.assertEqual(
            [
                {"input_item": 1, "status": "ok", "value": 2, "error_message": None},
                {"input_item": 2, "status": "ok", "value": 4, "error_message": None},
                {"input_item": 3, "status": "ok", "value": 6, "error_message": None},
            ],
            result["results"],
        )

    def test_preserves_input_order_even_when_completion_order_differs(self):
        def worker(item):
            time.sleep({1: 0.03, 2: 0.01, 3: 0.02}[item])
            return f"done-{item}"

        result = limit_concurrent_tasks([1, 2, 3], worker, max_workers=3)

        self.assertEqual(
            ["done-1", "done-2", "done-3"],
            [outcome["value"] for outcome in result["results"]],
        )

    def test_records_task_errors_without_stopping_other_tasks(self):
        def worker(item):
            if item == 2:
                raise RuntimeError("bad item")
            return item + 10

        result = limit_concurrent_tasks([1, 2, 3], worker, max_workers=2)

        self.assertEqual(3, result["submitted_count"])
        self.assertEqual(2, result["succeeded_count"])
        self.assertEqual(1, result["failed_count"])
        self.assertEqual("ok", result["results"][0]["status"])
        self.assertEqual("error", result["results"][1]["status"])
        self.assertIn("bad item", result["results"][1]["error_message"])
        self.assertEqual("ok", result["results"][2]["status"])

    def test_rejects_non_positive_max_workers(self):
        with self.assertRaises(ValueError):
            limit_concurrent_tasks([1], lambda item: item, max_workers=0)

        with self.assertRaises(ValueError):
            limit_concurrent_tasks([1], lambda item: item, max_workers=-1)

        with self.assertRaises(ValueError):
            limit_concurrent_tasks([1], lambda item: item, max_workers="2")

        with self.assertRaises(ValueError):
            limit_concurrent_tasks([1], lambda item: item, max_workers=None)

    def test_returns_empty_summary_for_empty_input(self):
        call_count = 0

        def worker(item):
            nonlocal call_count
            call_count += 1
            return item

        result = limit_concurrent_tasks([], worker, max_workers=3)

        self.assertEqual(0, result["submitted_count"])
        self.assertEqual(0, result["succeeded_count"])
        self.assertEqual(0, result["failed_count"])
        self.assertEqual([], result["results"])
        self.assertEqual(0, call_count)

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

        self.assertEqual(6, result["succeeded_count"])
        self.assertLessEqual(max_seen, 2)

    def test_consumes_generator_input(self):
        produced_count = 0
        result_holder = {}
        ready_to_assert = threading.Event()
        release_workers = threading.Event()
        lock = threading.Lock()
        started_count = 0

        def generate_items():
            nonlocal produced_count
            for item in (1, 2, 3):
                produced_count += 1
                yield item

        def worker(item):
            nonlocal started_count
            with lock:
                started_count += 1
                if started_count == 2:
                    ready_to_assert.set()
            if item in {1, 2}:
                release_workers.wait(timeout=1)
            return item * 3

        thread = threading.Thread(
            target=lambda: result_holder.setdefault(
                "result",
                limit_concurrent_tasks(generate_items(), worker, max_workers=2),
            )
        )
        thread.start()

        self.assertTrue(ready_to_assert.wait(timeout=1))
        self.assertEqual(2, produced_count)
        release_workers.set()
        thread.join(timeout=2)

        result = result_holder["result"]

        self.assertEqual(3, result["submitted_count"])
        self.assertEqual([3, 6, 9], [outcome["value"] for outcome in result["results"]])


if __name__ == "__main__":
    unittest.main()