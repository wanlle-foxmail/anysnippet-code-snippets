# Compress Text for Storage with Python

Compress structured text into a base64-encoded gzip string so it uses less space in text-only storage.

This snippet is useful when you want to store large JSON documents, editor state, or repeated log text in a cache or database column that only accepts strings.

## Highlights

- Base64 gzip string output
- Saves space for JSON text
- Unicode-safe round-trip

## Use Cases

- Store large JSON payloads in a text column
- Shrink cached editor state before saving it
- Reduce the size of repeated log or event text

## Code

```python
from src.compress_text_for_storage import compress_text, decompress_text


compressed_text = compress_text('{"page": "/docs", "events": ["view", "view", "view"]}')
print(compressed_text)
print(decompress_text(compressed_text))
```

## Notes

- Base64 keeps the output string-safe for text-only storage.
- Space savings are best for structured or repetitive text such as JSON, HTML, and logs.
- Invalid base64 or non-gzip input raises `ValueError` during decompression.

## Verification

Run the unit tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- JSON text round-trip
- empty string round-trip
- Unicode and emoji round-trip
- size reduction for repetitive JSON text
- invalid base64 errors
- invalid gzip payload errors
- non-string input rejection for both functions

## Files

- `src/compress_text_for_storage.py`
- `tests/test_compress_text_for_storage.py`
- `snippet.json`