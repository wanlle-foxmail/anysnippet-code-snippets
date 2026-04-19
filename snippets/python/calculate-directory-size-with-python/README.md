# Calculate Directory Size with Python

Recursively calculate total file size and count how many files and subdirectories exist under a target directory.

This snippet is useful when you need a quick directory inventory for build output checks, upload preparation, cleanup automation, or storage audits.

## Highlights

- Traverses nested directories with `os.walk`
- Sums total file bytes across the full directory tree
- Reports file and subdirectory counts separately
- Includes hidden files and hidden directories in the totals

## Use Cases

- Measure build artifact size before publishing
- Audit import or export folders before moving them elsewhere
- Track directory growth in local automation scripts

## Code

```python
from pathlib import Path

from src.calculate_directory_size import calculate_directory_size


result = calculate_directory_size(Path("build"))
print(result["total_bytes"])
print(result["file_count"])
```

## Notes

- `total_bytes` sums file sizes only; directory metadata size is not included.
- `subdirectory_count` counts directories under the root directory and excludes the root itself.
- Hidden files and hidden directories are included in the reported totals.
- The root path must be a real directory; nested symlinked files and directories are skipped.
- Symlinked files and symlinked directories are skipped so the result stays scoped to the scanned tree.
- Traversal errors such as unreadable descendant directories are raised instead of being silently ignored.

## Verification

Run the unit tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- flat directory size calculation
- recursive size accumulation across nested directories
- empty directory handling
- zero-byte file counting
- hidden file and hidden directory inclusion
- subdirectory counting that excludes the root directory
- absolute root directory reporting
- missing directory errors
- non-directory path errors

## Files

- `src/calculate_directory_size.py`
- `tests/test_calculate_directory_size.py`
- `snippet.json`