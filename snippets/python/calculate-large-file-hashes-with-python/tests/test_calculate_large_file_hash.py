import hashlib
import sys
import tempfile
import unittest
from pathlib import Path


SNIPPET_ROOT = Path(__file__).resolve().parent.parent
if str(SNIPPET_ROOT) not in sys.path:
    sys.path.insert(0, str(SNIPPET_ROOT))

from src.calculate_large_file_hash import calculate_large_file_hash


class CalculateLargeFileHashTests(unittest.TestCase):
    def test_returns_sha256_digest_for_small_file(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            file_path = Path(tmp_dir) / "example.txt"
            content = b"hello hashing"
            file_path.write_bytes(content)

            digest = calculate_large_file_hash(file_path)

            self.assertEqual(hashlib.sha256(content).hexdigest(), digest)

    def test_returns_sha256_digest_for_multi_chunk_file(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            file_path = Path(tmp_dir) / "large.bin"
            content = (b"chunk-data-" * 2048) + b"tail"
            file_path.write_bytes(content)

            digest = calculate_large_file_hash(file_path, chunk_size=128)

            self.assertEqual(hashlib.sha256(content).hexdigest(), digest)

    def test_returns_digest_for_empty_file(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            file_path = Path(tmp_dir) / "empty.txt"
            file_path.write_bytes(b"")

            digest = calculate_large_file_hash(file_path)

            self.assertEqual(hashlib.sha256(b"").hexdigest(), digest)

    def test_supports_md5_algorithm(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            file_path = Path(tmp_dir) / "checksum.txt"
            content = b"compatibility-check"
            file_path.write_bytes(content)

            digest = calculate_large_file_hash(file_path, algorithm="md5")

            self.assertEqual(hashlib.md5(content).hexdigest(), digest)

    def test_accepts_string_path_for_unicode_filename(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            file_path = Path(tmp_dir) / "space and unicode café.txt"
            content = b"path handling"
            file_path.write_bytes(content)

            digest = calculate_large_file_hash(str(file_path))

            self.assertEqual(hashlib.sha256(content).hexdigest(), digest)

    def test_raises_value_error_for_directory_path(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            directory_path = Path(tmp_dir) / "folder"
            directory_path.mkdir()

            with self.assertRaises(ValueError):
                calculate_large_file_hash(directory_path)

    def test_raises_file_not_found_for_missing_file(self):
        with self.assertRaises(FileNotFoundError):
            calculate_large_file_hash("missing-file.bin")

    def test_raises_value_error_for_invalid_algorithm(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            file_path = Path(tmp_dir) / "example.txt"
            file_path.write_bytes(b"content")

            with self.assertRaises(ValueError):
                calculate_large_file_hash(file_path, algorithm="sha1")

            with self.assertRaises(ValueError):
                calculate_large_file_hash(file_path, algorithm=["sha1"])

    def test_raises_value_error_for_invalid_chunk_size(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            file_path = Path(tmp_dir) / "example.txt"
            file_path.write_bytes(b"content")

            with self.assertRaises(ValueError):
                calculate_large_file_hash(file_path, chunk_size=0)

            with self.assertRaises(ValueError):
                calculate_large_file_hash(file_path, chunk_size=-1)


if __name__ == "__main__":
    unittest.main()