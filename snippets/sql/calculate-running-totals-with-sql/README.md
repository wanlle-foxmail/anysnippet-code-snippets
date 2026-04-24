# Calculate Running Totals in SQLite

Return one running total per row in SQLite with `SUM() OVER (...)` and a stable tie-breaker.

This snippet is useful when an account, wallet, or ledger table stores many transactions and you need the cumulative balance after each row.

## Highlights

- Uses `SUM()` as a window function
- Keeps running totals deterministic
- Separates balances by account

## Use Cases

- Show account balances after each transaction
- Build a ledger view with cumulative totals
- Compute per-customer running balances in reports

## Code

```sql
-- Use this query when each account has many ledger rows and you need one
-- running balance per row in a stable order.
-- The secondary sort on transaction_id keeps the running total deterministic
-- when multiple rows share the same posted_at value.

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
```

## Notes

- SQLite window functions require SQLite 3.25 or newer.
- Add a stable secondary sort key such as `transaction_id ASC` when timestamps can match.
- The query returns rows grouped by `account_id` so each account's running total stays easy to scan.

## Verification

Run the integration tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- per-account running totals
- stable ordering for matching timestamps
- single-row accounts
- empty tables
- grouped output ordering
- negative amount handling

## Files

- `src/calculate_running_totals.sql`
- `tests/test_calculate_running_totals.py`
- `snippet.json`