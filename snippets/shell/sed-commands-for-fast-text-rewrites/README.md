Rewrite files and streams faster with `sed` commands for substitutions, whitespace cleanup, line filtering, range printing, and quick text reshaping.

## What This Snippet Covers

- Replacing the first match on each line
- Replacing all matches and saving a new file
- Removing comment lines from config text
- Trimming trailing whitespace
- Printing a fixed line range
- Printing only a marker-delimited section
- Removing blank lines from output

## Before Using

- Replace the example file names with real local files.
- Prefer redirecting to a new file when testing rewrite commands on important data.
- Check your patterns carefully before using them on large files.

## Code

```sh
sed 's/http:/https:/' config.txt
# Replace the first matching text on each line.

sed 's/ERROR/FAILED/g' app.log > cleaned.log
# Replace all matching text and save the result to a new file.

sed '/^#/d' config.txt
# Drop comment lines from a config file.

sed 's/[[:space:]]*$//' input.txt
# Trim trailing whitespace from each line.

sed -n '1,20p' server.log
# Print only a chosen line range.

sed -n '/BEGIN CONFIG/,/END CONFIG/p' app.txt
# Print only the lines between two marker patterns.

sed '/^$/d' input.txt
# Remove blank lines from the output.
```

## Why These Commands Are Useful

- They cover high-value `sed` tasks without forcing platform-specific in-place syntax.
- They help you reshape text streams quickly when a full script would be overkill.
- They keep common substitution and filtering patterns easy to copy and adapt.

## Limitations

- This snippet stays `Draft` because it depends on local files.
- Regular expressions can match more broadly than expected if the pattern is too loose.
- Redirecting to a new file is safer, but it also means the original file is not changed automatically.

## Manual Verification

1. Confirm `sed --version` or `sed -h` works in your environment.
2. Replace the example file names with real local files.
3. Run each command on sample text with known matches and line ranges.
4. Confirm the rewritten or filtered output matches the intended result.

## Files

- `src/sed_commands_for_fast_text_rewrites.sh`
- `snippet.json`