import sys
import tempfile
import unittest
from pathlib import Path

import pandas


SNIPPET_ROOT = Path(__file__).resolve().parent.parent
if str(SNIPPET_ROOT) not in sys.path:
    sys.path.insert(0, str(SNIPPET_ROOT))

from src.read_parquet_records import read_parquet_records


class ReadParquetRecordsTests(unittest.TestCase):
    def test_reads_multiple_records_in_order(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            parquet_path = Path(tmp_dir) / "events.parquet"
            pandas.DataFrame(
                [
                    {"id": 1, "name": "alpha", "active": True},
                    {"id": 2, "name": "beta", "active": False},
                ]
            ).to_parquet(parquet_path, engine="pyarrow", index=False)

            records = read_parquet_records(str(parquet_path))

            self.assertEqual(
                [
                    {"id": 1, "name": "alpha", "active": True},
                    {"id": 2, "name": "beta", "active": False},
                ],
                records,
            )

    def test_returns_empty_list_for_empty_parquet_file(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            parquet_path = Path(tmp_dir) / "empty.parquet"
            pandas.DataFrame(columns=["id", "name"]).to_parquet(parquet_path, engine="pyarrow", index=False)

            records = read_parquet_records(str(parquet_path))

            self.assertEqual([], records)

    def test_preserves_unicode_text_values(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            parquet_path = Path(tmp_dir) / "unicode.parquet"
            pandas.DataFrame(
                [{"message": "你好", "emoji": "🚀"}]
            ).to_parquet(parquet_path, engine="pyarrow", index=False)

            records = read_parquet_records(str(parquet_path))

            self.assertEqual([{"message": "你好", "emoji": "🚀"}], records)

    def test_preserves_none_values_in_object_columns(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            parquet_path = Path(tmp_dir) / "nulls.parquet"
            pandas.DataFrame(
                [
                    {"name": "alpha", "notes": None},
                    {"name": None, "notes": "beta"},
                ]
            ).to_parquet(parquet_path, engine="pyarrow", index=False)

            records = read_parquet_records(str(parquet_path))

            self.assertEqual(
                [
                    {"name": "alpha", "notes": None},
                    {"name": None, "notes": "beta"},
                ],
                records,
            )

    def test_raises_file_not_found_for_missing_file(self):
        with self.assertRaises(FileNotFoundError):
            read_parquet_records("missing.parquet")

    def test_raises_type_error_for_non_string_path(self):
        with self.assertRaises(TypeError):
            read_parquet_records(None)

        with self.assertRaises(TypeError):
            read_parquet_records(123)

        with self.assertRaises(TypeError):
            read_parquet_records(True)

    def test_raises_value_error_for_empty_path(self):
        with self.assertRaises(ValueError):
            read_parquet_records("")

        with self.assertRaises(ValueError):
            read_parquet_records("   ")


if __name__ == "__main__":
    unittest.main()