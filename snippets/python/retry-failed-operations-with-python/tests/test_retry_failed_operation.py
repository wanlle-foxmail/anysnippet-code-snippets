import sys
import unittest
from pathlib import Path
from unittest.mock import patch


SNIPPET_ROOT = Path(__file__).resolve().parent.parent
if str(SNIPPET_ROOT) not in sys.path:
    sys.path.insert(0, str(SNIPPET_ROOT))

from src.retry_failed_operation import retry_failed_operation


class RetryFailedOperationTests(unittest.TestCase):
    def test_returns_result_on_first_success(self):
        calls = []

        def operation():
            calls.append("run")
            return "ok"

        result = retry_failed_operation(operation)

        self.assertEqual("ok", result)
        self.assertEqual(["run"], calls)

    def test_retries_retryable_error_until_success(self):
        call_count = 0

        def operation():
            nonlocal call_count
            call_count += 1
            if call_count < 3:
                raise TimeoutError("temporary failure")
            return "ok"

        with patch("src.retry_failed_operation.time.sleep") as mocked_sleep:
            result = retry_failed_operation(operation, max_attempts=3, delay_seconds=0.5, retry_exceptions=(TimeoutError,))

        self.assertEqual("ok", result)
        self.assertEqual(3, call_count)
        self.assertEqual([(0.5,), (0.5,)], [call.args for call in mocked_sleep.call_args_list])

    def test_raises_last_retryable_error_after_max_attempts(self):
        call_count = 0

        def operation():
            nonlocal call_count
            call_count += 1
            raise TimeoutError("still failing")

        with patch("src.retry_failed_operation.time.sleep") as mocked_sleep:
            with self.assertRaises(TimeoutError) as raised:
                retry_failed_operation(operation, max_attempts=3, delay_seconds=0.25, retry_exceptions=(TimeoutError,))

        self.assertEqual("still failing", str(raised.exception))
        self.assertEqual(3, call_count)
        self.assertEqual(2, mocked_sleep.call_count)

    def test_does_not_retry_non_retryable_error(self):
        call_count = 0

        def operation():
            nonlocal call_count
            call_count += 1
            raise ValueError("bad input")

        with patch("src.retry_failed_operation.time.sleep") as mocked_sleep:
            with self.assertRaises(ValueError):
                retry_failed_operation(operation, max_attempts=3, delay_seconds=0.5, retry_exceptions=(TimeoutError,))

        self.assertEqual(1, call_count)
        self.assertEqual(0, mocked_sleep.call_count)

    def test_skips_sleep_when_success_happens_on_last_allowed_attempt(self):
        call_count = 0

        def operation():
            nonlocal call_count
            call_count += 1
            if call_count < 2:
                raise TimeoutError("retry once")
            return "done"

        with patch("src.retry_failed_operation.time.sleep") as mocked_sleep:
            result = retry_failed_operation(operation, max_attempts=2, delay_seconds=1.0, retry_exceptions=(TimeoutError,))

        self.assertEqual("done", result)
        self.assertEqual(1, mocked_sleep.call_count)

    def test_rejects_invalid_max_attempts(self):
        with self.assertRaises(ValueError):
            retry_failed_operation(lambda: None, max_attempts=0)

        with self.assertRaises(TypeError):
            retry_failed_operation(lambda: None, max_attempts=True)

    def test_rejects_invalid_delay_seconds(self):
        with self.assertRaises(ValueError):
            retry_failed_operation(lambda: None, delay_seconds=-1.0)

        with self.assertRaises(TypeError):
            retry_failed_operation(lambda: None, delay_seconds=False)


if __name__ == "__main__":
    unittest.main()