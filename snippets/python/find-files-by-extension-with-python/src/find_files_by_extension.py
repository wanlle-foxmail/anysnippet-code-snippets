import os
from pathlib import Path
from typing import Iterable, TypedDict, Union


PathLike = Union[str, Path]


class FileExtensionSearchResult(TypedDict):
    root_directory: str
    normalized_extensions: list[str]
    total_files: int
    matched_files: int
    skipped_files: int
    matching_file_paths: list[str]


def find_files_by_extension(
    directory: PathLike,
    extensions: Union[str, Iterable[str]],
) -> FileExtensionSearchResult:
    """Recursively collect files that match one or more normalized suffixes.

    Args:
        directory: Directory path as a string or Path object.
        extensions: One extension string or an iterable of extension strings.

    Returns:
        A dictionary with the normalized absolute root directory, normalized
        extension list, total file count, matched file count, skipped file
        count, and matching file paths in deterministic traversal order.

    Raises:
        FileNotFoundError: If the target directory does not exist.
        ValueError: If the target path exists but is not a real directory.
        ValueError: If no valid extension strings are provided.
    """
    target_dir = Path(directory)
    if not target_dir.exists():
        raise FileNotFoundError(f"Directory not found: {target_dir}")
    if target_dir.is_symlink():
        raise ValueError(f"Path must point to a real directory, not a symlink: {target_dir}")
    if not target_dir.is_dir():
        raise ValueError(f"Path must point to a directory: {target_dir}")

    normalized_extensions = _normalize_extensions(extensions)
    total_files = 0
    matching_file_paths: list[str] = []

    for current_root, dirnames, filenames in os.walk(target_dir, onerror=_raise_walk_error):
        dirnames[:] = [
            directory_name
            for directory_name in sorted(dirnames)
            if not (Path(current_root) / directory_name).is_symlink()
        ]
        filenames.sort()

        for filename in filenames:
            total_files += 1
            if _matches_extension(filename, normalized_extensions):
                matching_file_paths.append(os.path.abspath(os.path.join(current_root, filename)))

    matched_files = len(matching_file_paths)
    return {
        "root_directory": os.path.abspath(str(target_dir)),
        "normalized_extensions": normalized_extensions,
        "total_files": total_files,
        "matched_files": matched_files,
        "skipped_files": total_files - matched_files,
        "matching_file_paths": matching_file_paths,
    }


def _normalize_extensions(extensions: Union[str, Iterable[str]]) -> list[str]:
    if isinstance(extensions, str):
        candidates = [extensions]
    else:
        try:
            candidates = list(extensions)
        except TypeError as error:
            raise ValueError("extensions must be a string or an iterable of strings") from error

    normalized_extensions: list[str] = []
    seen_extensions: set[str] = set()
    for extension in candidates:
        if not isinstance(extension, str):
            raise ValueError("extensions must contain only non-empty strings")

        normalized_extension = extension.strip().lower()
        if not normalized_extension:
            raise ValueError("extensions must contain only non-empty strings")
        if not normalized_extension.startswith("."):
            normalized_extension = f".{normalized_extension}"
        if normalized_extension == ".":
            raise ValueError("extensions must contain only non-empty strings")
        if normalized_extension not in seen_extensions:
            seen_extensions.add(normalized_extension)
            normalized_extensions.append(normalized_extension)

    if not normalized_extensions:
        raise ValueError("extensions must contain at least one valid extension")

    return normalized_extensions


def _matches_extension(filename: str, normalized_extensions: list[str]) -> bool:
    lowercase_filename = filename.lower()
    return any(lowercase_filename.endswith(extension) for extension in normalized_extensions)


def _raise_walk_error(error: OSError) -> None:
    raise error