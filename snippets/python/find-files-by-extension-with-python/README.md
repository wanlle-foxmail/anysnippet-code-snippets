# Find Files by Extension with Python

Recursively traverse a directory, filter files by one or more extensions, and return absolute paths with traversal counts.

This snippet is useful when you need to build a file inventory for automation, packaging, static checks, or batch processing tasks.

## Highlights

- Traverses nested directories with `os.walk`
- Accepts one or many extensions and normalizes them automatically
- Matches file suffixes case-insensitively, including compound suffixes such as `.tar.gz`
- Returns total, matched, and skipped file counts in one dictionary

## Use Cases

- Collect all source files before running a code generation step
- Find documentation or asset files during build automation
- Create an input list for batch processing by file type

## Code

```python
from pathlib import Path

from src.find_files_by_extension import find_files_by_extension


result = find_files_by_extension(Path("project"), ["py", ".md"])
print(result["matched_files"])
print(result["matching_file_paths"])
```

## Notes

- Extension inputs are normalized, so `txt` becomes `.txt`.
- The root path must be a real directory and cannot be a symlinked directory.
- Matching is case-insensitive, so `.txt` also matches filenames such as `README.TXT`.
- Returned file paths are normalized absolute paths and follow deterministic traversal order based on sorted directory and file names.
- Symlinked files are reported by the path found inside the scanned directory, not by the resolved target path.
- Symlinked directories are not traversed.
- Traversal errors such as unreadable descendant directories are raised instead of being silently ignored.

## Verification

Run the unit tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- single-extension filtering with counts
- nested directory traversal with multiple extensions
- extension normalization and deduplication
- compound suffix matching such as `.tar.gz`
- case-insensitive suffix matching
- deterministic ordering of returned paths
- hidden file inclusion
- empty directory handling
- missing directory errors
- invalid directory or extension input

## Files

- `src/find_files_by_extension.py`
- `tests/test_find_files_by_extension.py`
- `snippet.json`