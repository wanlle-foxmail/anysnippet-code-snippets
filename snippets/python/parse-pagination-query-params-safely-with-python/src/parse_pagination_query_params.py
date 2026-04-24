from dataclasses import dataclass
from typing import Mapping


DEFAULT_PAGE = 1
DEFAULT_PAGE_SIZE = 20
MAX_PAGE_SIZE = 100


@dataclass(frozen=True)
class PaginationParams:
    page: int
    page_size: int
    offset: int


def parse_pagination_query_params(query_params: Mapping[str, str]) -> PaginationParams:
    """Parse page and page_size query params with defaults and limits."""
    # Flow:
    #   read page and page_size from the query mapping
    #      |
    #      +-> valid integers -> apply defaults and max limit -> return params
    #      `-> invalid or out-of-range values -> raise ValueError
    page = _read_positive_int(query_params, "page", default=DEFAULT_PAGE)
    page_size = _read_positive_int(query_params, "page_size", default=DEFAULT_PAGE_SIZE)
    bounded_page_size = min(page_size, MAX_PAGE_SIZE)
    offset = (page - 1) * bounded_page_size
    return PaginationParams(page=page, page_size=bounded_page_size, offset=offset)


def _read_positive_int(query_params: Mapping[str, str], key: str, default: int) -> int:
    raw_value = query_params.get(key)
    if raw_value is None:
        return default
    if not isinstance(raw_value, str):
        raise ValueError(f"{key} must be a string integer")

    try:
        value = int(raw_value.strip())
    except ValueError as error:
        raise ValueError(f"{key} must be a valid integer") from error

    if value < 1:
        raise ValueError(f"{key} must be greater than or equal to 1")
    return value


if __name__ == "__main__":
    params = parse_pagination_query_params({"page": "3", "page_size": "50"})
    print(params)