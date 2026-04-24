import sys
import tempfile
import unittest
from pathlib import Path

import ijson


SNIPPET_ROOT = Path(__file__).resolve().parent.parent
if str(SNIPPET_ROOT) not in sys.path:
    sys.path.insert(0, str(SNIPPET_ROOT))

from src.read_large_json_array import read_large_json_array


class ReadLargeJsonArrayTests(unittest.TestCase):
    def test_reads_multiple_items_in_order(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            file_path = Path(tmp_dir) / "items.json"
            file_path.write_text(
                '[{"question": "q1", "answer": "a1"}, {"question": "q2", "answer": "a2"}]',
                encoding="utf-8",
            )

            items = list(read_large_json_array(str(file_path)))

            self.assertEqual(
                [
                    {"question": "q1", "answer": "a1"},
                    {"question": "q2", "answer": "a2"},
                ],
                items,
            )

    def test_preserves_nested_objects_and_lists(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            file_path = Path(tmp_dir) / "nested.json"
            file_path.write_text(
                '[{"meta": {"source": "crawler"}, "tags": ["a", "b"]}]',
                encoding="utf-8",
            )

            items = list(read_large_json_array(str(file_path)))

            self.assertEqual(
                [{"meta": {"source": "crawler"}, "tags": ["a", "b"]}],
                items,
            )

    def test_returns_no_items_for_empty_array(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            file_path = Path(tmp_dir) / "empty.json"
            file_path.write_text("[]", encoding="utf-8")

            items = list(read_large_json_array(str(file_path)))

            self.assertEqual([], items)

    def test_reads_unicode_values(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            file_path = Path(tmp_dir) / "unicode.json"
            file_path.write_text(
                '[{"question": "你好", "answer": "🚀"}]',
                encoding="utf-8",
            )

            items = list(read_large_json_array(str(file_path)))

            self.assertEqual([{"question": "你好", "answer": "🚀"}], items)

    def test_raises_file_not_found_for_missing_file(self):
        with self.assertRaises(FileNotFoundError):
            list(read_large_json_array("missing.json"))

    def test_raises_json_error_for_invalid_array_content(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            file_path = Path(tmp_dir) / "broken.json"
            file_path.write_text('[{"question": "q1"},', encoding="utf-8")

            with self.assertRaises(ijson.JSONError):
                list(read_large_json_array(str(file_path)))


if __name__ == "__main__":
    unittest.main()