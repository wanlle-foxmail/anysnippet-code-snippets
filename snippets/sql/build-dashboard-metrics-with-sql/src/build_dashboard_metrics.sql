-- Use this query when a dashboard needs order counts, paid revenue,
-- and refund totals for one reporting window.
-- Order metrics use created_at. Refund metrics use refunded_at.
-- Flow:
-- reporting window
--   |- created_at  -> order_metrics  -> counts + paid revenue
--   |- refunded_at -> refund_metrics -> refund count + refund amount
--   `- cross join both one-row summaries into one dashboard row

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