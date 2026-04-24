# Deduplicate Large JSONL Records by Key with Python

Read a large JSONL file and write only the first record for each unique top-level string key.

This snippet is useful when event exports or crawler output contain duplicate records and you want to keep the first occurrence of each record ID or email.

## Highlights

- Keeps first record per key
- Preserves original order
- Reports bad input lines

## Use Cases

- Deduplicate exported records by ID
- Keep the first event per message key
- Clean JSONL input before downstream imports

## Code

```python
from src.deduplicate_jsonl_records_by_key import deduplicate_jsonl_records_by_key


written = deduplicate_jsonl_records_by_key("events.jsonl", "unique-events.jsonl", "id")
print(written)
```

## Notes

- The key must exist as a top-level field in every JSON object.
- The key value must be a string.
- Memory usage grows with the number of unique keys because seen keys are tracked in memory.

## Verification

Run the unit tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- first-record deduplication
- blank-line skipping
- Unicode preservation
- missing input files
- invalid JSON lines
- missing key errors
- invalid key-name or key-value errors

## Files

- `src/deduplicate_jsonl_records_by_key.py`
- `tests/test_deduplicate_jsonl_records_by_key.py`
- `snippet.json`