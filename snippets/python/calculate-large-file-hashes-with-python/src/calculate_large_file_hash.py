import hashlib
from pathlib import Path
from typing import Literal, Union


PathLike = Union[str, Path]
HashAlgorithm = Literal["md5", "sha256"]


def calculate_large_file_hash(
    path: PathLike,
    algorithm: HashAlgorithm = "sha256",
    chunk_size: int = 1024 * 1024,
) -> str:
    """Return the lowercase hex digest for a file read chunk by chunk.

    Args:
        path: File path as a string or Path object.
        algorithm: Hash algorithm to use. Choose "sha256" for modern integrity
            checks or "md5" only for compatibility with legacy systems.
        chunk_size: Number of bytes to read per iteration.

    Returns:
        The lowercase hexadecimal digest string for the file.

    Raises:
        ValueError: If the algorithm is unsupported or chunk_size is not positive.
        ValueError: If the target path exists but is not a regular file.
        FileNotFoundError: If the target path does not exist.
    """
    if not isinstance(algorithm, str) or algorithm not in {"md5", "sha256"}:
        raise ValueError("algorithm must be 'md5' or 'sha256'")
    if chunk_size <= 0:
        raise ValueError("chunk_size must be greater than 0")

    file_path = Path(path)
    if not file_path.exists():
        raise FileNotFoundError(f"File not found: {file_path}")
    if not file_path.is_file():
        raise ValueError(f"Path must point to a regular file: {file_path}")

    hasher = hashlib.new(algorithm)

    with file_path.open("rb") as file_handle:
        while True:
            # Stream fixed-size chunks to keep memory usage stable for large files.
            chunk = file_handle.read(chunk_size)
            if not chunk:
                break
            hasher.update(chunk)

    return hasher.hexdigest()