# Select Latest Row per Group with ROW_NUMBER()

Return one newest status row per device with `ROW_NUMBER()` and a stable tie-breaker.

This snippet is useful when a device or entity keeps many history rows and you need exactly one latest status row for each group.

## Highlights

- Uses `ROW_NUMBER()` to rank rows inside each group
- Breaks timestamp ties with a secondary sort on the primary key
- Returns one deterministic result row per group

## Use Cases

- Fetch the latest device status for every device
- Return the newest order state from a status history table
- Select the most recent heartbeat per worker or service

## Code

```sql
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
```

## Notes

- SQLite window functions require SQLite 3.25 or newer.
- Add a stable secondary sort key, such as `event_id DESC`, whenever timestamps can collide.

## Verification

Run the integration tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- selecting the newest event for each device
- breaking ties with a secondary key
- single-row groups
- empty history tables
- stable ordering by group key
- ignoring older rows once a newer row exists

## Files

- `src/select_latest_row_per_group.sql`
- `tests/test_select_latest_row_per_group.py`
- `snippet.json`