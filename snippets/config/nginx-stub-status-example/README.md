# Nginx Stub Status Example

This snippet uses a Markdown guide with a Mermaid diagram and the official Nginx stub status module to explain a local-only status endpoint.

## What This Snippet Covers

- One loopback-only server example
- One `/nginx_status` endpoint
- One deny-by-default access pattern

## Official Correctness Note

This snippet follows the official Nginx stub status module documentation:

- https://nginx.org/en/docs/http/ngx_http_stub_status_module.html

It also uses the standard access module for the local-only restriction:

- https://nginx.org/en/docs/http/ngx_http_access_module.html

## Entry File

- `nginx-stub-status.md`
- The entry file intentionally stays in Markdown so the access model and the config fragment remain together on the first screen.

## How to Use It

- Open `nginx-stub-status.md` first.
- Paste the example into a staging Nginx configuration.
- Confirm your package or build includes `ngx_http_stub_status_module`.
- One quick check is `nginx -V 2>&1 | grep http_stub_status_module`.
- Keep the endpoint bound to loopback unless you have a separate access control plan.
- Run `nginx -t`, then reload Nginx.

## Manual Verification

1. Open `nginx-stub-status.md` and confirm the diagram matches your intended exposure model.
2. Paste the config block into a staging Nginx configuration.
3. Run `nginx -t` on the staging host.
4. Reload Nginx.
5. Run `curl http://127.0.0.1:8080/nginx_status` from the host and confirm the status page returns connection counters.
6. Confirm remote access is denied.

## Notes

- This module is not built by default.
- A local-only binding is the safest default for a simple status page snippet.
- If the module is missing, Nginx will reject the configuration when you test or reload it.