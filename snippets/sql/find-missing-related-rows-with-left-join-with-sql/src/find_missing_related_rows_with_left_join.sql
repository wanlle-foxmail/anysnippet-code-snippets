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