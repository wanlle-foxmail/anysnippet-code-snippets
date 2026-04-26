Triage API responses faster with `curl` and `jq` commands for health checks, auth checks, failure extraction, item counts, and quick status summaries.

## What This Snippet Covers

- Pretty-printing a health response
- Pulling a single status field from JSON
- Counting returned items
- Filtering failed items down to the key fields
- Inspecting identity fields from an authenticated endpoint
- Grouping API results into a status summary

## Before Using

- Replace the example URLs, query strings, and bearer token placeholders.
- Confirm the JSON field names match your API responses.
- Use `-s` only when you do not need curl progress output.

## Code

```sh
curl -s https://api.example.com/health | jq '.'
# Pretty-print a health response before deeper debugging.

curl -s https://api.example.com/health | jq '.status'
# Pull out only the status field from a health endpoint.

curl -s https://api.example.com/items | jq '.items | length'
# Count how many items the API returned.

curl -s 'https://api.example.com/items?status=failed' | jq '.items[] | {id, error}'
# Keep only failed item IDs and error messages.

curl -s -H 'Authorization: Bearer YOUR_TOKEN' https://api.example.com/me | jq '{id, email, role}'
# Inspect the key identity fields from an authenticated API response.

curl -s https://api.example.com/items | jq '.items | group_by(.status) | map({status: .[0].status, count: length})'
# Turn raw items into a quick status breakdown.
```

## Why These Commands Are Useful

- They turn raw API JSON into answers you can use immediately during triage.
- They keep the request and the JSON filter together so copy-paste debugging stays fast.
- They cover the common questions people ask first: is it healthy, who am I, what failed, and how many are affected.

## Limitations

- This snippet stays `Draft` because it depends on real endpoints, authentication, and response shapes.
- Local validation can confirm the command and filter chain with sample JSON, but real APIs still need matching auth and field layouts.
- Grouping by status assumes the API returns a list under `.items`.
- Auth examples require a valid token and a compatible endpoint.

## Manual Verification

1. Confirm `curl` and `jq` are available in your shell.
2. Replace the example URLs, tokens, and JSON field names.
3. Run the commands against real endpoints with known response shapes.
4. Confirm the extracted fields and grouped counts match the raw JSON.

## Files

- `src/curl_and_jq_commands_for_faster_api_triage.sh`
- `snippet.json`