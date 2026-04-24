import sys
import unittest
from pathlib import Path


SNIPPET_ROOT = Path(__file__).resolve().parent.parent
if str(SNIPPET_ROOT) not in sys.path:
    sys.path.insert(0, str(SNIPPET_ROOT))

from src.parse_sort_query_params import SortField, parse_sort_query_params


class ParseSortQueryParamsTests(unittest.TestCase):
    def test_parses_single_ascending_field(self):
        sort_fields = parse_sort_query_params({"sort": "name"}, ["name", "created_at"])

        self.assertEqual([SortField(field_name="name", descending=False)], sort_fields)

    def test_parses_multiple_fields_with_mixed_directions(self):
        sort_fields = parse_sort_query_params(
            {"sort": "name,-created_at"},
            ["name", "created_at"],
        )

        self.assertEqual(
            [
                SortField(field_name="name", descending=False),
                SortField(field_name="created_at", descending=True),
            ],
            sort_fields,
        )

    def test_returns_empty_list_when_sort_param_is_missing(self):
        sort_fields = parse_sort_query_params({}, ["name", "created_at"])

        self.assertEqual([], sort_fields)

    def test_rejects_unknown_sort_fields(self):
        with self.assertRaises(ValueError) as raised:
            parse_sort_query_params({"sort": "role"}, ["name", "created_at"])

        self.assertIn("not allowed", str(raised.exception))

    def test_rejects_empty_sort_terms(self):
        with self.assertRaises(ValueError) as raised:
            parse_sort_query_params({"sort": "name, "}, ["name", "created_at"])

        self.assertIn("must not be empty", str(raised.exception))

    def test_rejects_duplicate_sort_fields(self):
        with self.assertRaises(ValueError) as raised:
            parse_sort_query_params({"sort": "name,-name"}, ["name", "created_at"])

        self.assertIn("must not be repeated", str(raised.exception))


if __name__ == "__main__":
    unittest.main()