Search codebases more precisely with `find` and `grep` commands for language-specific file selection, recursive keyword scans, and noise reduction.

## What This Snippet Covers

- Searching only one language or file type
- Skipping `.git` while scanning source files
- Finding heading patterns inside markdown files
- Searching YAML files for one configuration key
- Returning only files that contain one error word
- Searching several file types for common HTTP client calls

## Before Using

- Replace the example paths and patterns with real local values.
- Narrow the search root when working in very large repositories.
- Keep your file globs specific so the results stay focused.

## Code

```sh
find ./src -type f -name '*.ts' -exec grep -n 'TODO' {} +
# Search only TypeScript files for one keyword.

find . -name '.git' -prune -o -type f -name '*.py' -exec grep -n 'import requests' {} +
# Skip .git while searching Python files for one import pattern.

find . -type f -name '*.md' -exec grep -n '^## ' {} +
# Find markdown files that contain second-level headings.

find . -type f \( -name '*.yml' -o -name '*.yaml' \) -exec grep -n 'image:' {} +
# Search YAML files for image definitions.

find . -type f -name '*.log' -exec grep -l 'ERROR' {} +
# Return only the log files that contain one error word.

find . -type f \( -name '*.js' -o -name '*.ts' \) -exec grep -nE 'fetch|axios' {} +
# Search JavaScript and TypeScript files for common HTTP client calls.
```

## Why These Commands Are Useful

- They help you search the right files first instead of scanning the whole tree blindly.
- They keep language filters and grep patterns in the same command so the intent stays obvious.
- They reduce search noise in mixed-language repositories.

## Limitations

- This snippet stays `Draft` because it depends on local codebases and file layouts.
- Large repositories can still be slow to scan if the search root is too broad.
- When `grep` receives only one file, some environments print `line:content` without the file name.
- Results depend on the exact file names, globs, and text patterns used.

## Manual Verification

1. Confirm `find` and `grep` are available in your shell.
2. Replace the example roots, file globs, and search patterns.
3. Run the commands in a repository that contains known matches.
4. Confirm the returned files or lines match the intended filter rules, and note that single-file `grep` output may omit the file name.

## Files

- `src/find_and_grep_commands_for_targeted_code_search.sh`
- `snippet.json`