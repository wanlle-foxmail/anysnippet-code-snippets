# Nginx Static Files Example

This snippet uses a Markdown guide with a Mermaid diagram and the official Nginx `root` directive to explain a minimal static-files setup.

## What This Snippet Covers

- One static prefix location
- One disk path built from `root` plus the request URI
- One direct file-serving path with no upstream proxy

## Official Correctness Note

This snippet follows the Nginx beginner's guide and the official core module documentation for `root`:

- https://nginx.org/en/docs/beginners_guide.html
- https://nginx.org/en/docs/http/ngx_http_core_module.html#root

## Entry File

- `nginx-static-files.md`
- The entry file intentionally stays in Markdown so the request flow and the config fragment remain together on the first screen.

## How to Use It

- Open `nginx-static-files.md` first.
- Paste the example into your active Nginx configuration.
- Replace the URL prefix and root path with your real file layout.
- Put your static files under the mapped disk path.
- Run `nginx -t`, then reload Nginx.

## Manual Verification

1. Open `nginx-static-files.md` and confirm the request-flow diagram matches your intended behavior.
2. Paste the config block into a staging Nginx configuration.
3. Replace the sample URL prefix and root path.
4. Put a test file in the mapped directory.
5. Run `nginx -t` on the staging host.
6. Reload Nginx.
7. Request the test file and confirm Nginx serves it directly from disk.

## Notes

- This snippet intentionally uses `root`, not `alias`, so the request URI stays part of the resulting disk path.
- If you want to serve `/assets/logo.png` from a directory that does not include `/assets/` in its on-disk path, use an `alias`-focused snippet instead.