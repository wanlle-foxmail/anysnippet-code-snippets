import sys
import tempfile
import unittest
from pathlib import Path


SNIPPET_ROOT = Path(__file__).resolve().parent.parent
if str(SNIPPET_ROOT) not in sys.path:
    sys.path.insert(0, str(SNIPPET_ROOT))

from src.read_headerless_csv_in_chunks import read_headerless_csv_in_chunks


class ReadHeaderlessCsvInChunksTests(unittest.TestCase):
    def test_yields_row_dictionaries_for_single_chunk(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            csv_path = Path(tmp_dir) / "people.csv"
            csv_path.write_text("Alice,29,Shanghai\nBob,31,Beijing\n", encoding="utf-8")

            result = list(read_headerless_csv_in_chunks(str(csv_path), ["name", "age", "city"]))

            self.assertEqual(
                [
                    {"name": "Alice", "age": 29, "city": "Shanghai"},
                    {"name": "Bob", "age": 31, "city": "Beijing"},
                ],
                result,
            )

    def test_yields_all_rows_across_multiple_chunks_in_order(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            csv_path = Path(tmp_dir) / "rows.csv"
            csv_path.write_text("1,A\n2,B\n3,C\n4,D\n5,E\n", encoding="utf-8")

            result = list(read_headerless_csv_in_chunks(str(csv_path), ["row_id", "name"], chunk_size=2))

            self.assertEqual(
                [
                    {"row_id": 1, "name": "A"},
                    {"row_id": 2, "name": "B"},
                    {"row_id": 3, "name": "C"},
                    {"row_id": 4, "name": "D"},
                    {"row_id": 5, "name": "E"},
                ],
                result,
            )

    def test_returns_empty_list_for_empty_csv(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            csv_path = Path(tmp_dir) / "empty.csv"
            csv_path.write_text("", encoding="utf-8")

            result = list(read_headerless_csv_in_chunks(str(csv_path), ["name", "age"]))

            self.assertEqual([], result)

    def test_yields_single_row_for_one_line_csv(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            csv_path = Path(tmp_dir) / "single.csv"
            csv_path.write_text("Alice,29,Shanghai\n", encoding="utf-8")

            result = list(read_headerless_csv_in_chunks(str(csv_path), ["name", "age", "city"]))

            self.assertEqual(
                [{"name": "Alice", "age": 29, "city": "Shanghai"}],
                result,
            )

    def test_raises_file_not_found_for_missing_csv(self):
        with self.assertRaises(FileNotFoundError):
            list(read_headerless_csv_in_chunks("missing.csv", ["name", "age"]))


if __name__ == "__main__":
    unittest.main()