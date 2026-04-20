import os


def calculate_directory_size(source: str) -> dict[str, int]:
    """Walk a directory tree and return size and item counts."""
    if not os.path.isdir(source):
        raise ValueError("source must be a directory")

    data = {"total_bytes": 0, "file_count": 0, "subdirectory_count": 0}

    for root, dirs, files in os.walk(source):
        data["subdirectory_count"] += len(dirs)
        for file in files:
            filename = os.path.join(root, file)
            data["total_bytes"] += os.path.getsize(filename)
            data["file_count"] += 1

    return data


if __name__ == "__main__":
    result = calculate_directory_size("Your Directory")
    print(result)