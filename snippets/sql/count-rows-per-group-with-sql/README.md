# Count Rows per Group in MySQL

Return one row count per group value in MySQL with `GROUP BY` and `COUNT(*)`.

This snippet is useful when you need a quick grouped summary without adding extra joins or dashboard logic.

## Highlights

- Counts rows per group value
- Sorts larger groups first
- Keeps tie ordering stable

## Use Cases

- Count rows for each category in a table
- Summarize grouped data before a report export
- Check how evenly values are distributed across groups

## Code

```sql
-- Use this query when you need a quick row count for each group value.
-- The ORDER BY clause puts larger groups first and keeps ties stable.

SELECT
    category_name,
    COUNT(*) AS row_count
FROM items
GROUP BY category_name
ORDER BY row_count DESC, category_name ASC;
```

## Notes

- `COUNT(*)` counts every row in each group, even when other columns contain `NULL`.
- Add `HAVING` only when you need to filter out smaller groups.

## Verification

Run the integration tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- counting rows across multiple groups
- sorting larger groups before smaller groups
- breaking equal counts by the group name
- returning one group when all rows match
- counting rows even when other columns are `NULL`

The tests start a temporary MySQL 8 instance with local `mysqld`, `mysqladmin`, and `mysql` binaries.

## Files

- `src/count_rows_per_group.sql`
- `tests/test_count_rows_per_group.py`
- `snippet.json`