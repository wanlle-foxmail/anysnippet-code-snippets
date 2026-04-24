# Parse Pagination Query Params Safely with Python

Parse `page` and `page_size` query parameters with defaults, validation, and a fixed maximum page size.

This snippet is useful when an API should accept simple page-based pagination without letting callers request invalid or unbounded page sizes.

- Uses defaults when params are missing
- Rejects invalid or non-positive values
- Caps large page sizes at a fixed maximum

## Example

```python
from src.parse_pagination_query_params import parse_pagination_query_params

params = parse_pagination_query_params({"page": "3", "page_size": "50"})
print(params)
```

Output:

```python
PaginationParams(page=3, page_size=50, offset=100)
```

## Notes

- Missing values fall back to `page=1` and `page_size=20`.
- Large page sizes are capped at `100`.
- This helper parses only `page` and `page_size`; filtering and sorting should stay separate.

## Verification

Run the tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

Verified behavior covers:

- valid page and page size parsing
- default values for missing params
- page size capping
- invalid integer input
- non-positive values
- surrounding whitespace in integer values

## Files

- `src/parse_pagination_query_params.py`
- `tests/test_parse_pagination_query_params.py`