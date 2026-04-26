Find the exact lines you need with `grep` commands for recursive search, context output, exact matches, inverted matches, and multi-pattern filters.

## What This Snippet Covers

- Matching a word with line numbers
- Searching a directory tree recursively
- Case-insensitive search for messy logs
- Hiding comment lines from config files
- Showing nearby lines around an important match
- Checking several error words in one pass

## Before Using

- Replace the example file and directory paths with real local paths.
- Make sure the pattern you use matches the kind of text in the target files.
- Use quotes around patterns that contain spaces or shell metacharacters.

## Code

```sh
grep -n 'ERROR' app.log
# Show matching lines with line numbers.

grep -rn 'TODO' ./src
# Search recursively through a source directory.

grep -i 'warning' app.log
# Match text without caring about letter case.

grep -v '^#' .env.example
# Hide comment lines and keep active settings visible.

grep -C 2 'panic' server.log
# Show two lines of context before and after a match.

grep -E 'timeout|refused|reset' app.log
# Match any of several error words with one command.
```

## Why These Commands Are Useful

- They cover the most common search shapes used in logs, configs, and source code.
- They make it easier to move from a vague keyword to the exact lines you need.
- They keep a few high-value options close at hand instead of forcing a man page lookup.

## Limitations

- This snippet stays `Draft` because it depends on local files and directories.
- Regular expressions can match more broadly than expected if the pattern is too loose.
- Recursive searches can be slow in large folders with many generated files.

## Manual Verification

1. Confirm `grep --version` or `grep -V` works in your environment.
2. Replace the example paths with real local files and folders.
3. Run each command on sample text that contains known matches.
4. Confirm the returned lines and context match the intended search behavior.

## Files

- `src/grep_commands_for_faster_terminal_search.sh`
- `snippet.json`