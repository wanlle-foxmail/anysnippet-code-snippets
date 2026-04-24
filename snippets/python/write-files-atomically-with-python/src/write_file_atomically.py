import os
import tempfile
from pathlib import Path
from typing import Optional


def write_file_atomically(path: str, content: str) -> None:
    """Write text to a temporary file in the target directory, then replace the target file."""
    target_path = Path(path)
    parent_dir = target_path.parent

    if not parent_dir.exists():
        raise FileNotFoundError(f"parent directory does not exist: {parent_dir}")
    if not parent_dir.is_dir():
        raise NotADirectoryError(f"parent path is not a directory: {parent_dir}")

    temp_path: Optional[Path] = None

    try:
        with tempfile.NamedTemporaryFile(
            mode="w",
            encoding="utf-8",
            dir=parent_dir,
            prefix=f".{target_path.name}.",
            suffix=".tmp",
            delete=False,
        ) as temp_file:
            temp_path = Path(temp_file.name)
            temp_file.write(content)
            temp_file.flush()
            os.fsync(temp_file.fileno())

        os.replace(temp_path, target_path)
    except Exception:
        if temp_path is not None and temp_path.exists():
            temp_path.unlink()
        raise


if __name__ == "__main__":
    write_file_atomically("settings.json", '{"status": "ok"}')
    print("settings.json")