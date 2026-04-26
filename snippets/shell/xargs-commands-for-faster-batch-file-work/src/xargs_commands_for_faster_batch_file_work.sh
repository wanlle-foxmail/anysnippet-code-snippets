find . -type f -name '*.log' -print0 | xargs -0 rm -f
# Delete matching files safely, even when file names contain spaces.

find . -type f -name '*.json' -print0 | xargs -0 grep -n 'TODO'
# Search across selected JSON files.

printf '%s\n' alpha beta gamma | xargs -n1 echo
# Run one command per input value.

printf '%s\n' https://example.com/a https://example.com/b https://example.com/c | xargs -n1 -P4 curl -O
# Download several files in parallel.

find . -type f -name '*.tmp' -print0 | xargs -0 -I{} mv '{}' ./tmp-backup/
# Move selected files into another directory.

find . -type f -name '*.jpg' -print0 | xargs -0 -n1 file
# Run one inspection command per file while keeping file names intact.