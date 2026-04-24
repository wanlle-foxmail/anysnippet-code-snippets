import sys
import tempfile
import unittest
from pathlib import Path


SNIPPET_ROOT = Path(__file__).resolve().parent.parent
if str(SNIPPET_ROOT) not in sys.path:
    sys.path.insert(0, str(SNIPPET_ROOT))

from src.merge_sorted_files import merge_sorted_text_files


class MergeSortedFilesTests(unittest.TestCase):
    def test_merges_two_sorted_files_in_order(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            first_path = Path(tmp_dir) / "part-1.txt"
            second_path = Path(tmp_dir) / "part-2.txt"
            first_path.write_text("apple\npear\n", encoding="utf-8")
            second_path.write_text("banana\nplum\n", encoding="utf-8")

            merged_lines = list(merge_sorted_text_files([str(first_path), str(second_path)]))

            self.assertEqual(["apple", "banana", "pear", "plum"], merged_lines)

    def test_preserves_duplicate_lines(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            first_path = Path(tmp_dir) / "part-1.txt"
            second_path = Path(tmp_dir) / "part-2.txt"
            first_path.write_text("apple\npear\n", encoding="utf-8")
            second_path.write_text("apple\nplum\n", encoding="utf-8")

            merged_lines = list(merge_sorted_text_files([str(first_path), str(second_path)]))

            self.assertEqual(["apple", "apple", "pear", "plum"], merged_lines)

    def test_handles_empty_files(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            empty_path = Path(tmp_dir) / "empty.txt"
            data_path = Path(tmp_dir) / "data.txt"
            empty_path.write_text("", encoding="utf-8")
            data_path.write_text("alpha\nomega\n", encoding="utf-8")

            merged_lines = list(merge_sorted_text_files([str(empty_path), str(data_path)]))

            self.assertEqual(["alpha", "omega"], merged_lines)

    def test_returns_no_lines_for_empty_path_list(self):
        merged_lines = list(merge_sorted_text_files([]))

        self.assertEqual([], merged_lines)

    def test_raises_file_not_found_for_missing_file(self):
        with self.assertRaises(FileNotFoundError):
            list(merge_sorted_text_files(["missing.txt"]))

    def test_reads_unicode_lines(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            first_path = Path(tmp_dir) / "part-1.txt"
            second_path = Path(tmp_dir) / "part-2.txt"
            first_path.write_text("你好\n再见\n", encoding="utf-8")
            second_path.write_text("世界\n", encoding="utf-8")

            merged_lines = list(merge_sorted_text_files([str(first_path), str(second_path)]))

            self.assertEqual(["世界", "你好", "再见"], merged_lines)


if __name__ == "__main__":
    unittest.main()