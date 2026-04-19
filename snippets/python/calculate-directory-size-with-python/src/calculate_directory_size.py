import os
from pathlib import Path
from typing import TypedDict, Union


PathLike = Union[str, Path]


class DirectorySizeResult(TypedDict):
    root_directory: str
    total_bytes: int
    file_count: int
    subdirectory_count: int


def calculate_directory_size(directory: PathLike) -> DirectorySizeResult:
    """Recursively calculate total file size and directory counts for a tree.

    Args:
        directory: Directory path as a string or Path object.

    Returns:
        A dictionary containing the normalized absolute root directory, total
        file bytes, file count, and subdirectory count excluding the root
        directory and any nested symlink entries.

    Raises:
        FileNotFoundError: If the target directory does not exist.
        ValueError: If the target path exists but is not a real directory.
    """
    target_dir = Path(directory)
    if not target_dir.exists():
        raise FileNotFoundError(f"Directory not found: {target_dir}")
    if target_dir.is_symlink():
        raise ValueError(f"Path must point to a real directory, not a symlink: {target_dir}")
    if not target_dir.is_dir():
        raise ValueError(f"Path must point to a directory: {target_dir}")

    total_bytes = 0
    file_count = 0
    subdirectory_count = 0

    for current_root, dirnames, filenames in os.walk(target_dir, onerror=_raise_walk_error):
        dirnames[:] = [
            directory_name
            for directory_name in sorted(dirnames)
            if not (Path(current_root) / directory_name).is_symlink()
        ]
        filenames.sort()
        subdirectory_count += len(dirnames)

        for filename in filenames:
            file_path = Path(current_root) / filename
            if file_path.is_symlink():
                continue
            total_bytes += file_path.stat().st_size
            file_count += 1

    return {
        "root_directory": os.path.abspath(str(target_dir)),
        "total_bytes": total_bytes,
        "file_count": file_count,
        "subdirectory_count": subdirectory_count,
    }


def _raise_walk_error(error: OSError) -> None:
    raise error