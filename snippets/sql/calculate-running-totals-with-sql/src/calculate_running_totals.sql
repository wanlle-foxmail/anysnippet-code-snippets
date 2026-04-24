-- Use this query when each account has many ledger rows and you need one
-- running balance per row in a stable order.
-- The secondary sort on transaction_id keeps the running total deterministic
-- when multiple rows share the same posted_at value.
-- Flow:
-- ledger_entries
--   -> partition rows by account_id
--   -> order rows by posted_at ASC, transaction_id ASC
--   -> apply SUM(amount) over each ordered partition
--   -> return one running total per row

SELECT
    transaction_id,
    account_id,
    posted_at,
    amount,
    SUM(amount) OVER (
        PARTITION BY account_id
        ORDER BY posted_at ASC, transaction_id ASC
        ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW
    ) AS running_total
FROM ledger_entries
ORDER BY account_id ASC, posted_at ASC, transaction_id ASC;