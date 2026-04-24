# Nginx HTTPS Redirect Example

This snippet uses a Markdown guide with a Mermaid diagram and the official Nginx `return` directive to explain the smallest HTTP-to-HTTPS redirect setup.

## What This Snippet Covers

- One HTTP listener on port 80
- One permanent redirect to the HTTPS URL
- Host and request URI preserved in the redirect target

## Official Correctness Note

This snippet follows the official Nginx `ngx_http_rewrite_module` documentation for `return` redirects:

- https://nginx.org/en/docs/http/ngx_http_rewrite_module.html

It also relies on standard `listen` and `server_name` directives from the Nginx core module:

- https://nginx.org/en/docs/http/ngx_http_core_module.html

## Entry File

- `nginx-https-redirect.md`
- The entry file intentionally stays in Markdown so the redirect flow and the config fragment remain together on the first screen.

## How to Use It

- Open `nginx-https-redirect.md` first.
- Paste the example into your active Nginx configuration.
- Replace the example host names.
- Make sure an HTTPS server already exists for the same hosts.
- Run `nginx -t`, then reload Nginx.

## Manual Verification

1. Open `nginx-https-redirect.md` and confirm the request-flow diagram matches your intended behavior.
2. Paste the config block into a staging Nginx configuration.
3. Replace the example host names with your real hosts.
4. Confirm the HTTPS server on port 443 is already working.
5. Run `nginx -t` on the staging host.
6. Reload Nginx.
7. Request the HTTP URL and confirm Nginx responds with a `301` redirect to the HTTPS URL.

## Notes

- This snippet only handles the redirecting HTTP server, not the TLS certificate or the HTTPS server block itself.
- The redirect preserves both the original host and the original request URI.