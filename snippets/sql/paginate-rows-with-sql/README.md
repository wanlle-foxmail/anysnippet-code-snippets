# Paginate Rows with LIMIT and OFFSET in MySQL

Return one stable page of rows in MySQL with `ORDER BY`, `LIMIT`, and `OFFSET`.

This snippet is useful when you need a simple page-based list and deterministic row order.

## Highlights

- Uses stable two-column sorting
- Shows one fixed page slice
- Works for basic page-based lists

## Use Cases

- Return the second page of a table sorted by creation time
- Build a simple admin list without cursor pagination
- Keep page contents deterministic when timestamps match

## Code

```sql
-- Use this query when you need one stable page from a larger result set.
-- Replace LIMIT and OFFSET with the page size and starting row you need.

SELECT
    item_id,
    item_name,
    created_at
FROM items
ORDER BY created_at DESC, item_id DESC
LIMIT 3 OFFSET 3;
```

## Notes

- Always include `ORDER BY` before `LIMIT` and `OFFSET`, or page contents can shift unpredictably.
- Offset-based pagination is simple, but large offsets get slower as tables grow.

## Verification

Run the integration tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- returning the second page of rows from a larger result set
- enforcing the page size when more rows are available
- using `item_id` as a tie-breaker for matching timestamps
- returning a partial page when fewer rows remain
- returning an empty result when the offset reaches the table end

The tests start a temporary MySQL 8 instance with local `mysqld`, `mysqladmin`, and `mysql` binaries.

## Files

- `src/paginate_rows.sql`
- `tests/test_paginate_rows.py`
- `snippet.json`