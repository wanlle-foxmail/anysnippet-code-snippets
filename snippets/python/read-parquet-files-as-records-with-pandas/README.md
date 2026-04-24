# Read Parquet Files as Records with pandas

Load a Parquet file with pandas and return Python row dictionaries.

This snippet is useful when you want to inspect a downloaded Parquet dataset, pass rows into plain Python code, or serialize row records for another step in your pipeline.

## Highlights

- Reads parquet with pyarrow
- Returns one dict per row
- Keeps row order stable

## Code

```python
from src.read_parquet_records import read_parquet_records


items = read_parquet_records("your_file.parquet")
for index, item in enumerate(items):
    print(index, item)
```

## Notes

- This snippet uses the `pyarrow` engine explicitly.
- `to_dict("records")` returns one dictionary per row.
- Null values in object columns stay as `None`, but numeric nulls follow pandas behavior and may appear as `NaN` after conversion.
- The full file is loaded into a pandas DataFrame before conversion, so this pattern is best for small or medium Parquet files.

## Verification

Install the dependencies and run the unit tests from the snippet root:

```bash
python -m pip install pandas pyarrow
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- multi-row record reads in order
- empty parquet file handling
- Unicode value preservation
- null value preservation in object columns
- missing file behavior
- non-string path rejection
- empty path validation

## Files

- `src/read_parquet_records.py`
- `tests/test_read_parquet_records.py`
- `snippet.json`