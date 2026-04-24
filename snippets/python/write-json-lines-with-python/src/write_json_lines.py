import json
import os
import tempfile
from collections.abc import Iterable
from pathlib import Path
from typing import Any, Optional


def write_json_lines(path: str, items: Iterable[Any]) -> None:
    """Write one JSON value per line to a temporary file, then replace the target file."""
    # Flow: items -> write temp file -> flush + fsync -> atomic replace
    #        any error -> remove temp file and re-raise
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
            for item in items:
                temp_file.write(json.dumps(item, ensure_ascii=False))
                temp_file.write("\n")
            temp_file.flush()
            os.fsync(temp_file.fileno())

        os.replace(temp_path, target_path)
    except Exception:
        if temp_path is not None and temp_path.exists():
            temp_path.unlink()
        raise


if __name__ == "__main__":
    write_json_lines("events.jsonl", [{"id": 1}, {"id": 2}])
    print("events.jsonl")