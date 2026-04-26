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