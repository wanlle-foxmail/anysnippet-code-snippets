curl -i https://api.example.com/health
# Show response headers and body in one request.

curl -s -o /dev/null -w '%{http_code}\n' https://api.example.com/health
# Print only the HTTP status code for a quick health check.

curl --get https://api.example.com/search \
  --data-urlencode 'q=hello world' \
  --data-urlencode 'page=2'
# Send query parameters safely, including spaces and punctuation.

curl -X POST https://api.example.com/items \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer YOUR_TOKEN' \
  -d '{"name":"demo","enabled":true}'
# Send a JSON POST request with a bearer token header.

curl -X POST https://api.example.com/upload \
  -H 'Authorization: Bearer YOUR_TOKEN' \
  -F 'file=@./report.pdf' \
  -F 'title=Quarterly Report'
# Upload a file with multipart form fields.

curl -L -o release.zip https://example.com/download/latest
# Follow redirects and save the final response to a file.

curl --retry 5 --retry-delay 2 --retry-connrefused https://api.example.com/data
# Retry on temporary network or connection failures.