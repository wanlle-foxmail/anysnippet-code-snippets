# Nginx CORS Preflight Example

This snippet uses a Markdown guide with a Mermaid diagram and standard Nginx header and rewrite directives to explain a small CORS preflight response pattern for one API origin.

## What This Snippet Covers

- One `/api/` location
- One `OPTIONS` preflight response with status `204`
- One explicit allowed origin for actual API responses

## Official Correctness Note

This snippet uses documented Nginx directives from these official modules:

- https://nginx.org/en/docs/http/ngx_http_headers_module.html
- https://nginx.org/en/docs/http/ngx_http_rewrite_module.html

## Entry File

- `nginx-cors-preflight.md`
- The entry file intentionally stays in Markdown so the request flow and the config fragment remain together on the first screen.

## How to Use It

- Open `nginx-cors-preflight.md` first.
- Paste the example into the API location that should answer preflight requests.
- Replace the sample frontend origin and adjust methods and headers.
- Keep the `if` block narrow and return immediately if you follow this pattern.
- If the upstream API already emits CORS headers, remove one side to avoid duplicate response headers.
- Keep the origin explicit when credentialed cross-site requests are involved.
- Run `nginx -t`, then reload Nginx.

## Manual Verification

1. Open `nginx-cors-preflight.md` and confirm the diagram matches your intended browser flow.
2. Paste the config block into a staging Nginx configuration.
3. Replace the sample frontend origin and sample upstream.
4. Run `nginx -t` on the staging host.
5. Reload Nginx.
6. Send an `OPTIONS` request to the protected API path and confirm the response is `204` with the expected CORS headers.
7. Send a normal API request and confirm the actual response still includes `Access-Control-Allow-Origin`.

## Notes

- This snippet keeps the allowed origin explicit instead of using `*`.
- It intentionally does not add credential headers because not every API needs them.
- `Access-Control-Max-Age 86400` caches the preflight for one day; reduce it if you need faster policy rollback.