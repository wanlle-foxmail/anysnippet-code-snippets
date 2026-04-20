# Compute Order and Refund Metrics in SQLite

Aggregate order counts, paid revenue, and refunds for one reporting window in SQLite.

This snippet is useful when an admin view or reporting endpoint needs created-order counts, paid revenue, and refund activity for the same reporting window without issuing several separate queries.

## Highlights

- Counts order states in one query
- Sums paid and refunded amounts
- Tracks refunds by refunded_at

## Use Cases

- Show order counts in an operations dashboard
- Build daily revenue and refund summaries
- Feed a reporting API with one time-window query

## Code

```sql
-- Use this query when a dashboard needs order counts, paid revenue,
-- and refund totals for one reporting window.
-- Order metrics use created_at. Refund metrics use refunded_at.

WITH order_metrics AS (
    SELECT
        COUNT(*) AS total_orders,
        COUNT(CASE WHEN status = 'paid' THEN 1 END) AS paid_orders,
        COUNT(CASE WHEN status = 'pending' THEN 1 END) AS pending_orders,
        COUNT(CASE WHEN status = 'cancelled' THEN 1 END) AS cancelled_orders,
        COALESCE(SUM(CASE WHEN status = 'paid' THEN total_amount END), 0) AS paid_revenue
    FROM orders
    WHERE created_at >= :window_start AND created_at < :window_end
),
refund_metrics AS (
    SELECT
        COUNT(*) AS refunded_orders,
        COALESCE(SUM(total_amount), 0) AS refunded_amount
    FROM orders
    WHERE refunded_at >= :window_start AND refunded_at < :window_end
)
SELECT
    order_metrics.total_orders,
    order_metrics.paid_orders,
    order_metrics.pending_orders,
    order_metrics.cancelled_orders,
    refund_metrics.refunded_orders,
    order_metrics.paid_revenue,
    refund_metrics.refunded_amount
FROM order_metrics
CROSS JOIN refund_metrics;
```

## Notes

- `COUNT(CASE WHEN ...)` keeps the status counts compact and SQLite-friendly.
- The time filter is start-inclusive and end-exclusive.
- Order counts and paid revenue use `created_at`, while refund metrics use `refunded_at`.
- If you adapt this query to another database, update the parameter syntax for that driver or dialect.

## Verification

Run the integration tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- status counts in a single query
- paid revenue and refunded amount aggregation
- refund metrics based on the refund timestamp window
- empty-window zero values
- null amount handling
- boundary behavior for created and refund timestamps

## Files

- `src/build_dashboard_metrics.sql`
- `tests/test_build_dashboard_metrics.py`
- `snippet.json`