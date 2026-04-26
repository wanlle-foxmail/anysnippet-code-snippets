Reshape plain-text columns quickly with `cut` and `paste` commands for field extraction, character slicing, delimiter handling, and lightweight file merging.

## What This Snippet Covers

- Pulling one field from a delimited file
- Keeping several selected fields
- Extracting a fixed character range
- Splitting key names from key-value lines
- Combining two files column by column
- Merging two files with a custom delimiter

## Before Using

- Replace the example file names with real local files.
- Make sure the delimiter in the command matches the actual file format.
- For `paste`, use files with aligned line counts when you want rows to match cleanly.

## Code

```sh
cut -d',' -f1 users.csv
# Extract the first comma-separated field.

cut -d',' -f1,3 users.csv
# Keep only selected comma-separated fields.

cut -c1-12 access.log
# Slice a fixed character range from each line.

cut -d'=' -f1 .env
# Keep only the key name from key-value lines.

paste first_names.txt last_names.txt
# Merge two files side by side with tab separators.

paste -d',' ids.txt scores.txt
# Merge two files side by side with a custom delimiter.
```

## Why These Commands Are Useful

- They give you quick column reshaping without switching to a spreadsheet or script.
- They work well for small exports, config files, and quick one-off transformations.
- They show where field-based extraction ends and simple file merging begins.

## Limitations

- This snippet stays `Draft` because it depends on local files.
- `cut` works best on consistently delimited or fixed-width input.
- `paste` combines lines by position, so mismatched file lengths can produce uneven output.

## Manual Verification

1. Confirm `cut` and `paste` are available in your shell.
2. Replace the example file names with real local files.
3. Run the commands on small sample files with known columns.
4. Confirm the extracted or merged output matches the intended layout.

## Files

- `src/cut_and_paste_commands_for_fast_column_shaping.sh`
- `snippet.json`