import sys
import unittest
from pathlib import Path
from unittest.mock import patch


SNIPPET_ROOT = Path(__file__).resolve().parent.parent
if str(SNIPPET_ROOT) not in sys.path:
    sys.path.insert(0, str(SNIPPET_ROOT))

from src.cache_function_results_with_ttl import cache_function_results_with_ttl


class CacheFunctionResultsWithTtlTests(unittest.TestCase):
    def test_reuses_cached_result_within_ttl(self):
        calls = []

        @cache_function_results_with_ttl(10.0)
        def load_value(value):
            calls.append(value)
            return f"result-{len(calls)}"

        with patch("src.cache_function_results_with_ttl.time.monotonic", side_effect=[100.0, 105.0]):
            first = load_value("alpha")
            second = load_value("alpha")

        self.assertEqual("result-1", first)
        self.assertEqual("result-1", second)
        self.assertEqual(["alpha"], calls)

    def test_recomputes_result_after_ttl_expires(self):
        calls = []

        @cache_function_results_with_ttl(10.0)
        def load_value(value):
            calls.append(value)
            return f"result-{len(calls)}"

        with patch("src.cache_function_results_with_ttl.time.monotonic", side_effect=[100.0, 111.0]):
            first = load_value("alpha")
            second = load_value("alpha")

        self.assertEqual("result-1", first)
        self.assertEqual("result-2", second)
        self.assertEqual(["alpha", "alpha"], calls)

    def test_separates_cache_entries_by_arguments(self):
        calls = []

        @cache_function_results_with_ttl(10.0)
        def load_value(value):
            calls.append(value)
            return f"result-{value}-{len(calls)}"

        with patch("src.cache_function_results_with_ttl.time.monotonic", side_effect=[100.0, 101.0, 102.0]):
            first = load_value("alpha")
            second = load_value("beta")
            third = load_value("alpha")

        self.assertEqual("result-alpha-1", first)
        self.assertEqual("result-beta-2", second)
        self.assertEqual("result-alpha-1", third)
        self.assertEqual(["alpha", "beta"], calls)

    def test_normalizes_keyword_argument_order(self):
        calls = []

        @cache_function_results_with_ttl(10.0)
        def build_label(*, name, prefix):
            calls.append((name, prefix))
            return f"{prefix}-{name}-{len(calls)}"

        with patch("src.cache_function_results_with_ttl.time.monotonic", side_effect=[100.0, 105.0]):
            first = build_label(name="worker", prefix="job")
            second = build_label(prefix="job", name="worker")

        self.assertEqual("job-worker-1", first)
        self.assertEqual("job-worker-1", second)
        self.assertEqual([("worker", "job")], calls)

    def test_does_not_cache_exceptions(self):
        call_count = 0

        @cache_function_results_with_ttl(10.0)
        def load_value():
            nonlocal call_count
            call_count += 1
            if call_count == 1:
                raise RuntimeError("temporary failure")
            return "ok"

        with patch("src.cache_function_results_with_ttl.time.monotonic", side_effect=[100.0, 101.0]):
            with self.assertRaises(RuntimeError):
                load_value()
            result = load_value()

        self.assertEqual("ok", result)
        self.assertEqual(2, call_count)

    def test_rejects_non_positive_or_boolean_ttl(self):
        with self.assertRaises(ValueError):
            cache_function_results_with_ttl(0)

        with self.assertRaises(TypeError):
            cache_function_results_with_ttl(True)


if __name__ == "__main__":
    unittest.main()