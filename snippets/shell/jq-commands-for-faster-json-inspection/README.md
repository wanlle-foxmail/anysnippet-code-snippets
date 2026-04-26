Read messy JSON faster and pull out exactly what you need with `jq` commands for formatting, filtering, field selection, aggregation, and clean terminal output.

## What This Snippet Covers

- Pretty-printing JSON from a local file
- Extracting a single string field as plain text
- Filtering array entries by status
- Reshaping objects to keep selected fields only
- Counting items inside an array
- Emitting compact one-line JSON for pipelines
- Combining `curl` with `jq` to aggregate API response data

## Before Using

- Replace `response.json` with a real JSON file path.
- Replace `https://api.example.com/items` with a reachable API endpoint.
- Make sure the example field names match your JSON shape.

## Code

```sh
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
```

## Why These Commands Are Useful

- They cover the fast JSON checks developers do during API work and debugging.
- They help you move from raw JSON noise to one useful field, count, or filtered view.
- They show both file-based and live-API `jq` usage in one small snippet.

## Limitations

- This snippet stays `Draft` because it depends on placeholder files, endpoints, and JSON shapes.
- Some expressions assume array fields such as `.items` already exist in the input.
- The last example also requires `curl` and a reachable endpoint.

## Manual Verification

1. Confirm `jq --version` works.
2. Replace the placeholder file path and API URL.
3. Run the commands against real JSON that matches the example field names.
4. Confirm the filtered or aggregated output matches the expected result.

## Files

- `src/jq_commands_for_faster_json_inspection.sh`
- `snippet.json`