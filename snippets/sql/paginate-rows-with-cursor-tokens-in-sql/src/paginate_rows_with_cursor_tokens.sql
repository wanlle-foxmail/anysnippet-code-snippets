-- Use this query when the client already decoded the cursor into the last
-- seen created_at and item_id values from the previous page.
-- The WHERE clause advances to rows strictly older than the cursor row while
-- keeping matching timestamps deterministic with item_id.
-- The example cursor values are shown inline for readability.
-- In application code, replace them with bound parameters instead of
-- interpolating dynamic values into the SQL string.

SELECT
    item_id,
    item_name,
    created_at
FROM items
WHERE created_at < '2026-04-01 12:00:00'
   OR (created_at = '2026-04-01 12:00:00' AND item_id < 5)
ORDER BY created_at DESC, item_id DESC
LIMIT 3;