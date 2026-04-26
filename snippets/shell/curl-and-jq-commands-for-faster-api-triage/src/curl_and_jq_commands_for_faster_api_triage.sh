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