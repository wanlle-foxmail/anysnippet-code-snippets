# Sample Large Streams with Reservoir Sampling in Python

Sample a fixed number of items from a large stream without loading it all into memory.

This snippet is useful when logs, exports, or queue messages are too large to buffer fully but you still need a random sample for debugging or analysis.

## Highlights

- Samples without full buffering
- Works with one-pass iterables
- Supports repeatable tests

## Use Cases

- Sample large log streams
- Inspect a random subset of exported records
- Keep a fixed-size sample from one-pass input

## Code

```python
from src.reservoir_sampling import reservoir_sample


result = reservoir_sample(range(10_000), 5, seed=7)
print(result)
```

## Notes

- `sample_size` must be a positive integer.
- If the stream is shorter than `sample_size`, the function returns all items.
- Use `seed` when you want deterministic results for tests or repeatable debugging.

## Verification

Run the unit tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- deterministic sampling with a seed
- short-stream behavior
- empty input handling
- generator input
- full-length sampling
- invalid sample-size errors

## Files

- `src/reservoir_sampling.py`
- `tests/test_reservoir_sampling.py`
- `snippet.json`