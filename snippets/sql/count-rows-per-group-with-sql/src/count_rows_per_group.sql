-- Use this query when you need a quick row count for each group value.
-- The ORDER BY clause puts larger groups first and keeps ties stable.

SELECT
    category_name,
    COUNT(*) AS row_count
FROM items
GROUP BY category_name
ORDER BY row_count DESC, category_name ASC;