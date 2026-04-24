# Write JSON Lines with Python

Write Python values to a UTF-8 JSON Lines file and replace the output file atomically.

This snippet is useful when you need to export records for batch jobs, line-based logs, or append-friendly data pipelines that expect JSONL output.

## Highlights

- Writes files atomically
- Preserves readable Unicode
- Serializes one value per line

## Use Cases

- Export API results to a JSONL file
- Save batch job results for later processing
- Produce line-based input for downstream tools

## Code

```python
from src.write_json_lines import write_json_lines


write_json_lines("events.jsonl", [{"id": 1}, {"id": 2}])
print("events.jsonl")
```

## Notes

- The parent directory must already exist.
- Each call replaces the target file atomically.
- This snippet keeps Unicode readable with `ensure_ascii=False`.
- A serialization failure leaves the previous output file unchanged.

## Verification

Run the unit tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- ordered multi-record writes
- empty iterable output
- Unicode text preservation
- atomic replacement of an existing file
- missing parent directory errors
- serialization failure cleanup

## Files

- `src/write_json_lines.py`
- `tests/test_write_json_lines.py`
- `snippet.json`