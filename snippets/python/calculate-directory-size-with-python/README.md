# Calculate Directory Size with Python

Walk a directory tree and return total bytes, file count, and subdirectory count.

This snippet is useful when you need a quick directory inventory for build output checks, upload preparation, cleanup automation, or storage audits.

## Highlights

- Walks nested directories
- Adds file sizes
- Counts files and folders

## Use Cases

- Measure build artifact size before publishing
- Audit import or export folders before moving them elsewhere
- Track directory growth in local automation scripts

## Code

```python
from src.calculate_directory_size import calculate_directory_size


result = calculate_directory_size("build")
print(result["total_bytes"])
print(result["file_count"])
print(result["subdirectory_count"])
```

## Notes

- `total_bytes` sums file sizes only; directory metadata size is not included.
- `subdirectory_count` counts directories under the root directory and excludes the root itself.
- Hidden files and hidden directories are included when `os.walk` sees them.
- Pass a real directory path.

## Verification

Run the unit tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- flat directory totals
- nested directory totals
- empty directory handling
- zero-byte file counting
- non-directory path errors

## Files

- `src/calculate_directory_size.py`
- `tests/test_calculate_directory_size.py`
- `snippet.json`