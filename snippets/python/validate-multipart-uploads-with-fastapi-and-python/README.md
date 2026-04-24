# Validate Multipart Uploads with FastAPI

Validate one multipart upload in FastAPI with filename, content type, and file size checks.

This snippet is useful when an API should accept one small upload and reject empty files, oversized bodies, or unsupported content types before further processing.

## Highlights

- Checks filename and MIME type
- Rejects empty or oversized files
- Returns accepted file metadata

## What It Does

- Defines a `POST /upload` endpoint with one `UploadFile`
- Requires a non-empty filename
- Allows only `text/plain` and `text/csv`
- Reads the upload in chunks and rejects empty or oversized bodies
- Returns the accepted filename, content type, and size

## Usage

```python
from upload_validation import app

# Run directly:
# python src/upload_validation.py
# Then send one multipart POST request to /upload
```

## Notes

- This snippet checks file size incrementally so it can reject oversized uploads without keeping the whole body in memory.
- Replace the allowed content types and size limit with values that match your upload endpoint.
- FastAPI multipart parsing requires the `python-multipart` package.
- `content_type` is still client-supplied metadata, so production-sensitive uploads should add file signature checks too.
- Add authentication and rate limiting before using this pattern on a public upload endpoint.

## Verification

Run the unit tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- accepting a valid text upload
- rejecting blank filenames
- rejecting unsupported content types
- rejecting empty files
- rejecting oversized files
- accepting files at the exact size limit
- returning a sanitized basename for uploaded filenames
- keeping the allowlist explicit

## Files

- `src/upload_validation.py`
- `tests/test_upload_validation.py`
- `snippet.json`