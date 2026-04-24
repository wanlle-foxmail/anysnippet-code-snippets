# Find Missing Related Rows with LEFT JOIN in MySQL

Return rows in MySQL that still have no related row by using `LEFT JOIN` and `IS NULL`.

This snippet is useful when you need to spot missing related data without rewriting the query as a subquery.

## Highlights

- Uses a standard anti-join
- Returns only unmatched rows
- Keeps result ordering simple

## Use Cases

- List users that do not have a profile row yet
- Find parent rows missing a child record
- Check which entities still need a related settings row

## Code

```sql
-- Use this query when you need rows from one table that still have no
-- related row in another table.
-- LEFT JOIN plus IS NULL is the standard anti-join form.

SELECT
    u.user_id,
    u.username
FROM users AS u
LEFT JOIN profiles AS p
    ON p.user_id = u.user_id
WHERE p.user_id IS NULL
ORDER BY u.user_id ASC;
```

## Notes

- Put the join match in the `ON` clause, not in `WHERE`, or the query stops behaving like a left join.
- The `IS NULL` filter should check a right-table column that is non-null for real matches.

## Verification

Run the integration tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- finding rows with no related match
- returning an empty result when every row is matched
- returning all rows when the related table is empty
- ignoring unrelated `NULL` keys on the right table
- keeping unmatched rows ordered by the left-table key

The tests start a temporary MySQL 8 instance with local `mysqld`, `mysqladmin`, and `mysql` binaries.

## Files

- `src/find_missing_related_rows_with_left_join.sql`
- `tests/test_find_missing_related_rows_with_left_join.py`
- `snippet.json`