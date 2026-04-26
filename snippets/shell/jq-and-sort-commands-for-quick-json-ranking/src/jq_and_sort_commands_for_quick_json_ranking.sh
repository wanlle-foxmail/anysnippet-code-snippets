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