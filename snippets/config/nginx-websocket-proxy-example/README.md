# Nginx WebSocket Proxy Example

This snippet uses a Markdown guide with a Mermaid diagram and the official Nginx WebSocket proxying guide to explain a minimal reverse-proxy setup for upgraded connections.

## What This Snippet Covers

- One dedicated WebSocket location
- One explicit HTTP/1.1 proxy setting
- One explicit Upgrade plus Connection header pass-through

## Official Correctness Note

This snippet follows the official Nginx WebSocket proxying guide:

- https://nginx.org/en/docs/http/websocket.html

It also uses standard proxy module directives:

- https://nginx.org/en/docs/http/ngx_http_proxy_module.html

## Entry File

- `nginx-websocket-proxy.md`
- The entry file intentionally stays in Markdown so the request flow and the config fragment remain together on the first screen.

## How to Use It

- Open `nginx-websocket-proxy.md` first.
- Paste the example into the location that fronts your WebSocket service.
- Replace the sample upstream.
- Keep `proxy_http_version 1.1;` explicit so the config behaves consistently across mixed Nginx versions.
- Run `nginx -t`, then reload Nginx.

## Manual Verification

1. Open `nginx-websocket-proxy.md` and confirm the request-flow diagram matches your intended behavior.
2. Paste the config block into a staging Nginx configuration.
3. Replace the sample upstream with your real WebSocket service.
4. Run `nginx -t` on the staging host.
5. Reload Nginx.
6. Open a WebSocket connection through Nginx and confirm the upstream receives the protocol upgrade and keeps the connection open.

## Notes

- This snippet assumes the location handles only WebSocket traffic.
- If the upstream stays idle longer than the default timeout, increase `proxy_read_timeout` or configure application-level ping frames.