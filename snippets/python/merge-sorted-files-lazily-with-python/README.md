# Merge Sorted Files Lazily with Python

Lazily merge multiple sorted UTF-8 text files into one sorted output stream.

This snippet is useful when large text files have already been sorted and you want to combine them without loading every line into memory.

## Highlights

- Merges files lazily
- Preserves duplicate lines
- Uses heap-based ordering

## Use Cases

- Merge sorted log exports
- Combine sorted batch output files
- Stream one merged view from sorted text parts

## Code

```python
from src.merge_sorted_files import merge_sorted_text_files


for item in merge_sorted_text_files(["part-1.txt", "part-2.txt"]):
    print(item)
```

## Notes

- Each input file must already be sorted lexicographically.
- Output lines are yielded without trailing newline characters.
- Duplicate lines are preserved.

## Verification

Run the unit tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- two-file merges
- duplicate-line preservation
- empty-file handling
- empty path lists
- missing file behavior
- Unicode lines

## Files

- `src/merge_sorted_files.py`
- `tests/test_merge_sorted_files.py`
- `snippet.json`