# Nginx Static Cache Headers Example

This snippet uses a Markdown guide with a Mermaid diagram and the official Nginx headers module to explain a minimal static-asset cache policy.

## What This Snippet Covers

- One asset location with cache headers
- One `expires` directive for Expires plus max-age
- One extra `Cache-Control` hint for versioned assets

## Official Correctness Note

This snippet follows the official Nginx headers module documentation for `expires` and `add_header`:

- https://nginx.org/en/docs/http/ngx_http_headers_module.html

It also uses the standard `root` directive for the asset path:

- https://nginx.org/en/docs/http/ngx_http_core_module.html#root

## Entry File

- `nginx-static-cache-headers.md`
- The entry file intentionally stays in Markdown so the request flow and the config fragment remain together on the first screen.

## How to Use It

- Open `nginx-static-cache-headers.md` first.
- Paste the example into your active Nginx configuration.
- Replace the asset URL prefix and root path.
- Use hashed or versioned filenames before keeping `immutable`.
- Run `nginx -t`, then reload Nginx.

## Manual Verification

1. Open `nginx-static-cache-headers.md` and confirm the request-flow diagram matches your intended behavior.
2. Paste the config block into a staging Nginx configuration.
3. Replace the sample asset path and root path.
4. Run `nginx -t` on the staging host.
5. Reload Nginx.
6. Request one asset and confirm the response includes `Expires` and `Cache-Control` headers.

## Notes

- The official docs say `expires` already sets `Cache-Control: max-age=...`, so the extra `add_header` line is only for additional hints such as `public, immutable`.
- This pattern is safer for versioned static assets than for HTML responses that change frequently.