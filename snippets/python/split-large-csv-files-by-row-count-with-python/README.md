# Split Large CSV Files by Row Count with Python

Split a large CSV file into smaller UTF-8 CSV files with a fixed number of data rows per part.

This snippet is useful when upload limits, import tools, or manual review workflows work better with smaller CSV files.

## Highlights

- Repeats the header row
- Creates output parts in order
- Uses Python stdlib only

## Use Cases

- Break large imports into smaller files
- Prepare CSV uploads for row-limited tools
- Split exported data for manual review

## Code

```python
from src.split_csv_file_by_row_count import split_csv_file_by_row_count


result = split_csv_file_by_row_count("events.csv", "csv-parts", 1000)
print(result)
```

## Notes

- This snippet treats the first row as the CSV header.
- Each output file repeats the header row.
- If the input file has only a header row, no output files are created.

## Verification

Run the unit tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- multi-part splits with repeated headers
- output directory creation
- single-file output when the limit is large
- header-only input handling
- missing input files
- invalid row-limit errors

## Files

- `src/split_csv_file_by_row_count.py`
- `tests/test_split_csv_file_by_row_count.py`
- `snippet.json`