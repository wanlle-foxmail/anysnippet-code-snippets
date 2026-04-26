Clean and rank repeated text fast with `sort` and `uniq` commands for deduplication, duplicate checks, frequency reports, and field-based ordering.

## What This Snippet Covers

- Sorting plain text lines alphabetically
- Sorting and deduplicating in one step
- Removing repeated lines after sorting
- Showing only duplicated entries
- Ranking repeated values by frequency
- Sorting CSV-like rows by a chosen field

## Before Using

- Replace the example file names with real local files.
- Make sure delimiter-based sorts match the actual file format.
- Remember that `uniq` works on adjacent duplicates, so sorting usually comes first.

## Code

```sh
sort names.txt
# Sort lines alphabetically.

sort -u emails.txt
# Sort and drop duplicate lines in one step.

sort events.txt | uniq
# Collapse adjacent duplicates after sorting.

sort events.txt | uniq -d
# Show only lines that appear more than once.

sort events.txt | uniq -c | sort -nr
# Count repeated lines and rank the most common ones first.

sort -t',' -k2,2 report.csv
# Sort a CSV-like file by the second field.
```

## Why These Commands Are Useful

- They solve common cleanup and reporting tasks without needing a script or spreadsheet.
- They make duplicate-heavy text easier to inspect and rank.
- They show where `sort` ends and where `uniq` becomes useful in the pipeline.

## Limitations

- This snippet stays `Draft` because it depends on local files.
- Delimiter-based field sorting assumes the file really uses the chosen separator.
- Locale settings can change sorting results for some character sets.

## Manual Verification

1. Confirm `sort --version` or `sort -h` works in your environment.
2. Replace the example file names with real local files.
3. Run each command on text that contains known duplicates or sortable fields.
4. Confirm the output order, duplicates, or counts match the intended result.

## Files

- `src/sort_and_uniq_commands_for_quick_dedup_and_ranking.sh`
- `snippet.json`