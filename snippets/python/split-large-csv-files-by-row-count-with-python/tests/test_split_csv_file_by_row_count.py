import csv
import sys
import tempfile
import unittest
from pathlib import Path


SNIPPET_ROOT = Path(__file__).resolve().parent.parent
if str(SNIPPET_ROOT) not in sys.path:
    sys.path.insert(0, str(SNIPPET_ROOT))

from src.split_csv_file_by_row_count import split_csv_file_by_row_count


def read_csv_rows(path: str) -> list[list[str]]:
    with open(path, "r", encoding="utf-8", newline="") as input_file:
        return list(csv.reader(input_file))


class SplitCsvFileByRowCountTests(unittest.TestCase):
    def test_splits_rows_and_repeats_header(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            input_path = Path(tmp_dir) / "events.csv"
            output_dir = Path(tmp_dir) / "parts"
            input_path.write_text(
                "id,name\n1,alpha\n2,beta\n3,gamma\n4,delta\n5,epsilon\n",
                encoding="utf-8",
            )

            written_files = split_csv_file_by_row_count(str(input_path), str(output_dir), 2)

            self.assertEqual(3, len(written_files))
            self.assertEqual(
                [["id", "name"], ["1", "alpha"], ["2", "beta"]],
                read_csv_rows(written_files[0]),
            )
            self.assertEqual(
                [["id", "name"], ["3", "gamma"], ["4", "delta"]],
                read_csv_rows(written_files[1]),
            )
            self.assertEqual(
                [["id", "name"], ["5", "epsilon"]],
                read_csv_rows(written_files[2]),
            )

    def test_creates_output_directory_when_missing(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            input_path = Path(tmp_dir) / "events.csv"
            output_dir = Path(tmp_dir) / "nested" / "parts"
            input_path.write_text("id,name\n1,alpha\n", encoding="utf-8")

            written_files = split_csv_file_by_row_count(str(input_path), str(output_dir), 1)

            self.assertTrue(output_dir.exists())
            self.assertEqual(1, len(written_files))

    def test_returns_single_file_when_limit_exceeds_row_count(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            input_path = Path(tmp_dir) / "events.csv"
            output_dir = Path(tmp_dir) / "parts"
            input_path.write_text("id,name\n1,alpha\n2,beta\n", encoding="utf-8")

            written_files = split_csv_file_by_row_count(str(input_path), str(output_dir), 10)

            self.assertEqual(1, len(written_files))

    def test_returns_no_files_for_header_only_csv(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            input_path = Path(tmp_dir) / "events.csv"
            output_dir = Path(tmp_dir) / "parts"
            input_path.write_text("id,name\n", encoding="utf-8")

            written_files = split_csv_file_by_row_count(str(input_path), str(output_dir), 2)

            self.assertEqual([], written_files)

    def test_raises_file_not_found_for_missing_input_file(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            output_dir = Path(tmp_dir) / "parts"

            with self.assertRaises(FileNotFoundError):
                split_csv_file_by_row_count(str(Path(tmp_dir) / "missing.csv"), str(output_dir), 2)

    def test_rejects_invalid_rows_per_file(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            input_path = Path(tmp_dir) / "events.csv"
            output_dir = Path(tmp_dir) / "parts"
            input_path.write_text("id,name\n1,alpha\n", encoding="utf-8")

            with self.assertRaises(ValueError):
                split_csv_file_by_row_count(str(input_path), str(output_dir), 0)

            with self.assertRaises(ValueError):
                split_csv_file_by_row_count(str(input_path), str(output_dir), -1)

            with self.assertRaises(ValueError):
                split_csv_file_by_row_count(str(input_path), str(output_dir), True)


if __name__ == "__main__":
    unittest.main()