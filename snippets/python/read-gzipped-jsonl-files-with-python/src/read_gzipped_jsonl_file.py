import gzip
import json
from collections.abc import Iterator
from typing import Any


def read_gzipped_jsonl_file(path: str) -> Iterator[Any]:
    """Yield one parsed JSON value per non-empty line from a gzipped JSONL file."""
    # Flow:
    #   gzip stream -> strip each incoming line
    #                  |
    #                  +-> empty line -> skip it
    #                  +-> valid JSON -> parse and yield one item
    #                  `-> invalid JSON -> raise a line-number error
    with gzip.open(path, "rt", encoding="utf-8") as file_handle:
        for line_number, raw_line in enumerate(file_handle, start=1):
            line = raw_line.strip()
            if not line:
                continue

            try:
                yield json.loads(line)
            except json.JSONDecodeError as error:
                raise ValueError(f"invalid JSON on line {line_number}") from error


if __name__ == "__main__":
    for item in read_gzipped_jsonl_file("events.jsonl.gz"):
        print(item)