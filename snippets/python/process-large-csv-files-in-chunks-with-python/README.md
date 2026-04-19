# Process Large CSV Files in Chunks with pandas

Read large CSV files in pandas chunks and process each chunk without loading the full dataset into memory.

This snippet is useful when a CSV file is too large to fit comfortably in memory, but you still want to apply pandas transformations, aggregations, or validations chunk by chunk.

## Highlights

- Uses `pandas.read_csv(..., chunksize=...)` for streaming-style processing
- Preserves chunk order in the returned results
- Passes `read_csv_kwargs` through for options like `usecols` and `dtype`
- Skips empty chunks so header-only CSV files return an empty result list

## Use Cases

- Aggregate metrics from very large CSV exports
- Clean or validate a large CSV before loading it into another system
- Process selected columns with `usecols` to reduce memory use
- Apply chunk-level business rules without building one huge DataFrame

## Code

```python
from src.process_csv_in_chunks import process_csv_in_chunks


chunk_totals = process_csv_in_chunks(
    "sales.csv",
    lambda chunk: int(chunk["amount"].sum()),
    chunk_size=5000,
    read_csv_kwargs={"usecols": ["customer_id", "amount"]},
)

print(chunk_totals)
print(sum(chunk_totals))
```

## Notes

- `read_csv_kwargs` is passed to `pandas.read_csv`, but it must not include `chunksize` because the snippet controls that value explicitly.
- Exceptions raised by your chunk processor are propagated so calling code can fail fast or log the exact chunk-level problem.
- Missing files and CSV parsing errors are not hidden by the snippet.

## Verification

Install the dependency and run the unit tests from the snippet root:

```bash
python -m pip install pandas
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- single-chunk processing for small CSV files
- stable processing order across multiple chunks
- `read_csv_kwargs` passthrough for column selection and dtypes
- header-only CSV handling
- missing file behavior
- invalid `chunk_size` values
- invalid `chunk_processor` values
- exception propagation from chunk-level processing
- conflicting `chunksize` arguments in `read_csv_kwargs`

## Files

- `src/process_csv_in_chunks.py`
- `tests/test_process_csv_in_chunks.py`
- `snippet.json`