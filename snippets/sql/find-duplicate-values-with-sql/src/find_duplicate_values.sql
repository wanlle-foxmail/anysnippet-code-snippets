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