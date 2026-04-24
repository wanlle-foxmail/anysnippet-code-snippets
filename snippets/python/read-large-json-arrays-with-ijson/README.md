# Read Large JSON Arrays with ijson

Stream a top-level JSON array with ijson and yield one parsed item at a time.

This snippet is useful when one file contains a very large JSON array and loading the whole document would risk running out of memory.

## Highlights

- Streams one array item at a time
- Avoids full-file memory loads
- Keeps nested JSON objects intact

## Use Cases

- Process large crawler exports wrapped in one JSON array
- Read question and answer records without loading the whole file
- Feed downstream ETL steps one item at a time

## Code

```python
from src.read_large_json_array import read_large_json_array


for item in read_large_json_array("your_file.json"):
    print(item)
```

## Input Shape

The file must contain one top-level JSON array.

```json
[
    {"question": "xxx111", "answer": "aaa", "crawled_time": "2025-05-01 12:13:14"},
    {"question": "xxx222", "answer": "bbb", "crawled_time": "2025-05-01 12:13:14"}
]
```

## Notes

- The file must contain one top-level JSON array such as `[{}, {}]`.
- `ijson.items(..., "item")` streams one element from that top-level array at a time.
- The source opens the file in binary mode because `ijson` works best with binary file handles.
- The file extension does not matter. A file named `.jsonl` still works if its actual content is one JSON array.

## Verification

Run the unit tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

Install the dependency first if needed:

```bash
python -m pip install ijson
```

The verified test suite covers:

- ordered reads across multiple items
- nested object and list preservation
- empty array handling
- malformed JSON array errors
- missing file behavior
- Unicode values

## Files

- `src/read_large_json_array.py`
- `tests/test_read_large_json_array.py`
- `snippet.json`