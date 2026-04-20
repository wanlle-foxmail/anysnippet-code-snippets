# Read Headerless CSV in Chunks with pandas

Read a headerless CSV file in chunks and yield row dictionaries.

This snippet is useful when the source file does not include a header row and each line already contains data, for example:

```text
Alice,29,Shanghai
Bob,31,Beijing
```

## Highlights

- Reads headerless CSV files
- Defines column names
- Yields row dicts

## Code

```python
import pandas
from typing import Iterator


def read_headerless_csv_in_chunks(csv_path: str, cols: list[str], chunk_size: int = 10000) -> Iterator[dict[str, object]]:
    chunk_iter = pandas.read_csv(csv_path, header=None, names=cols, chunksize=chunk_size)

    for chunk in chunk_iter:
        # Turn each chunk into row dictionaries.
        items = chunk.to_dict(orient="records")
        for item in items:
            yield item


cols = ["name", "age", "city"]
for item in read_headerless_csv_in_chunks("your_file.csv", cols):
    print(item)
```

## Notes

- This snippet keeps a small function so the basic call shape is visible in the source file.
- The function yields one row dictionary at a time.
- Set `header=None` because the first row is data, not column names.
- Replace `cols` with the actual column names for your CSV layout.
- Increase or decrease `chunksize` based on file size and available memory.

## Verification

Install the dependency and run the unit tests from the snippet root:

```bash
python -m pip install pandas
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- row dictionaries for single-chunk input
- multi-chunk iteration in order
- empty file handling
- single-row CSV handling
- missing file behavior

## Files

- `src/read_headerless_csv_in_chunks.py`
- `tests/test_read_headerless_csv_in_chunks.py`
- `snippet.json`