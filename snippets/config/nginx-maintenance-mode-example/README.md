# Nginx Maintenance Mode Example

This snippet uses a Markdown guide with a Mermaid diagram and the official Nginx `return`, `error_page`, and `internal` directives to explain a minimal full-site maintenance page.

## What This Snippet Covers

- One global `503` response path
- One internal maintenance page location
- One exact file target for the maintenance page

## Official Correctness Note

This snippet follows the official Nginx core and rewrite module documentation for `error_page`, `internal`, and `return`:

- https://nginx.org/en/docs/http/ngx_http_core_module.html#error_page
- https://nginx.org/en/docs/http/ngx_http_core_module.html#internal
- https://nginx.org/en/docs/http/ngx_http_rewrite_module.html

## Entry File

- `nginx-maintenance-mode.md`
- The entry file intentionally stays in Markdown so the request flow and the config fragment remain together on the first screen.

## How to Use It

- Open `nginx-maintenance-mode.md` first.
- Paste the example into your active Nginx configuration.
- Replace the maintenance page root path if needed.
- Ensure the Nginx worker user and any host security policy can read that directory.
- Add any exceptions for health checks or ACME challenge paths before the catch-all location.
- Run `nginx -t`, then reload Nginx.

## Manual Verification

1. Open `nginx-maintenance-mode.md` and confirm the request-flow diagram matches your intended behavior.
2. Put a real `maintenance.html` file under the configured root path.
3. Paste the config block into a staging Nginx configuration.
4. Run `nginx -t` on the staging host.
5. Reload Nginx.
6. Request the site and confirm Nginx responds with status `503` and the maintenance page body.
7. Request `/maintenance.html` directly and confirm the `internal` location does not expose it as a normal public file.

## Notes

- This is an always-on maintenance mode example, not a toggleable flag-file pattern.
- If your deployment must keep health checks, ACME paths, or admin IPs available, add those exceptions before the catch-all location.