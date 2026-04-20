# Calculate Large File Hashes with Python

Hash a large file in chunks with SHA-256 or MD5.

This snippet is useful when you need a stable checksum for large files such as backups, release archives, or imported datasets.

## Highlights

- Reads files chunk by chunk
- Supports sha256 and md5
- Uses Python stdlib only

## Use Cases

- Verify archive integrity after download or deployment
- Generate checksums for large exports or backups
- Compare file content across machines without loading the file into memory

## Code

```python
from src.calculate_large_file_hash import calculate_large_file_hash


digest = calculate_large_file_hash("release.tar.gz")
print(digest)
```

## Notes

- `sha256` is the default and recommended option for modern integrity checks.
- `md5` is included for compatibility with older workflows.
- `chunk_size` controls how many bytes are read at a time.

## Verification

Run the unit tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- SHA-256 hashing for a small file
- chunked hashing for a multi-chunk file
- MD5 compatibility mode
- invalid algorithm errors
- invalid chunk size errors

## Files

- `src/calculate_large_file_hash.py`
- `tests/test_calculate_large_file_hash.py`
- `snippet.json`