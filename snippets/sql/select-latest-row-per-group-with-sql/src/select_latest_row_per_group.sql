-- Use this query when each device has many history rows and you need one
-- latest status row per device.
-- The secondary sort on event_id keeps results deterministic when
-- recorded_at values match.

WITH ranked_events AS (
    SELECT
        event_id,
        device_id,
        status,
        recorded_at,
        ROW_NUMBER() OVER (
            PARTITION BY device_id
            ORDER BY recorded_at DESC, event_id DESC
        ) AS row_rank
    FROM device_status_events
)
SELECT
    event_id,
    device_id,
    status,
    recorded_at
FROM ranked_events
WHERE row_rank = 1
ORDER BY device_id ASC;