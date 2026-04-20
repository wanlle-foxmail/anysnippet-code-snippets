# Upsert Tenant Contacts via SQLite ON CONFLICT

Idempotently insert or update tenant-scoped contacts with one SQLite-compatible statement.

This snippet is useful when sync jobs, imports, or admin tools need repeated writes to the same tenant contact table without creating duplicates.

## Highlights

- Uses `ON CONFLICT ... DO UPDATE` for idempotent writes
- Targets a composite unique key for tenant-scoped records
- Works for both single-row execution and repeated batch bindings

## Use Cases

- Sync customer contacts from external systems
- Import tenant-scoped reference data without duplicate rows
- Refresh cached profile data after a scheduled pull

## Code

```sql
-- Use this query when repeated sync writes should insert new tenant contacts
-- or update existing ones without creating duplicates.
-- Requires a UNIQUE constraint on (tenant_id, email).

INSERT INTO customer_contacts (
    tenant_id,
    email,
    full_name,
    phone,
    updated_at
)
VALUES (?, ?, ?, ?, ?)
ON CONFLICT (tenant_id, email) DO UPDATE SET
    full_name = excluded.full_name,
    phone = excluded.phone,
    updated_at = excluded.updated_at;
```

## Notes

- The conflict target must match a real unique constraint, such as `UNIQUE (tenant_id, email)`.
- This query is verified with SQLite. PostgreSQL uses similar `ON CONFLICT` syntax, while MySQL and SQL Server require different upsert forms.

## Verification

Run the integration tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- inserting a new contact when no conflict exists
- updating an existing contact on a composite-key conflict
- keeping the same email separate across tenants
- applying null updates safely
- handling batch upserts with repeated bindings
- leaving unrelated rows unchanged

## Files

- `src/upsert_rows_on_conflict.sql`
- `tests/test_upsert_rows_on_conflict.py`
- `snippet.json`