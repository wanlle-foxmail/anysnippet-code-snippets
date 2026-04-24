# Merge Overlapping Time Ranges with Python

Sort numeric time ranges by start time and merge overlapping intervals.

This snippet is useful when bookings, maintenance windows, or schedule blocks may overlap and you want one merged list of occupied time ranges.

## Highlights

- Sorts ranges before merging
- Merges overlaps in one pass
- Rejects invalid ranges

## Use Cases

- Merge overlapping booking windows
- Collapse maintenance time ranges
- Normalize busy schedule blocks before display

## Code

```python
from src.merge_overlapping_time_ranges import merge_overlapping_time_ranges


result = merge_overlapping_time_ranges([(60, 120), (90, 180), (240, 300)])
print(result)
```

## Notes

- Each range must be a 2-item tuple of integers like `(start, end)`.
- The function sorts ranges by start time before merging.
- Ranges are treated as closed intervals, so `(1, 3)` and `(3, 5)` merge into `(1, 5)`.

## Verification

Run the unit tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- unsorted overlapping ranges
- separate non-overlapping ranges
- shared-boundary merges
- nested ranges
- empty input handling
- reversed range errors
- bool rejection for integer fields

## Files

- `src/merge_overlapping_time_ranges.py`
- `tests/test_merge_overlapping_time_ranges.py`
- `snippet.json`