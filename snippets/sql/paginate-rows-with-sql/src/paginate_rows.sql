-- Use this query when you need one stable page from a larger result set.
-- Replace LIMIT and OFFSET with the page size and starting row you need.

SELECT
    item_id,
    item_name,
    created_at
FROM items
ORDER BY created_at DESC, item_id DESC
LIMIT 3 OFFSET 3;