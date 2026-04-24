# Nginx Real IP Behind Proxy Example

This snippet uses a Markdown guide with a Mermaid diagram and the official Nginx real-IP module to explain a minimal trusted-proxy client-IP restore setup.

## What This Snippet Covers

- One trusted proxy CIDR list
- One forwarded-header source for the client IP
- One recursive lookup through the forwarded chain

## Official Correctness Note

This snippet follows the official Nginx real-IP module documentation:

- https://nginx.org/en/docs/http/ngx_http_realip_module.html

The proxied location itself uses the standard proxy module:

- https://nginx.org/en/docs/http/ngx_http_proxy_module.html

## Entry File

- `nginx-real-ip-behind-proxy.md`
- The entry file intentionally stays in Markdown so the request flow and the config fragment remain together on the first screen.

## How to Use It

- Open `nginx-real-ip-behind-proxy.md` first.
- Paste the example into your active Nginx configuration.
- Replace the trusted CIDR with the exact proxy or load balancer range you control.
- Confirm your package or build includes `ngx_http_realip_module`.
- Run `nginx -t`, then reload Nginx.

## Manual Verification

1. Open `nginx-real-ip-behind-proxy.md` and confirm the request-flow diagram matches your intended behavior.
2. Paste the config block into a staging Nginx configuration.
3. Replace the trusted CIDR and sample upstream.
4. Run `nginx -t` on the staging host.
5. Reload Nginx.
6. Send a request through the trusted proxy and confirm Nginx logs or the upstream app see the real client IP instead of the proxy IP.

## Notes

- This snippet is only safe when the trusted range is narrow and correct.
- If your deployment uses the PROXY protocol instead of forwarded headers, use the `proxy_protocol` form of `real_ip_header` instead.