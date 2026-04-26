find . -type f -name '*.log'
# Find files by name pattern.

find . -type f -mtime -7
# Find files changed within the last seven days.

find . -type f -size +100M
# Find files larger than 100 megabytes.

find . -type f -empty
# Find empty files.

find . -name '.git' -prune -o -type f -name '*.py' -print
# Skip .git directories while searching for Python files.

find ./downloads -type f -exec ls -lh {} \;
# Run one inspection command for each matched file.