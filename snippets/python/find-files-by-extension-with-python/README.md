# Find Files by Extension with Python

Walk through a directory and collect files whose extensions match a target list.

This snippet is useful when you need to build a file inventory for automation, packaging, static checks, or batch processing tasks.

## Highlights

- Walks nested directories
- Counts scanned and matched files
- Returns matched file paths

## Use Cases

- Collect all source files before running a code generation step
- Find documentation or asset files during build automation
- Create an input list for batch processing by file type

## Code

```python
from src.find_files_by_extension import find_files_by_extension


result = find_files_by_extension("project", [".py", ".md"])
print(result["hit"])
print(result["items"])
```

## Notes

- `count` is the number of files scanned.
- `hit` is the number of matched files.
- `items` contains the matched file paths.
- Put a leading dot in each extension, such as `.csv` or `.txt`.

## Verification

Run the unit tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- single-extension filtering
- nested directory traversal
- multiple extension matching
- no-match results
- empty directory handling

## Files

- `src/find_files_by_extension.py`
- `tests/test_find_files_by_extension.py`
- `snippet.json`