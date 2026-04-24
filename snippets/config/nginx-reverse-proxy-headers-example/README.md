# Nginx Reverse Proxy Headers Example

This snippet uses a Markdown guide with a Mermaid diagram and the official Nginx proxy module to explain a minimal forwarded-header set for reverse proxying.

## What This Snippet Covers

- One proxied location
- One preserved host header
- One forwarded client IP chain and request scheme

## Official Correctness Note

This snippet follows the official Nginx proxy module documentation for `proxy_set_header` and its embedded variables:

- https://nginx.org/en/docs/http/ngx_http_proxy_module.html

It also uses the standard `$scheme` variable from the core module:

- https://nginx.org/en/docs/http/ngx_http_core_module.html

## Entry File

- `nginx-reverse-proxy-headers.md`
- The entry file intentionally stays in Markdown so the request flow and the config fragment remain together on the first screen.

## How to Use It

- Open `nginx-reverse-proxy-headers.md` first.
- Paste the example into the location that proxies requests to your app.
- Replace the sample upstream.
- If Nginx is behind another trusted proxy, configure real-IP handling before relying on the forwarded headers.
- Run `nginx -t`, then reload Nginx.

## Manual Verification

1. Open `nginx-reverse-proxy-headers.md` and confirm the request-flow diagram matches your intended behavior.
2. Paste the config block into a staging Nginx configuration.
3. Replace the sample upstream with your real app.
4. Run `nginx -t` on the staging host.
5. Reload Nginx.
6. Send a request through Nginx and confirm the upstream app receives the expected host, client IP, forwarded chain, and scheme headers.

## Notes

- This snippet is intentionally small and does not include every possible proxy header.
- `X-Forwarded-For` is safest when combined with trusted-proxy real-IP handling in front of the app.