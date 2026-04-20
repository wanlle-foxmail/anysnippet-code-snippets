import sys
import tempfile
import unittest
from pathlib import Path


SNIPPET_ROOT = Path(__file__).resolve().parent.parent
if str(SNIPPET_ROOT) not in sys.path:
    sys.path.insert(0, str(SNIPPET_ROOT))

from src.calculate_directory_size import calculate_directory_size


class CalculateDirectorySizeTests(unittest.TestCase):
    def test_calculates_total_bytes_for_flat_directory(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root_dir = Path(tmp_dir)
            (root_dir / "alpha.bin").write_bytes(b"abc")
            (root_dir / "beta.bin").write_bytes(b"12345")

            result = calculate_directory_size(str(root_dir))

            self.assertEqual(8, result["total_bytes"])
            self.assertEqual(2, result["file_count"])
            self.assertEqual(0, result["subdirectory_count"])

    def test_calculates_nested_directory_totals(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root_dir = Path(tmp_dir)
            (root_dir / "root.txt").write_bytes(b"root")
            first_nested = root_dir / "nested"
            second_nested = first_nested / "deeper"
            second_nested.mkdir(parents=True)
            (first_nested / "inner.txt").write_bytes(b"12")
            (second_nested / "deep.txt").write_bytes(b"123")

            result = calculate_directory_size(str(root_dir))

            self.assertEqual(9, result["total_bytes"])
            self.assertEqual(3, result["file_count"])
            self.assertEqual(2, result["subdirectory_count"])

    def test_returns_zeroes_for_empty_directory(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            result = calculate_directory_size(str(tmp_dir))

            self.assertEqual(0, result["total_bytes"])
            self.assertEqual(0, result["file_count"])
            self.assertEqual(0, result["subdirectory_count"])

    def test_counts_zero_byte_files(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root_dir = Path(tmp_dir)
            (root_dir / "empty.txt").write_bytes(b"")

            result = calculate_directory_size(str(root_dir))

            self.assertEqual(0, result["total_bytes"])
            self.assertEqual(1, result["file_count"])
            self.assertEqual(0, result["subdirectory_count"])

    def test_raises_value_error_for_non_directory_path(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            file_path = Path(tmp_dir) / "not-a-directory.txt"
            file_path.write_text("content", encoding="utf-8")

            with self.assertRaises(ValueError):
                calculate_directory_size(str(file_path))


if __name__ == "__main__":
    unittest.main()