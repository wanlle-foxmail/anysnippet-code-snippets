# Paginate Rows with Cursor Tokens in MySQL

Return the next stable page of rows in MySQL with keyset pagination on `created_at` and `item_id`.

This snippet is useful when you want cursor-based pagination instead of `OFFSET` so deep pages stay predictable and efficient.

## Highlights

- Uses keyset pagination
- Keeps timestamp ties stable
- Avoids large OFFSET scans

## Use Cases

- Return the next page of an activity feed
- Page through large admin tables with better deep-page performance
- Keep cursor pagination deterministic when timestamps match

## Code

```sql
-- Use this query when the client already decoded the cursor into the last
-- seen created_at and item_id values from the previous page.
-- The WHERE clause advances to rows strictly older than the cursor row while
-- keeping matching timestamps deterministic with item_id.

SELECT
    item_id,
    item_name,
    created_at
FROM items
WHERE created_at < '2026-04-01 12:00:00'
   OR (created_at = '2026-04-01 12:00:00' AND item_id < 5)
ORDER BY created_at DESC, item_id DESC
LIMIT 3;
```

## Notes

- The cursor values come from the last row of the previous page.
- The first page omits the cursor `WHERE` clause and starts directly from `ORDER BY ... LIMIT ...`.
- The secondary comparison on `item_id` keeps pagination deterministic when timestamps match.
- Back this query with a composite index on `(created_at DESC, item_id DESC)` before using it on large tables.
- In application code, pass the cursor values with parameter binding instead of string interpolation.

## Verification

Run the integration tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- returning the next page after a cursor row
- enforcing the page size limit
- using `item_id` as the cursor tie-breaker
- returning a partial page near the end
- returning an empty result when no rows remain after the cursor

The tests start a temporary MySQL 8 instance with local `mysqld`, `mysqladmin`, and `mysql` binaries.

## Files

- `src/paginate_rows_with_cursor_tokens.sql`
- `tests/test_paginate_rows_with_cursor_tokens.py`
- `snippet.json`