import sys
import tempfile
import unittest
from pathlib import Path


SNIPPET_ROOT = Path(__file__).resolve().parent.parent
if str(SNIPPET_ROOT) not in sys.path:
    sys.path.insert(0, str(SNIPPET_ROOT))

from src.read_jsonl_file import read_jsonl_file


class ReadJsonlFileTests(unittest.TestCase):
    def test_reads_multiple_json_objects_in_order(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            file_path = Path(tmp_dir) / "events.jsonl"
            file_path.write_text('{"id": 1}\n{"id": 2}\n', encoding="utf-8")

            items = list(read_jsonl_file(str(file_path)))

            self.assertEqual([{"id": 1}, {"id": 2}], items)

    def test_skips_blank_lines(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            file_path = Path(tmp_dir) / "events.jsonl"
            file_path.write_text('{"id": 1}\n\n{"id": 2}\n   \n', encoding="utf-8")

            items = list(read_jsonl_file(str(file_path)))

            self.assertEqual([{"id": 1}, {"id": 2}], items)

    def test_returns_no_items_for_empty_file(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            file_path = Path(tmp_dir) / "events.jsonl"
            file_path.write_text("", encoding="utf-8")

            items = list(read_jsonl_file(str(file_path)))

            self.assertEqual([], items)

    def test_raises_value_error_with_line_number_for_invalid_json(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            file_path = Path(tmp_dir) / "events.jsonl"
            file_path.write_text('{"id": 1}\nnot-json\n', encoding="utf-8")

            with self.assertRaises(ValueError) as raised:
                list(read_jsonl_file(str(file_path)))

            self.assertIn("line 2", str(raised.exception))

    def test_raises_file_not_found_for_missing_file(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            file_path = Path(tmp_dir) / "missing.jsonl"

            with self.assertRaises(FileNotFoundError):
                list(read_jsonl_file(str(file_path)))

    def test_reads_unicode_values(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            file_path = Path(tmp_dir) / "events.jsonl"
            file_path.write_text('{"message": "你好"}\n', encoding="utf-8")

            items = list(read_jsonl_file(str(file_path)))

            self.assertEqual([{"message": "你好"}], items)


if __name__ == "__main__":
    unittest.main()