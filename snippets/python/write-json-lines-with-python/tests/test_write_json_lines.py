import json
import sys
import tempfile
import unittest
from pathlib import Path


SNIPPET_ROOT = Path(__file__).resolve().parent.parent
if str(SNIPPET_ROOT) not in sys.path:
    sys.path.insert(0, str(SNIPPET_ROOT))

from src.write_json_lines import write_json_lines


class WriteJsonLinesTests(unittest.TestCase):
    def test_writes_multiple_records_in_order(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            file_path = Path(tmp_dir) / "events.jsonl"
            records = [{"id": 1}, {"id": 2, "active": True}]

            write_json_lines(str(file_path), records)

            expected_lines = [json.dumps(record, ensure_ascii=False) for record in records]
            self.assertEqual(expected_lines, file_path.read_text(encoding="utf-8").splitlines())

    def test_writes_empty_iterable_to_empty_file(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            file_path = Path(tmp_dir) / "events.jsonl"

            write_json_lines(str(file_path), [])

            self.assertEqual("", file_path.read_text(encoding="utf-8"))

    def test_preserves_unicode_characters(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            file_path = Path(tmp_dir) / "events.jsonl"

            write_json_lines(str(file_path), [{"message": "你好"}])

            self.assertIn("你好", file_path.read_text(encoding="utf-8"))

    def test_overwrites_existing_file(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            file_path = Path(tmp_dir) / "events.jsonl"
            file_path.write_text("stale-data\n", encoding="utf-8")

            write_json_lines(str(file_path), [{"id": 99}])

            self.assertEqual('{"id": 99}\n', file_path.read_text(encoding="utf-8"))

    def test_raises_file_not_found_when_parent_directory_is_missing(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            file_path = Path(tmp_dir) / "missing" / "events.jsonl"

            with self.assertRaises(FileNotFoundError):
                write_json_lines(str(file_path), [{"id": 1}])

    def test_raises_type_error_for_non_serializable_item(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            file_path = Path(tmp_dir) / "events.jsonl"
            file_path.write_text('{"id": 1}\n', encoding="utf-8")
            before_names = sorted(path.name for path in Path(tmp_dir).iterdir())

            with self.assertRaises(TypeError):
                write_json_lines(str(file_path), [{"bad": {1, 2, 3}}])

            after_names = sorted(path.name for path in Path(tmp_dir).iterdir())
            self.assertEqual(before_names, after_names)
            self.assertEqual('{"id": 1}\n', file_path.read_text(encoding="utf-8"))


if __name__ == "__main__":
    unittest.main()