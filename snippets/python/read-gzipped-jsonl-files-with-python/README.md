# Read Gzipped JSONL Files with Python

Read a gzipped JSON Lines file lazily and yield one parsed value per non-empty line.

This snippet is useful when logs, exports, or queue payloads are stored as `.jsonl.gz` files and you want to process them one record at a time without decompressing the whole file into memory.

## Highlights

- Reads `.jsonl.gz` lazily
- Skips blank lines safely
- Reports invalid line numbers

## Use Cases

- Process compressed event logs line by line
- Iterate through gzipped export records during imports
- Feed streaming ETL jobs from compressed JSONL input

## Code

```python
from src.read_gzipped_jsonl_file import read_gzipped_jsonl_file


for item in read_gzipped_jsonl_file("events.jsonl.gz"):
    print(item)
```

## Notes

- The input file must be a gzip-compressed UTF-8 JSONL file.
- Each non-empty line is decompressed and parsed with `json.loads`.
- Invalid JSON raises `ValueError` with the 1-based line number.

## Verification

Run the unit tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- ordered reads across multiple lines
- blank line skipping
- empty file handling
- invalid JSON line errors
- missing file behavior
- Unicode values

## Files

- `src/read_gzipped_jsonl_file.py`
- `tests/test_read_gzipped_jsonl_file.py`
- `snippet.json`