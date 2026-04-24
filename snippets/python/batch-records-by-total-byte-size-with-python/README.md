# Batch Records by Total Byte Size with Python

Group UTF-8 string records into batches whose total encoded byte size stays under a limit.

This snippet is useful when an API, queue, or storage service limits payload size by bytes instead of by record count.

## Highlights

- Batches by byte budget
- Preserves input order
- Handles Unicode byte sizes

## Use Cases

- Prepare API upload batches by payload budget
- Group queue messages under a byte cap
- Split text records for size-limited storage writes

## Code

```python
from src.batch_records_by_total_byte_size import batch_records_by_total_byte_size


result = list(batch_records_by_total_byte_size(["alpha", "beta", "gamma"], 9))
print(result)
```

## Notes

- This snippet expects an iterable of strings.
- Batch size is calculated from each record's UTF-8 encoded byte length.
- Separator or wrapper bytes are not included in the total.

## Verification

Run the unit tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- ASCII batching by byte count
- Unicode byte handling
- generator input ordering
- empty input handling
- oversized record errors
- invalid byte-limit errors
- non-string record errors

## Files

- `src/batch_records_by_total_byte_size.py`
- `tests/test_batch_records_by_total_byte_size.py`
- `snippet.json`