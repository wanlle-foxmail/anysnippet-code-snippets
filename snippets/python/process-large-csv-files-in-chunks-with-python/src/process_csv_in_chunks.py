from os import PathLike
from pathlib import Path
from typing import Any, Callable, List, Mapping, Optional, TypeVar, Union

import pandas as pd


ProcessedChunk = TypeVar("ProcessedChunk")
CsvPath = Union[str, PathLike[str]]


def process_csv_in_chunks(
    csv_path: CsvPath,
    chunk_processor: Callable[[pd.DataFrame], ProcessedChunk],
    *,
    chunk_size: int = 10000,
    read_csv_kwargs: Optional[Mapping[str, Any]] = None,
) -> List[ProcessedChunk]:
    if isinstance(chunk_size, bool) or not isinstance(chunk_size, int) or chunk_size <= 0:
        raise ValueError("chunk_size must be a positive integer")

    if not callable(chunk_processor):
        raise ValueError("chunk_processor must be callable")

    if read_csv_kwargs is None:
        reader_kwargs = {}
    else:
        if not isinstance(read_csv_kwargs, Mapping):
            raise ValueError("read_csv_kwargs must be a mapping")
        if "chunksize" in read_csv_kwargs:
            raise ValueError("read_csv_kwargs must not include chunksize")
        reader_kwargs = dict(read_csv_kwargs)

    try:
        normalized_path = Path(csv_path)
    except TypeError as exc:
        raise ValueError("csv_path must be a string or Path-like value") from exc

    processed_results = []

    with pd.read_csv(normalized_path, chunksize=chunk_size, **reader_kwargs) as chunk_reader:
        for chunk in chunk_reader:
            if chunk.empty:
                continue
            processed_results.append(chunk_processor(chunk))

    return processed_results