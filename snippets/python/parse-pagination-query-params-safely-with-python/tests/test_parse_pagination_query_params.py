import sys
import unittest
from pathlib import Path


SNIPPET_ROOT = Path(__file__).resolve().parent.parent
if str(SNIPPET_ROOT) not in sys.path:
    sys.path.insert(0, str(SNIPPET_ROOT))

from src.parse_pagination_query_params import (
    DEFAULT_PAGE,
    DEFAULT_PAGE_SIZE,
    MAX_PAGE_SIZE,
    PaginationParams,
    parse_pagination_query_params,
)


class ParsePaginationQueryParamsTests(unittest.TestCase):
    def test_parses_valid_page_and_page_size(self):
        params = parse_pagination_query_params({"page": "3", "page_size": "25"})

        self.assertEqual(PaginationParams(page=3, page_size=25, offset=50), params)

    def test_uses_defaults_when_values_are_missing(self):
        params = parse_pagination_query_params({})

        self.assertEqual(
            PaginationParams(page=DEFAULT_PAGE, page_size=DEFAULT_PAGE_SIZE, offset=0),
            params,
        )

    def test_caps_page_size_at_maximum(self):
        params = parse_pagination_query_params({"page": "2", "page_size": "500"})

        self.assertEqual(
            PaginationParams(page=2, page_size=MAX_PAGE_SIZE, offset=MAX_PAGE_SIZE),
            params,
        )

    def test_raises_value_error_for_non_integer_page(self):
        with self.assertRaises(ValueError) as raised:
            parse_pagination_query_params({"page": "first"})

        self.assertIn("page", str(raised.exception))

    def test_raises_value_error_for_non_positive_page_size(self):
        with self.assertRaises(ValueError) as raised:
            parse_pagination_query_params({"page_size": "0"})

        self.assertIn("page_size", str(raised.exception))

    def test_accepts_whitespace_around_integer_values(self):
        params = parse_pagination_query_params({"page": " 2 ", "page_size": " 10 "})

        self.assertEqual(PaginationParams(page=2, page_size=10, offset=10), params)


if __name__ == "__main__":
    unittest.main()