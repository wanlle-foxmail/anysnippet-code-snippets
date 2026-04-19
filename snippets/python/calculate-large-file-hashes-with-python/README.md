# Calculate Large File Hashes with Python

Calculate a file hash without loading the whole file into memory.

This snippet is useful when you need a stable checksum for large files such as backups, release archives, or imported datasets.

## Highlights

- Reads files in fixed-size chunks to keep memory usage stable
- Accepts both `str` and `Path` inputs
- Supports `sha256` and `md5` with explicit validation
- Uses only the Python standard library

## Use Cases

- Verify archive integrity after download or deployment
- Generate checksums for large exports or backups
- Compare file content across machines without loading the file into memory

## Code

```python
from pathlib import Path

from src.calculate_large_file_hash import calculate_large_file_hash


digest = calculate_large_file_hash(Path("release.tar.gz"), algorithm="sha256")
print(digest)
```

## Notes

- `sha256` is the default and recommended option for modern integrity checks.
- `md5` is included for compatibility with legacy systems and non-security workflows. Do not use it for security-sensitive verification.

## Verification

Run the unit tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- SHA-256 hashing for a small file
- chunked hashing for a multi-chunk file
- hashing an empty file
- MD5 compatibility mode
- `str` and `Path` input handling
- directory path validation
- missing file errors
- invalid algorithm errors
- invalid chunk size errors

## Files

- `src/calculate_large_file_hash.py`
- `tests/test_calculate_large_file_hash.py`
- `snippet.json`