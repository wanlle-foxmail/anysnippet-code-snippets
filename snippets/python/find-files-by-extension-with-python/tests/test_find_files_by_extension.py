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

from src.find_files_by_extension import find_files_by_extension


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


class FindFilesByExtensionTests(unittest.TestCase):
    def test_finds_files_for_single_extension_with_counts(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root_dir = Path(tmp_dir) / "scan-root"
            root_dir.mkdir()
            matching_file = root_dir / "alpha.txt"
            matching_file.write_text("alpha", encoding="utf-8")
            (root_dir / "beta.md").write_text("beta", encoding="utf-8")
            (root_dir / "gamma.txt").write_text("gamma", encoding="utf-8")
            input_dir = root_dir.parent / root_dir.name / ".." / root_dir.name

            result = find_files_by_extension(input_dir, ".txt")

            self.assertEqual(os.path.abspath(str(input_dir)), result["root_directory"])
            self.assertEqual([".txt"], result["normalized_extensions"])
            self.assertEqual(3, result["total_files"])
            self.assertEqual(2, result["matched_files"])
            self.assertEqual(1, result["skipped_files"])
            self.assertEqual(os.path.abspath(str(matching_file)), result["matching_file_paths"][0])

    def test_finds_files_for_multiple_extensions_in_nested_directories(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root_dir = Path(tmp_dir)
            (root_dir / "keep.py").write_text("print('x')", encoding="utf-8")
            nested_dir = root_dir / "nested"
            nested_dir.mkdir()
            nested_match = nested_dir / "notes.txt"
            nested_match.write_text("notes", encoding="utf-8")
            (nested_dir / "skip.log").write_text("skip", encoding="utf-8")

            result = find_files_by_extension(root_dir, [".py", ".txt"])

            self.assertEqual([".py", ".txt"], result["normalized_extensions"])
            self.assertEqual(3, result["total_files"])
            self.assertEqual(2, result["matched_files"])
            self.assertEqual(1, result["skipped_files"])
            self.assertIn(os.path.abspath(str(nested_match)), result["matching_file_paths"])

    def test_normalizes_extensions_and_removes_duplicates(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root_dir = Path(tmp_dir)
            (root_dir / "first.txt").write_text("one", encoding="utf-8")
            (root_dir / "second.TXT").write_text("two", encoding="utf-8")

            result = find_files_by_extension(root_dir, ["txt", ".TXT", " txt "])

            self.assertEqual([".txt"], result["normalized_extensions"])
            self.assertEqual(2, result["matched_files"])
            self.assertEqual(0, result["skipped_files"])

    def test_matches_compound_extensions(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root_dir = Path(tmp_dir)
            archive_file = root_dir / "backup.tar.gz"
            archive_file.write_text("archive", encoding="utf-8")
            (root_dir / "backup.gz").write_text("skip", encoding="utf-8")

            result = find_files_by_extension(root_dir, ".tar.gz")

            self.assertEqual(2, result["total_files"])
            self.assertEqual(1, result["matched_files"])
            self.assertEqual([os.path.abspath(str(archive_file))], result["matching_file_paths"])

    def test_matches_extensions_case_insensitively(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root_dir = Path(tmp_dir)
            uppercase_file = root_dir / "README.TXT"
            uppercase_file.write_text("upper", encoding="utf-8")

            result = find_files_by_extension(root_dir, ".txt")

            self.assertEqual(1, result["matched_files"])
            self.assertEqual([os.path.abspath(str(uppercase_file))], result["matching_file_paths"])

    def test_returns_absolute_paths_in_deterministic_order(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            workspace_dir = Path(tmp_dir)
            root_dir = workspace_dir / "search-root"
            root_dir.mkdir()
            (root_dir / "c.txt").write_text("c", encoding="utf-8")
            first_file = root_dir / "a.txt"
            first_file.write_text("a", encoding="utf-8")
            external_file = workspace_dir / "outside.txt"
            external_file.write_text("outside", encoding="utf-8")
            nested_dir = root_dir / "nested"
            nested_dir.mkdir()
            nested_file = nested_dir / "b.txt"
            nested_file.write_text("b", encoding="utf-8")
            expected_paths = [
                os.path.abspath(str(first_file)),
                os.path.abspath(str(root_dir / "c.txt")),
                os.path.abspath(str(nested_file)),
            ]
            expected_total_files = 3

            try:
                symlink_file = root_dir / "link.txt"
                os.symlink(external_file, symlink_file)
                expected_paths.insert(2, os.path.abspath(str(symlink_file)))
                expected_total_files += 1

                external_dir = workspace_dir / "outside-dir"
                external_dir.mkdir()
                (external_dir / "ignored.txt").write_text("ignored", encoding="utf-8")
                os.symlink(external_dir, root_dir / "linked-dir")
            except (AttributeError, NotImplementedError):
                pass
            except OSError as error:
                if not _is_expected_symlink_error(error):
                    raise

            result = find_files_by_extension(root_dir, ".txt")

            self.assertEqual(expected_total_files, result["total_files"])
            self.assertEqual(expected_total_files, result["matched_files"])
            self.assertEqual(0, result["skipped_files"])
            self.assertEqual(expected_paths, result["matching_file_paths"])

    def test_includes_hidden_files_if_extension_matches(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            root_dir = Path(tmp_dir)
            hidden_file = root_dir / ".hidden.txt"
            hidden_file.write_text("hidden", encoding="utf-8")

            result = find_files_by_extension(root_dir, ".txt")

            self.assertEqual(1, result["total_files"])
            self.assertEqual(1, result["matched_files"])
            self.assertEqual([os.path.abspath(str(hidden_file))], result["matching_file_paths"])

    def test_returns_zero_matches_for_empty_directory(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            result = find_files_by_extension(tmp_dir, ".txt")

            self.assertEqual(0, result["total_files"])
            self.assertEqual(0, result["matched_files"])
            self.assertEqual(0, result["skipped_files"])
            self.assertEqual([], result["matching_file_paths"])

    def test_raises_file_not_found_for_missing_directory(self):
        with self.assertRaises(FileNotFoundError):
            find_files_by_extension("missing-find-files-by-extension", ".txt")

    def test_raises_value_error_for_file_path(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            file_path = Path(tmp_dir) / "not-a-directory.txt"
            file_path.write_text("content", encoding="utf-8")

            with self.assertRaises(ValueError):
                find_files_by_extension(file_path, ".txt")

    def test_raises_value_error_for_blank_extensions(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            with self.assertRaises(ValueError):
                find_files_by_extension(tmp_dir, ["   "])

    def test_propagates_walk_errors(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            def failing_walk(*args, **kwargs):
                kwargs["onerror"](PermissionError("denied"))
                return iter(())

            with mock.patch("src.find_files_by_extension.os.walk", side_effect=failing_walk):
                with self.assertRaises(PermissionError):
                    find_files_by_extension(tmp_dir, ".txt")

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
                find_files_by_extension(symlink_dir, ".txt")


if __name__ == "__main__":
    unittest.main()