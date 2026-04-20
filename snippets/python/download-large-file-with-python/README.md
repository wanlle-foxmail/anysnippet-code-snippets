# Download Large File with Python

Stream-download a large file with `requests`, compute MD5 on the fly, and clean up on failure.

## What It Does

- Opens an HTTP connection with `stream=True` so the response body is not loaded into memory
- Writes data chunk by chunk to a local file
- Computes the MD5 hash incrementally while writing — zero extra I/O
- If the download fails (timeout, HTTP error, mid-stream disconnect), removes the partial file

## Usage

```python
from download_large_file import download_large_file

result = download_large_file("https://example.com/large-file.zip", "large-file.zip")
print(result)
# {"path": "large-file.zip", "hash": "d41d8cd9...", "size": 104857600}
```

The return value is a dict with three keys:

| Key | Type | Description |
|---|---|---|
| `path` | `str` | The local file path where the file was saved |
| `hash` | `str` | MD5 hex digest computed during download |
| `size` | `int` | Total bytes written |

## Verification

```bash
cd snippets/python/download-large-file-with-python
pip install requests
python -m unittest discover -s tests -p "test_*.py"
```
