import os
import sys
import tempfile
import unittest
from pathlib import Path


SNIPPET_ROOT = Path(__file__).resolve().parent.parent
if str(SNIPPET_ROOT) not in sys.path:
    sys.path.insert(0, str(SNIPPET_ROOT))

from src.find_files_by_extension import find_files_by_extension


class FindFilesByExtensionTests(unittest.TestCase):
    def test_finds_files_for_single_extension(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root_dir = Path(tmp_dir)
            first_match = root_dir / "alpha.txt"
            first_match.write_text("alpha", encoding="utf-8")
            (root_dir / "skip.md").write_text("skip", encoding="utf-8")

            result = find_files_by_extension(str(root_dir), [".txt"])

            self.assertEqual(2, result["count"])
            self.assertEqual(1, result["hit"])
            self.assertEqual([os.path.join(str(root_dir), "alpha.txt")], result["items"])

    def test_finds_files_in_nested_directories(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root_dir = Path(tmp_dir)
            nested_dir = root_dir / "nested"
            nested_dir.mkdir()
            nested_match = nested_dir / "beta.txt"
            nested_match.write_text("beta", encoding="utf-8")
            (root_dir / "skip.md").write_text("skip", encoding="utf-8")

            result = find_files_by_extension(str(root_dir), [".txt"])

            self.assertEqual(2, result["count"])
            self.assertEqual(1, result["hit"])
            self.assertEqual([os.path.join(str(nested_dir), "beta.txt")], result["items"])

    def test_supports_multiple_extensions(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root_dir = Path(tmp_dir)
            (root_dir / "main.py").write_text("print('x')", encoding="utf-8")
            (root_dir / "README.md").write_text("notes", encoding="utf-8")
            (root_dir / "skip.txt").write_text("skip", encoding="utf-8")

            result = find_files_by_extension(str(root_dir), [".py", ".md"])

            self.assertEqual(3, result["count"])
            self.assertEqual(2, result["hit"])

    def test_returns_empty_items_when_nothing_matches(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root_dir = Path(tmp_dir)
            (root_dir / "notes.md").write_text("notes", encoding="utf-8")
            (root_dir / "image.png").write_text("png", encoding="utf-8")

            result = find_files_by_extension(str(root_dir), [".txt"])

            self.assertEqual(2, result["count"])
            self.assertEqual(0, result["hit"])
            self.assertEqual([], result["items"])

    def test_returns_zero_counts_for_empty_directory(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            result = find_files_by_extension(tmp_dir, [".txt"])

            self.assertEqual(0, result["count"])
            self.assertEqual(0, result["hit"])
            self.assertEqual([], result["items"])


if __name__ == "__main__":
    unittest.main()