# Write Files Atomically with Python

Write UTF-8 text to a temporary file in the target directory, then replace the target file atomically.

This snippet is useful when you need to update a config file, cache file, or generated JSON document without leaving half-written output behind.

## Highlights

- Replaces files atomically
- Keeps temp files local
- Cleans up on replace failure

## Use Cases

- Update JSON config files safely
- Replace generated index files after a batch job
- Refresh small cache files without partial writes

## Code

```python
from src.write_file_atomically import write_file_atomically


write_file_atomically("settings.json", '{"status": "ok"}')
print("settings.json")
```

## Notes

- The parent directory must already exist.
- This snippet writes UTF-8 text, not binary content.
- The temporary file is created in the target directory so `os.replace` stays on the same filesystem.

## Verification

Run the unit tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- writing a new file
- replacing an existing file
- empty text writes
- missing parent directory errors
- cleanup after replace failures
- Unicode text output

## Files

- `src/write_file_atomically.py`
- `tests/test_write_file_atomically.py`
- `snippet.json`