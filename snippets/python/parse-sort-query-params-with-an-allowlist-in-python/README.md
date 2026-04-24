# Parse Sort Query Params with an Allowlist in Python

Parse a comma-separated `sort` query parameter against a fixed allowlist of field names.

This snippet is useful when an API should accept sorting from callers without letting arbitrary field names leak into later query-building code.

- Parses comma-separated sort terms
- Supports descending terms with `-field`
- Rejects unknown, empty, or repeated fields

## Example

```python
from src.parse_sort_query_params import parse_sort_query_params

sort_fields = parse_sort_query_params(
    {"sort": "name,-created_at"},
    ["name", "created_at"],
)

print(sort_fields)
```

Output:

```python
[SortField(field_name='name', descending=False), SortField(field_name='created_at', descending=True)]
```

## Notes

- Missing `sort` returns an empty list.
- `-field` means descending order.
- This helper parses and validates only the sort input; it does not build SQL or ORM expressions.

## Verification

Run the tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

Verified behavior covers:

- one ascending field
- multiple fields with mixed directions
- a missing sort param
- unknown fields
- empty terms
- repeated fields

## Files

- `src/parse_sort_query_params.py`
- `tests/test_parse_sort_query_params.py`