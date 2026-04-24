from typing import Callable

import pandas


def process_csv_in_chunks(
    csv_path: str,
    chunk_processor: Callable[[pandas.DataFrame], object],
    chunk_size: int = 10000,
) -> list[object]:
    """Read a CSV in chunks and return one processed result per chunk."""
    # Flow:
    #   CSV reader -> yield one chunk at a time
    #                 |
    #                 +-> empty chunk -> skip it
    #                 `-> non-empty chunk -> process it and collect one result
    if isinstance(chunk_size, bool) or not isinstance(chunk_size, int) or chunk_size <= 0:
        raise ValueError("chunk_size must be a positive integer")

    results = []
    chunk_iter = pandas.read_csv(csv_path, chunksize=chunk_size)

    for chunk in chunk_iter:
        if chunk.empty:
            continue
        results.append(chunk_processor(chunk))

    return results


if __name__ == "__main__":
    def to_rows(chunk: pandas.DataFrame) -> list[dict[str, object]]:
        return chunk.to_dict(orient="records")


    result = process_csv_in_chunks("your_file.csv", to_rows)
    print(result)