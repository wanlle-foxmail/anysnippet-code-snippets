import json
import sys
import tempfile
import unittest
from pathlib import Path


SNIPPET_ROOT = Path(__file__).resolve().parent.parent
if str(SNIPPET_ROOT) not in sys.path:
    sys.path.insert(0, str(SNIPPET_ROOT))

from src.deduplicate_jsonl_records_by_key import deduplicate_jsonl_records_by_key


def read_jsonl(path: Path) -> list[dict[str, object]]:
    with open(path, "r", encoding="utf-8") as input_file:
        return [json.loads(line) for line in input_file if line.strip()]


class DeduplicateJsonlRecordsByKeyTests(unittest.TestCase):
    def test_deduplicates_by_key_and_keeps_first_record(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            input_path = Path(tmp_dir) / "events.jsonl"
            output_path = Path(tmp_dir) / "unique-events.jsonl"
            input_path.write_text(
                '{"id": "a", "value": 1}\n'
                '{"id": "b", "value": 2}\n'
                '{"id": "a", "value": 3}\n'
                '{"id": "c", "value": 4}\n',
                encoding="utf-8",
            )

            written = deduplicate_jsonl_records_by_key(str(input_path), str(output_path), "id")

            self.assertEqual(3, written)
            self.assertEqual(
                [
                    {"id": "a", "value": 1},
                    {"id": "b", "value": 2},
                    {"id": "c", "value": 4},
                ],
                read_jsonl(output_path),
            )

    def test_skips_blank_lines_and_preserves_unicode(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            input_path = Path(tmp_dir) / "events.jsonl"
            output_path = Path(tmp_dir) / "unique-events.jsonl"
            input_path.write_text(
                '\n{"id": "x", "message": "你好"}\n\n{"id": "x", "message": "世界"}\n',
                encoding="utf-8",
            )

            written = deduplicate_jsonl_records_by_key(str(input_path), str(output_path), "id")

            self.assertEqual(1, written)
            self.assertEqual([{"id": "x", "message": "你好"}], read_jsonl(output_path))

    def test_raises_file_not_found_for_missing_input_file(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            output_path = Path(tmp_dir) / "unique-events.jsonl"

            with self.assertRaises(FileNotFoundError):
                deduplicate_jsonl_records_by_key("missing.jsonl", str(output_path), "id")

    def test_raises_value_error_for_invalid_json_line(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            input_path = Path(tmp_dir) / "events.jsonl"
            output_path = Path(tmp_dir) / "unique-events.jsonl"
            input_path.write_text('{"id": "a"}\nnot-json\n', encoding="utf-8")

            with self.assertRaises(ValueError) as raised:
                deduplicate_jsonl_records_by_key(str(input_path), str(output_path), "id")

            self.assertIn("line 2", str(raised.exception))

    def test_raises_value_error_for_missing_key(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            input_path = Path(tmp_dir) / "events.jsonl"
            output_path = Path(tmp_dir) / "unique-events.jsonl"
            input_path.write_text('{"name": "a"}\n', encoding="utf-8")

            with self.assertRaises(ValueError) as raised:
                deduplicate_jsonl_records_by_key(str(input_path), str(output_path), "id")

            self.assertIn("missing key 'id'", str(raised.exception))

    def test_raises_value_error_for_non_string_key_value(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            input_path = Path(tmp_dir) / "events.jsonl"
            output_path = Path(tmp_dir) / "unique-events.jsonl"
            input_path.write_text('{"id": 1}\n', encoding="utf-8")

            with self.assertRaises(ValueError) as raised:
                deduplicate_jsonl_records_by_key(str(input_path), str(output_path), "id")

            self.assertIn("must be a string", str(raised.exception))

    def test_rejects_invalid_key_name(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            input_path = Path(tmp_dir) / "events.jsonl"
            output_path = Path(tmp_dir) / "unique-events.jsonl"
            input_path.write_text('{"id": "a"}\n', encoding="utf-8")

            with self.assertRaises(TypeError):
                deduplicate_jsonl_records_by_key(str(input_path), str(output_path), None)

            with self.assertRaises(ValueError):
                deduplicate_jsonl_records_by_key(str(input_path), str(output_path), "   ")


if __name__ == "__main__":
    unittest.main()