# Nginx SPA History Fallback Example

This snippet uses a Markdown guide with a Mermaid diagram and the official Nginx `try_files` directive to explain a minimal SPA route fallback.

## What This Snippet Covers

- One SPA catch-all location
- One `try_files` chain for assets first
- One fallback to `index.html` for client-side routes

## Official Correctness Note

This snippet follows the official Nginx core module documentation for `try_files` and `root`:

- https://nginx.org/en/docs/http/ngx_http_core_module.html#try_files
- https://nginx.org/en/docs/http/ngx_http_core_module.html#root

## Entry File

- `nginx-spa-history-fallback.md`
- The entry file intentionally stays in Markdown so the request flow and the config fragment remain together on the first screen.

## How to Use It

- Open `nginx-spa-history-fallback.md` first.
- Paste the example into your active Nginx configuration.
- Replace the document root with your real SPA build directory.
- Ensure the Nginx worker user and any host security policy can read that directory.
- Add more specific locations before this catch-all if your deployment also has API routes.
- Run `nginx -t`, then reload Nginx.

## Manual Verification

1. Open `nginx-spa-history-fallback.md` and confirm the request-flow diagram matches your intended behavior.
2. Paste the config block into a staging Nginx configuration.
3. Replace the sample document root with your real SPA build directory.
4. Run `nginx -t` on the staging host.
5. Reload Nginx.
6. Request a real asset file and confirm Nginx serves the asset directly.
7. Request a client-side route such as `/dashboard` and confirm Nginx serves `index.html`.

## Notes

- This snippet is intentionally scoped to a pure SPA catch-all and does not include separate API or upload locations.
- The fallback works because `try_files` internally redirects to `/index.html` when no matching file or directory exists.