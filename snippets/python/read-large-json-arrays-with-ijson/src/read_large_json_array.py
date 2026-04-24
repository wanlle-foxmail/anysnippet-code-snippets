from collections.abc import Iterator
from typing import Any

import ijson


def read_large_json_array(path: str) -> Iterator[Any]:
    """Yield each item from a large top-level JSON array like ``[{...}, {...}]``."""
    with open(path, "rb") as file_handle:
        # "item" tells ijson to stream each element from the top-level array.
        items = ijson.items(file_handle, "item")
        for item in items:
            yield item


if __name__ == "__main__":
    # Example input shape: [{"question": "q1"}, {"question": "q2"}]
    for item in read_large_json_array("your_file.json"):
        print(item)