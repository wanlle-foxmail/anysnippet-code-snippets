import pandas


def read_parquet_records(path: str) -> list[dict[str, object]]:
    """Read a Parquet file with pandas and return row dictionaries."""
    if not isinstance(path, str):
        raise TypeError("path must be a string")
    if not path.strip():
        raise ValueError("path cannot be empty")

    dataframe = pandas.read_parquet(path, engine="pyarrow")
    return dataframe.to_dict("records")


def main() -> None:
    items = read_parquet_records("your_file.parquet")
    for index, item in enumerate(items):
        print(index, item)


if __name__ == "__main__":
    main()