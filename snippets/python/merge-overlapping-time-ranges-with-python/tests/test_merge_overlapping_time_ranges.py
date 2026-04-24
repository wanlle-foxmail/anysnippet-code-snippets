import sys
import unittest
from pathlib import Path


SNIPPET_ROOT = Path(__file__).resolve().parent.parent
if str(SNIPPET_ROOT) not in sys.path:
    sys.path.insert(0, str(SNIPPET_ROOT))

from src.merge_overlapping_time_ranges import merge_overlapping_time_ranges


class MergeOverlappingTimeRangesTests(unittest.TestCase):
    def test_merges_unsorted_overlapping_ranges(self):
        ranges = [(5, 7), (1, 3), (2, 6)]

        merged_ranges = merge_overlapping_time_ranges(ranges)

        self.assertEqual([(1, 7)], merged_ranges)

    def test_keeps_separate_non_overlapping_ranges(self):
        ranges = [(1, 2), (4, 5), (7, 8)]

        merged_ranges = merge_overlapping_time_ranges(ranges)

        self.assertEqual([(1, 2), (4, 5), (7, 8)], merged_ranges)

    def test_merges_ranges_that_share_a_boundary(self):
        ranges = [(1, 3), (3, 5), (5, 8)]

        merged_ranges = merge_overlapping_time_ranges(ranges)

        self.assertEqual([(1, 8)], merged_ranges)

    def test_keeps_outer_range_when_inner_ranges_are_nested(self):
        ranges = [(1, 10), (2, 3), (4, 8)]

        merged_ranges = merge_overlapping_time_ranges(ranges)

        self.assertEqual([(1, 10)], merged_ranges)

    def test_returns_empty_list_for_empty_input(self):
        self.assertEqual([], merge_overlapping_time_ranges([]))

    def test_raises_value_error_for_reversed_range(self):
        with self.assertRaises(ValueError):
            merge_overlapping_time_ranges([(5, 3)])

    def test_rejects_bool_values_in_ranges(self):
        with self.assertRaises(TypeError):
            merge_overlapping_time_ranges([(True, 3)])

        with self.assertRaises(TypeError):
            merge_overlapping_time_ranges([(1, False)])


if __name__ == "__main__":
    unittest.main()