Rank JSON fields quickly with `jq` and `sort` commands for names, durations, timestamps, grouped pairs, and unique owner lists.

## What This Snippet Covers

- Alphabetizing names from a JSON array
- Sorting numeric durations
- Showing the slowest durations first
- Sorting two selected JSON fields together
- Ranking timestamps from newest to oldest
- Returning unique owner values

## Before Using

- Replace the example file name and field paths.
- Confirm the JSON shape includes an `.items` array.
- Use `jq -r` when you want clean text output for downstream sorting.

## Code

```sh
jq -r '.items[] | .name' response.json | sort
# Alphabetize item names from a JSON array.

jq -r '.items[] | .duration_ms' response.json | sort -n
# Sort numeric durations from low to high.

jq -r '.items[] | .duration_ms' response.json | sort -nr | head -10
# Show the ten slowest durations first.

jq -r '.items[] | [.status, .id] | @tsv' response.json | sort
# Sort status and ID pairs for quick scanning.

jq -r '.items[] | .updated_at' response.json | sort -r
# Rank timestamps from newest to oldest.

jq -r '.items[] | .owner' response.json | sort -u
# Show the unique owners that appear in the JSON array.
```

## Why These Commands Are Useful

- They make ranking and scanning JSON responses easier without opening a full scripting loop.
- They show when `jq` should extract clean text and when `sort` should take over.
- They work well for quick reports built from saved API responses.

## Limitations

- This snippet stays `Draft` because it depends on local JSON files and field paths.
- Numeric sorting only works correctly when the extracted field is a numeric value.
- Missing `.items` arrays or null fields require filter changes.

## Manual Verification

1. Confirm `jq` and `sort` are available in your shell.
2. Replace the example file name and field paths.
3. Run the commands on a JSON file with known names, timestamps, and numeric values.
4. Confirm the sorted output matches the source data.

## Files

- `src/jq_and_sort_commands_for_quick_json_ranking.sh`
- `snippet.json`