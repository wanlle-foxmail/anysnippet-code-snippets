# Nginx IP Rate Limit Example

This snippet uses a Markdown guide with a Mermaid diagram and the official Nginx request-rate limiting directives to explain the smallest per-IP throttling setup.

## What This Snippet Covers

- One shared rate-limit zone keyed by client IP
- One protected `location /api/` block
- One short burst allowance before rejection

## Official Correctness Note

This snippet follows the official Nginx `ngx_http_limit_req_module` documentation:

- https://nginx.org/en/docs/http/ngx_http_limit_req_module.html

It also uses the official `ngx_http_limit_conn_module` documentation only to distinguish connection limits from request-rate limits:

- https://nginx.org/en/docs/http/ngx_http_limit_conn_module.html

## Entry File

- `nginx-ip-rate-limit.md`
- The entry file intentionally stays in Markdown so the request flow and the config fragment remain together on the first screen.

## How to Use It

- Open `nginx-ip-rate-limit.md` first.
- Paste the example into your active Nginx configuration.
- Keep `limit_req_zone` inside the `http {}` block.
- Adjust the path, upstream, rate, and burst for your endpoint.
- Run `nginx -t`, then reload Nginx.

## Manual Verification

1. Open `nginx-ip-rate-limit.md` and confirm the request-flow diagram matches your intended behavior.
2. Paste the config block into a staging Nginx configuration.
3. Replace the sample upstream and path with your real service.
4. Run `nginx -t` on the staging host.
5. Reload Nginx.
6. Send repeated requests from one IP and confirm requests above the average rate are delayed or rejected according to the burst setting.

## Notes

- The official docs use `limit_req_zone $binary_remote_addr ...` for per-IP request-rate limiting, and this snippet follows that shape.
- The official docs say `limit_req zone=... burst=...;` delays excessive requests until the burst is exhausted, then rejects them.
- If you do not want delayed excess requests, the official docs say to add `nodelay`.