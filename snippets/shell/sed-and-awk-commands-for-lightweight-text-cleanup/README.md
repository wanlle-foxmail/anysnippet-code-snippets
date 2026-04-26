Clean messy text files quickly with `sed` and `awk` commands for trimming, replacing, filtering, splitting, and small report-style transforms.

## What This Snippet Covers

- Trimming trailing whitespace and dropping blank lines
- Keeping active key names from a config file
- Replacing delimiters before selecting fields
- Filtering HTTP error rows from a log sample
- Numbering non-empty lines
- Stripping quotes before selecting CSV-like fields

## Before Using

- Replace the example file names with real local files.
- Check that the field positions in your text match the examples.
- Redirect output to a new file when you want to keep the original input untouched.

## Code

```sh
sed 's/[[:space:]]*$//' messy.txt | awk 'NF'
# Trim trailing whitespace and drop blank lines.

sed '/^#/d' .env | awk -F'=' '{print $1}'
# Keep only active environment variable names.

sed 's/,/ /g' users.csv | awk '{print $1, $3}'
# Replace commas with spaces and print selected fields.

sed -n '1,100p' access.log | awk '$9 >= 500 {print $1, $7, $9}'
# Inspect the first 100 log lines and keep only HTTP error rows.

sed '/^$/d' input.txt | awk '{print NR ":" $0}'
# Remove blank lines and number the remaining lines.

sed 's/"//g' report.csv | awk -F',' '{print $1 "," $4}'
# Strip double quotes before selecting two CSV-like fields.
```

## Why These Commands Are Useful

- They handle common text cleanup steps without needing a larger parsing script.
- They show where `sed` is good for shaping lines and where `awk` is better for fields.
- They work well for small one-off cleanups in exports, configs, and logs.

## Limitations

- This snippet stays `Draft` because it depends on local text files and field layouts.
- Field positions can shift if the input format is inconsistent.
- These examples target lightweight cleanup, not full CSV parsing edge cases.

## Manual Verification

1. Confirm `sed` and `awk` are available in your shell.
2. Replace the example file names with real local files.
3. Run the commands on small inputs with known spacing, delimiters, or field positions.
4. Confirm the cleaned or selected output matches the intended transformation.

## Files

- `src/sed_and_awk_commands_for_lightweight_text_cleanup.sh`
- `snippet.json`