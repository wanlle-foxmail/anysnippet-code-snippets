jq . response.json
# Pretty-print a JSON file.

jq -r '.user.email' response.json
# Extract one string field without quotes.

jq '.items[] | select(.status == "failed")' response.json
# Filter array entries by a field value.

jq '[.items[] | {id, status}]' response.json
# Keep only the fields you want from each item.

jq '.items | length' response.json
# Count items in an array.

jq -c '.items[]' response.json
# Emit one compact JSON object per line for pipelines.

curl -s https://api.example.com/items | jq '.items | map(.duration_ms) | add'
# Pull JSON from an API and aggregate one numeric field.