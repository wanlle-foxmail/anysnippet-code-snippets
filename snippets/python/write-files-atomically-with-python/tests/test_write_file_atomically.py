import sys
import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch


SNIPPET_ROOT = Path(__file__).resolve().parent.parent
if str(SNIPPET_ROOT) not in sys.path:
    sys.path.insert(0, str(SNIPPET_ROOT))

from src.write_file_atomically import write_file_atomically


class WriteFileAtomicallyTests(unittest.TestCase):
    def test_writes_new_file(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            file_path = Path(tmp_dir) / "settings.json"

            write_file_atomically(str(file_path), '{"status": "ok"}')

            self.assertEqual('{"status": "ok"}', file_path.read_text(encoding="utf-8"))

    def test_replaces_existing_file_content(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            file_path = Path(tmp_dir) / "settings.json"
            file_path.write_text("old-value", encoding="utf-8")

            write_file_atomically(str(file_path), "new-value")

            self.assertEqual("new-value", file_path.read_text(encoding="utf-8"))

    def test_writes_empty_text(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            file_path = Path(tmp_dir) / "empty.txt"

            write_file_atomically(str(file_path), "")

            self.assertEqual("", file_path.read_text(encoding="utf-8"))

    def test_raises_file_not_found_when_parent_directory_is_missing(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            file_path = Path(tmp_dir) / "missing" / "settings.json"

            with self.assertRaises(FileNotFoundError):
                write_file_atomically(str(file_path), "content")

    def test_cleans_up_temp_file_and_keeps_original_content_when_replace_fails(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            file_path = Path(tmp_dir) / "settings.json"
            file_path.write_text("stable-content", encoding="utf-8")
            before_names = sorted(path.name for path in Path(tmp_dir).iterdir())

            with patch("src.write_file_atomically.os.replace", side_effect=OSError("replace failed")):
                with self.assertRaises(OSError):
                    write_file_atomically(str(file_path), "new-content")

            after_names = sorted(path.name for path in Path(tmp_dir).iterdir())
            self.assertEqual(before_names, after_names)
            self.assertEqual("stable-content", file_path.read_text(encoding="utf-8"))

    def test_writes_unicode_text(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            file_path = Path(tmp_dir) / "unicode.txt"

            write_file_atomically(str(file_path), "你好, AnySnippet")

            self.assertEqual("你好, AnySnippet", file_path.read_text(encoding="utf-8"))


if __name__ == "__main__":
    unittest.main()