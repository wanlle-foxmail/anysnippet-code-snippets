# Find Duplicate Values in MySQL

Return repeated non-null values in one MySQL column with `GROUP BY` and `HAVING`.

This snippet is useful when you need to spot duplicates without deleting data or building a cleanup workflow.

## Highlights

- Lists repeated column values
- Ignores `NULL` values cleanly
- Sorts larger duplicates first

## Use Cases

- Check which email addresses appear more than once
- Review duplicate import results before cleanup
- Find repeated identifiers in a staging table

## Code

```sql
-- Use this query when you need to list repeated non-null values in one column.
-- The HAVING clause keeps only values that appear more than once.

SELECT
    email,
    COUNT(*) AS duplicate_count
FROM users
WHERE email IS NOT NULL
GROUP BY email
HAVING COUNT(*) > 1
ORDER BY duplicate_count DESC, email ASC;
```

## Notes

- `HAVING COUNT(*) > 1` removes unique values after grouping.
- Keep the `WHERE email IS NOT NULL` filter when you want to ignore missing values.

## Verification

Run the integration tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- returning duplicate values with their counts
- excluding unique values from the result
- ignoring `NULL` values before grouping
- sorting larger duplicate groups first
- returning an empty result when all values are unique

The tests start a temporary MySQL 8 instance with local `mysqld`, `mysqladmin`, and `mysql` binaries.

## Files

- `src/find_duplicate_values.sql`
- `tests/test_find_duplicate_values.py`
- `snippet.json`