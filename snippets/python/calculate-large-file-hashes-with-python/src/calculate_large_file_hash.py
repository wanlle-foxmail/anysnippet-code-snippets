import hashlib


def calculate_large_file_hash(path: str, algorithm: str = "sha256", chunk_size: int = 1024 * 1024) -> str:
    """Hash a file in chunks so large files do not need to fit in memory."""
    if algorithm not in {"md5", "sha256"}:
        raise ValueError("algorithm must be 'md5' or 'sha256'")
    if isinstance(chunk_size, bool) or not isinstance(chunk_size, int) or chunk_size <= 0:
        raise ValueError("chunk_size must be a positive integer")

    hasher = hashlib.new(algorithm)

    with open(path, "rb") as file_handle:
        while True:
            chunk = file_handle.read(chunk_size)
            if not chunk:
                break
            hasher.update(chunk)

    return hasher.hexdigest()


if __name__ == "__main__":
    digest = calculate_large_file_hash("Your Filename")
    print(digest)