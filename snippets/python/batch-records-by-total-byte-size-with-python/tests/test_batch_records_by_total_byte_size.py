import sys
import unittest
from pathlib import Path


SNIPPET_ROOT = Path(__file__).resolve().parent.parent
if str(SNIPPET_ROOT) not in sys.path:
    sys.path.insert(0, str(SNIPPET_ROOT))

from src.batch_records_by_total_byte_size import batch_records_by_total_byte_size


class BatchRecordsByTotalByteSizeTests(unittest.TestCase):
    def test_batches_ascii_records_by_total_bytes(self):
        batches = list(batch_records_by_total_byte_size(["aaa", "bb", "c", "dddd"], 5))

        self.assertEqual([["aaa", "bb"], ["c", "dddd"]], batches)

    def test_handles_unicode_byte_sizes(self):
        batches = list(batch_records_by_total_byte_size(["你", "ab", "好"], 5))

        self.assertEqual([["你", "ab"], ["好"]], batches)

    def test_preserves_input_order_for_generator_input(self):
        batches = list(batch_records_by_total_byte_size((item for item in ["aa", "bbb", "c", "dd"]), 4))

        self.assertEqual([["aa"], ["bbb", "c"], ["dd"]], batches)

    def test_returns_no_batches_for_empty_input(self):
        batches = list(batch_records_by_total_byte_size([], 5))

        self.assertEqual([], batches)

    def test_raises_value_error_for_record_that_exceeds_limit(self):
        with self.assertRaises(ValueError) as raised:
            list(batch_records_by_total_byte_size(["toolong"], 3))

        self.assertIn("single record exceeds", str(raised.exception))

    def test_rejects_invalid_max_batch_bytes(self):
        with self.assertRaises(ValueError):
            list(batch_records_by_total_byte_size(["a"], 0))

        with self.assertRaises(ValueError):
            list(batch_records_by_total_byte_size(["a"], -1))

        with self.assertRaises(ValueError):
            list(batch_records_by_total_byte_size(["a"], True))

    def test_rejects_non_string_records(self):
        with self.assertRaises(TypeError):
            list(batch_records_by_total_byte_size(["a", 1], 5))

        with self.assertRaises(TypeError):
            list(batch_records_by_total_byte_size("abc", 5))


if __name__ == "__main__":
    unittest.main()