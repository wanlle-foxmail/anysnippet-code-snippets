# Update Rows with JOIN in MySQL

Fill missing row values in MySQL from a lookup table with `UPDATE ... JOIN`.

This snippet is useful when one table stores item rows and another table stores the default value for each category.

## Highlights

- Uses MySQL update join syntax
- Fills only missing target values
- Reuses one lookup table cleanly

## Use Cases

- Copy a default rank into items that still have no rank
- Backfill one column from a category lookup table
- Update many rows with one MySQL join statement

## Code

```sql
-- Use this query when one table stores rows and another table stores
-- the default value for each category.
-- The WHERE clause keeps existing values unchanged.

UPDATE items AS i
JOIN category_defaults AS d
    ON d.category_name = i.category_name
SET i.display_rank = d.display_rank
WHERE i.display_rank IS NULL;
```

## Notes

- Keep the lookup table unique on the join key, or one target row can match more than one source row.
- Add the `WHERE` clause only when you want to preserve already populated values.

## Verification

Run the integration tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- updating rows that match the lookup table
- preserving rows that already have a target value
- leaving unmatched rows unchanged
- updating multiple rows from the same lookup value
- doing nothing when the lookup table is empty

The tests start a temporary MySQL 8 instance with local `mysqld`, `mysqladmin`, and `mysql` binaries.

## Files

- `src/update_rows_with_join.sql`
- `tests/test_update_rows_with_join.py`
- `snippet.json`