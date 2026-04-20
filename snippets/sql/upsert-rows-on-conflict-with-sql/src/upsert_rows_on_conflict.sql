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