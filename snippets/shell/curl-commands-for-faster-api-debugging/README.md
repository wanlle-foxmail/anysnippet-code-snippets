Speed up API debugging and request testing with copy-paste `curl` commands for health checks, query strings, JSON requests, uploads, downloads, and retries.

## What This Snippet Covers

- Health checks with full headers and body output
- Quick status-code-only checks for scripts and terminals
- Safe query string encoding for search endpoints
- JSON POST requests with bearer token headers
- Multipart file uploads
- Redirect-aware downloads to a local file
- Built-in retries for temporary failures

## Before Using

- Replace `https://api.example.com/...` URLs with your real endpoints.
- Replace `YOUR_TOKEN` with a valid bearer token when the endpoint requires authentication.
- Replace `./report.pdf` with a real local file path for the upload example.
- Run download examples from a writable directory.

## Code

```sh
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
```

## Why These Commands Are Useful

- They cover the most common terminal API tasks without requiring a custom script.
- They show the command shapes developers often need to reconstruct from memory.
- They keep the options close to the use case, so the right variant is easy to copy.

## Limitations

- This snippet stays `Draft` because it depends on placeholder URLs, tokens, and files.
- Retry behavior depends on the target server and the installed `curl` version.
- Upload and download commands require real endpoints and local file permissions.

## Manual Verification

1. Confirm `curl --version` works.
2. Replace the placeholder URLs, token, and file path.
3. Run each command against a reachable endpoint.
4. Confirm the response, upload, or downloaded file matches the intended result.

## Files

- `src/curl_commands_for_faster_api_debugging.sh`
- `snippet.json`