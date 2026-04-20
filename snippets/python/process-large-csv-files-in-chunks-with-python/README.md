# Process Large CSV Files in Chunks with pandas

Read a large CSV with pandas in chunks and return one processed result per chunk.

This snippet is useful when a CSV file is too large to fit comfortably in memory and you want to handle one chunk at a time.

## Highlights

- Processes one chunk at a time
- Keeps chunk result order
- Uses pandas `read_csv`

## Use Cases

- Aggregate metrics from very large CSV exports
- Clean or validate a large CSV before loading it elsewhere
- Apply chunk-level business logic without building one huge DataFrame

## Code

```python
from src.process_csv_in_chunks import process_csv_in_chunks

chunk_totals = process_csv_in_chunks(
    "sales.csv",
    lambda chunk: int(chunk["amount"].sum()),
    chunk_size=5000,
)

print(chunk_totals)
print(sum(chunk_totals))
```

## Notes

- `chunk_processor` receives one pandas DataFrame at a time.
- `chunk_size` controls how many rows pandas reads per chunk.
- Missing files and CSV parsing errors come from pandas.

## Verification

Install the dependency and run the unit tests from the snippet root:

```bash
python -m pip install pandas
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- single-chunk processing for small CSV files
- stable processing order across multiple chunks
- header-only CSV handling
- missing file behavior
- invalid `chunk_size` values

## Files

- `src/process_csv_in_chunks.py`
- `tests/test_process_csv_in_chunks.py`
- `snippet.json`