# Read Large JSONL Files with Python

Read a JSON Lines file lazily and yield one parsed value per non-empty line.

This snippet is useful when logs, exports, or queue payloads are stored as JSONL and you want to process them one record at a time without loading the whole file into memory.

## Highlights

- Yields one JSON value per line
- Skips blank lines safely
- Reports invalid line numbers

## Use Cases

- Process large event logs line by line
- Iterate through exported records during imports
- Feed streaming ETL jobs with JSONL input

## Code

```python
from src.read_jsonl_file import read_jsonl_file


for item in read_jsonl_file("events.jsonl"):
    print(item)
```

## Notes

- Each non-empty line is parsed with `json.loads`.
- Blank lines are ignored.
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

- `src/read_jsonl_file.py`
- `tests/test_read_jsonl_file.py`
- `snippet.json`