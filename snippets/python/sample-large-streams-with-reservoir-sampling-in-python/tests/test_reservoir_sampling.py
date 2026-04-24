import sys
import unittest
from pathlib import Path


SNIPPET_ROOT = Path(__file__).resolve().parent.parent
if str(SNIPPET_ROOT) not in sys.path:
    sys.path.insert(0, str(SNIPPET_ROOT))

from src.reservoir_sampling import reservoir_sample


class ReservoirSamplingTests(unittest.TestCase):
    def test_returns_deterministic_sample_with_seed(self):
        sampled_items = reservoir_sample(range(10), 3, seed=7)

        self.assertEqual([7, 8, 3], sampled_items)

    def test_returns_all_items_when_stream_is_shorter_than_sample(self):
        sampled_items = reservoir_sample(["a", "b"], 5, seed=7)

        self.assertEqual(["a", "b"], sampled_items)

    def test_returns_empty_list_for_empty_iterable(self):
        sampled_items = reservoir_sample([], 3, seed=7)

        self.assertEqual([], sampled_items)

    def test_works_with_generator_input(self):
        sampled_items = reservoir_sample((item for item in range(10)), 4, seed=11)

        self.assertEqual([0, 1, 2, 4], sampled_items)

    def test_returns_original_order_when_sample_matches_stream_length(self):
        sampled_items = reservoir_sample([10, 20, 30], 3, seed=99)

        self.assertEqual([10, 20, 30], sampled_items)

    def test_rejects_invalid_sample_size(self):
        with self.assertRaises(ValueError):
            reservoir_sample([1, 2, 3], 0)

        with self.assertRaises(ValueError):
            reservoir_sample([1, 2, 3], -1)

        with self.assertRaises(ValueError):
            reservoir_sample([1, 2, 3], True)


if __name__ == "__main__":
    unittest.main()