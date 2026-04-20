import pandas
from typing import Iterator


def read_headerless_csv_in_chunks(csv_path: str, cols: list[str], chunk_size: int = 10000) -> Iterator[dict[str, object]]:
    """Read a headerless CSV in chunks and yield row dictionaries."""
    chunk_iter = pandas.read_csv(csv_path, header=None, names=cols, chunksize=chunk_size)

    for chunk in chunk_iter:
        # Turn each chunk into row dictionaries.
        items = chunk.to_dict(orient="records")
        for item in items:
            yield item


if __name__ == "__main__":
    # The file has no header row, so define the column names explicitly.
    cols = ["name", "age", "city"]

    for item in read_headerless_csv_in_chunks("your_file.csv", cols):
        print(item)