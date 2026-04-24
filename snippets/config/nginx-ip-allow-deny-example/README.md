# Nginx IP Allow Deny Example

This snippet uses a Markdown guide with a Mermaid diagram and the official Nginx access module to explain a small IP-based allowlist for one location.

## What This Snippet Covers

- One protected location
- One IPv4 allowlist example
- One IPv6 allowlist example with a default deny

## Official Correctness Note

This snippet follows the official Nginx access module documentation:

- https://nginx.org/en/docs/http/ngx_http_access_module.html

If Nginx is behind another trusted proxy, combine it with the official real-IP module first:

- https://nginx.org/en/docs/http/ngx_http_realip_module.html

## Entry File

- `nginx-ip-allow-deny.md`
- The entry file intentionally stays in Markdown so the access flow and the config fragment remain together on the first screen.

## How to Use It

- Open `nginx-ip-allow-deny.md` first.
- Paste the example into the location you want to protect.
- Replace the sample IP ranges.
- If the request source address is currently a proxy address, fix real-IP handling first.
- Run `nginx -t`, then reload Nginx.

## Manual Verification

1. Open `nginx-ip-allow-deny.md` and confirm the access-flow diagram matches your intended behavior.
2. Paste the config block into a staging Nginx configuration.
3. Replace the sample IP ranges.
4. Run `nginx -t` on the staging host.
5. Reload Nginx.
6. Confirm a client inside the allowed range can reach the location and a client outside the allowed range is denied.

## Notes

- This snippet is intentionally small and focuses on one allowlist pattern.
- The final `deny all;` is the line that turns the block into a closed allowlist.
- If Nginx is behind a proxy or load balancer, these rules are not trustworthy until real-IP handling is configured correctly.