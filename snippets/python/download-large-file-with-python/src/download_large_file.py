import hashlib
import os

import requests


def stream_write_and_hash(chunks, file_handle):
    """Write byte chunks to a file while computing their MD5 hash."""
    hasher = hashlib.md5()
    written = 0
    for chunk in chunks:
        file_handle.write(chunk)
        hasher.update(chunk)
        written += len(chunk)
    return written, hasher.hexdigest()


def download_large_file(url: str, save_path: str, chunk_size: int = 1024 * 1024, timeout: int = 30) -> dict:
    """Download a large file in chunks, compute MD5 on the fly, clean up on failure."""
    try:
        with requests.get(url, stream=True, timeout=timeout) as resp:
            resp.raise_for_status()
            with open(save_path, "wb") as f:
                size, md5 = stream_write_and_hash(
                    resp.iter_content(chunk_size=chunk_size), f
                )
    except (requests.exceptions.RequestException, OSError):
        if os.path.exists(save_path):
            os.remove(save_path)
        raise
    return {"path": save_path, "hash": md5, "size": size}


if __name__ == "__main__":
    result = download_large_file("https://example.com/large-file.zip", "large-file.zip")
    print(result)
