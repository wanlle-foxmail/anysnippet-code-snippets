import errno
import os
import sys
import tempfile
import unittest
from unittest import mock
from pathlib import Path


SNIPPET_ROOT = Path(__file__).resolve().parent.parent
if str(SNIPPET_ROOT) not in sys.path:
    sys.path.insert(0, str(SNIPPET_ROOT))

from src.calculate_directory_size import calculate_directory_size


def _is_expected_symlink_error(error):
    if isinstance(error, PermissionError):
        return True
    if getattr(error, "winerror", None) == 1314:
        return True
    return error.errno in {
        errno.EACCES,
        errno.EPERM,
        getattr(errno, "ENOTSUP", errno.EPERM),
        getattr(errno, "EOPNOTSUPP", errno.EPERM),
    }


class CalculateDirectorySizeTests(unittest.TestCase):
    def test_calculates_total_bytes_and_file_count_in_flat_directory(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root_dir = Path(tmp_dir)
            (root_dir / "alpha.bin").write_bytes(b"abc")
            (root_dir / "beta.bin").write_bytes(b"12345")

            result = calculate_directory_size(root_dir)

            self.assertEqual(8, result["total_bytes"])
            self.assertEqual(2, result["file_count"])
            self.assertEqual(0, result["subdirectory_count"])

    def test_calculates_total_bytes_recursively(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root_dir = Path(tmp_dir)
            (root_dir / "root.txt").write_bytes(b"root")
            first_nested = root_dir / "nested"
            second_nested = first_nested / "deeper"
            second_nested.mkdir(parents=True)
            (first_nested / "inner.txt").write_bytes(b"12")
            (second_nested / "deep.txt").write_bytes(b"123")

            result = calculate_directory_size(root_dir)

            self.assertEqual(9, result["total_bytes"])
            self.assertEqual(3, result["file_count"])
            self.assertEqual(2, result["subdirectory_count"])

    def test_returns_zeroes_for_empty_directory(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            result = calculate_directory_size(tmp_dir)

            self.assertEqual(0, result["total_bytes"])
            self.assertEqual(0, result["file_count"])
            self.assertEqual(0, result["subdirectory_count"])

    def test_counts_zero_byte_files(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root_dir = Path(tmp_dir)
            (root_dir / "empty.txt").write_bytes(b"")

            result = calculate_directory_size(root_dir)

            self.assertEqual(0, result["total_bytes"])
            self.assertEqual(1, result["file_count"])
            self.assertEqual(0, result["subdirectory_count"])

    def test_includes_hidden_files_and_hidden_directories_but_skips_symlinks(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            workspace_dir = Path(tmp_dir)
            root_dir = workspace_dir / "scan-root"
            root_dir.mkdir()
            (root_dir / ".hidden.txt").write_bytes(b"1234")
            hidden_dir = root_dir / ".cache"
            hidden_dir.mkdir()
            (hidden_dir / "cached.bin").write_bytes(b"abcdef")

            try:
                external_file = workspace_dir / "outside.bin"
                external_file.write_bytes(b"outside")
                external_dir = workspace_dir / "outside-dir"
                external_dir.mkdir()
                (external_dir / "external.bin").write_bytes(b"external")
                os.symlink(external_file, root_dir / "outside-link.bin")
                os.symlink(external_dir, root_dir / "outside-link-dir")
            except (AttributeError, NotImplementedError):
                pass
            except OSError as error:
                if not _is_expected_symlink_error(error):
                    raise

            result = calculate_directory_size(root_dir)

            self.assertEqual(10, result["total_bytes"])
            self.assertEqual(2, result["file_count"])
            self.assertEqual(1, result["subdirectory_count"])

    def test_subdirectory_count_excludes_root_directory(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root_dir = Path(tmp_dir)
            (root_dir / "first").mkdir()
            (root_dir / "second").mkdir()

            result = calculate_directory_size(root_dir)

            self.assertEqual(2, result["subdirectory_count"])

    def test_returns_absolute_root_directory(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root_dir = Path(tmp_dir) / "scan-root"
            root_dir.mkdir()
            input_dir = root_dir.parent / root_dir.name / ".." / root_dir.name
            result = calculate_directory_size(input_dir)

            self.assertEqual(os.path.abspath(str(input_dir)), result["root_directory"])

    def test_raises_file_not_found_for_missing_directory(self):
        with self.assertRaises(FileNotFoundError):
            calculate_directory_size("missing-calculate-directory-size")

    def test_raises_value_error_for_file_path(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            file_path = Path(tmp_dir) / "not-a-directory.txt"
            file_path.write_text("content", encoding="utf-8")

            with self.assertRaises(ValueError):
                calculate_directory_size(file_path)

    def test_propagates_walk_errors(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            def failing_walk(*args, **kwargs):
                kwargs["onerror"](PermissionError("denied"))
                return iter(())

            with mock.patch("src.calculate_directory_size.os.walk", side_effect=failing_walk):
                with self.assertRaises(PermissionError):
                    calculate_directory_size(tmp_dir)

    def test_raises_value_error_for_symlink_root(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            external_dir = Path(tmp_dir) / "outside-dir"
            external_dir.mkdir()

            try:
                symlink_dir = Path(tmp_dir) / "symlink-root"
                os.symlink(external_dir, symlink_dir)
            except (AttributeError, NotImplementedError):
                return
            except OSError as error:
                if _is_expected_symlink_error(error):
                    return
                raise

            with self.assertRaises(ValueError):
                calculate_directory_size(symlink_dir)


if __name__ == "__main__":
    unittest.main()