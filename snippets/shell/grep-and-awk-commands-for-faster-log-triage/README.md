Narrow noisy logs down to the fields that matter with `grep` and `awk` commands for error sampling, timeout checks, status-path extraction, and quick counts.

## What This Snippet Covers

- Sampling error lines down to a few key fields
- Keeping timeout rows with a visible timestamp prefix
- Extracting request paths from HTTP 500 lines
- Counting log lines by one level field
- Pulling values from key-value style logs
- Extracting duration fields from worker completion lines

## Before Using

- Replace the example log file names and field positions.
- Check a few raw lines first so the `awk` field numbers match your log format.
- Use these commands on copies or read-only logs when possible.

## Code

```sh
grep 'ERROR' app.log | awk '{print $1, $2, $NF}'
# Show the timestamp prefix and final field for error lines.

grep -i 'timeout' app.log | awk '{print $1, $2, $NF}'
# Keep a short timeout view with the timestamp prefix and final field.

grep ' 500 ' access.log | awk '{print $7}'
# Pull only the failing request path from HTTP 500 lines.

grep -E 'ERROR|WARN' app.log | awk '{count[$3]++} END {for (key in count) print key, count[key]}'
# Count how many log lines appear for each level field.

grep 'user_id=' app.log | awk -F'user_id=' '{split($2, parts, /[ \t]/); print parts[1]}'
# Extract only the user_id value from key-value style logs.

grep 'completed in' worker.log | awk '{print $(NF-1), $NF}'
# Pull the duration fields from job completion lines.
```

## Why These Commands Are Useful

- They help you shrink large logs into the exact fields needed for first-pass triage.
- They show a practical split of responsibilities: `grep` finds the lines, `awk` shapes the output.
- They work well when you want answers quickly without building a custom parser.

## Limitations

- This snippet stays `Draft` because it depends on local log formats and field positions.
- The `awk` field references need adjustment if your log structure differs.
- These commands are meant for quick triage, not full structured log parsing.

## Manual Verification

1. Confirm `grep` and `awk` are available in your shell.
2. Replace the example file names and field positions.
3. Run the commands on logs with known error, timeout, and status patterns.
4. Confirm the extracted fields match the original log lines.

## Files

- `src/grep_and_awk_commands_for_faster_log_triage.sh`
- `snippet.json`