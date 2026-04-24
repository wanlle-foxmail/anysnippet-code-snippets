import sys
import unittest
from pathlib import Path


SNIPPET_ROOT = Path(__file__).resolve().parent.parent
if str(SNIPPET_ROOT) not in sys.path:
    sys.path.insert(0, str(SNIPPET_ROOT))

from src.distribute_urls_evenly_by_domain import distribute_urls_evenly_by_domain


class DistributeUrlsEvenlyByDomainTests(unittest.TestCase):
    def test_returns_empty_list_for_empty_input(self):
        ordered_urls = distribute_urls_evenly_by_domain({})

        self.assertEqual([], ordered_urls)

    def test_preserves_input_order_for_single_domain(self):
        ordered_urls = distribute_urls_evenly_by_domain(
            {
                "https://news-a.example.com/story-001": "news-a.example.com",
                "https://news-a.example.com/story-002": "news-a.example.com",
                "https://news-a.example.com/story-003": "news-a.example.com",
            }
        )

        self.assertEqual(
            [
                "https://news-a.example.com/story-001",
                "https://news-a.example.com/story-002",
                "https://news-a.example.com/story-003",
            ],
            ordered_urls,
        )

    def test_alternates_domains_when_counts_are_balanced(self):
        ordered_urls = distribute_urls_evenly_by_domain(
            {
                "https://news-a.example.com/story-001": "news-a.example.com",
                "https://news-b.example.com/story-010": "news-b.example.com",
                "https://news-a.example.com/story-002": "news-a.example.com",
                "https://news-b.example.com/story-011": "news-b.example.com",
            }
        )

        self.assertEqual(
            [
                "https://news-a.example.com/story-001",
                "https://news-b.example.com/story-010",
                "https://news-a.example.com/story-002",
                "https://news-b.example.com/story-011",
            ],
            ordered_urls,
        )

    def test_avoids_repeating_previous_domain_when_alternatives_exist(self):
        ordered_urls = distribute_urls_evenly_by_domain(
            {
                "https://news-a.example.com/story-001": "news-a.example.com",
                "https://news-b.example.com/story-010": "news-b.example.com",
                "https://news-a.example.com/story-002": "news-a.example.com",
                "https://news-c.example.com/story-201": "news-c.example.com",
                "https://news-a.example.com/story-003": "news-a.example.com",
                "https://news-b.example.com/story-011": "news-b.example.com",
                "https://news-a.example.com/story-004": "news-a.example.com",
            }
        )

        self.assertEqual(
            [
                "https://news-a.example.com/story-001",
                "https://news-b.example.com/story-010",
                "https://news-a.example.com/story-002",
                "https://news-b.example.com/story-011",
                "https://news-a.example.com/story-003",
                "https://news-c.example.com/story-201",
                "https://news-a.example.com/story-004",
            ],
            ordered_urls,
        )

    def test_prefers_domain_with_more_remaining_urls(self):
        ordered_urls = distribute_urls_evenly_by_domain(
            {
                "https://news-a.example.com/story-001": "news-a.example.com",
                "https://news-b.example.com/story-010": "news-b.example.com",
                "https://news-c.example.com/story-201": "news-c.example.com",
                "https://news-a.example.com/story-002": "news-a.example.com",
                "https://news-c.example.com/story-202": "news-c.example.com",
                "https://news-a.example.com/story-003": "news-a.example.com",
            }
        )

        self.assertEqual(
            [
                "https://news-a.example.com/story-001",
                "https://news-c.example.com/story-201",
                "https://news-a.example.com/story-002",
                "https://news-b.example.com/story-010",
                "https://news-a.example.com/story-003",
                "https://news-c.example.com/story-202",
            ],
            ordered_urls,
        )

    def test_groups_domains_by_exact_string_value(self):
        ordered_urls = distribute_urls_evenly_by_domain(
            {
                "https://example.com/story-001": "example.com",
                "https://EXAMPLE.COM/story-010": "EXAMPLE.COM",
                "https://www.example.com/story-201": "www.example.com",
                "https://example.com/story-002": "example.com",
            }
        )

        self.assertEqual(
            [
                "https://example.com/story-001",
                "https://EXAMPLE.COM/story-010",
                "https://example.com/story-002",
                "https://www.example.com/story-201",
            ],
            ordered_urls,
        )

    def test_rejects_non_dict_inputs_and_non_string_entries(self):
        with self.assertRaises(TypeError):
            distribute_urls_evenly_by_domain([("https://example.com/story-001", "example.com")])

        with self.assertRaises(TypeError):
            distribute_urls_evenly_by_domain({1: "example.com"})

        with self.assertRaises(TypeError):
            distribute_urls_evenly_by_domain({"https://example.com/story-001": None})


if __name__ == "__main__":
    unittest.main()