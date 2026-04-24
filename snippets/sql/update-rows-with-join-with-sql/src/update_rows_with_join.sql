-- Use this query when one table stores rows and another table stores
-- the default value for each category.
-- The WHERE clause keeps existing values unchanged.

UPDATE items AS i
JOIN category_defaults AS d
    ON d.category_name = i.category_name
SET i.display_rank = d.display_rank
WHERE i.display_rank IS NULL;