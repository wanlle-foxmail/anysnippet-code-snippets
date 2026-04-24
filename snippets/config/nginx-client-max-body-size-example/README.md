# Nginx Client Max Body Size Example

This snippet uses a Markdown guide with a Mermaid diagram and the official Nginx `client_max_body_size` directive to explain a minimal upload-size override.

## What This Snippet Covers

- One upload location with a larger request size limit
- One upstream application behind the upload path
- One clear `413` rejection path for oversized requests

## Official Correctness Note

This snippet follows the official Nginx core module documentation for `client_max_body_size`:

- https://nginx.org/en/docs/http/ngx_http_core_module.html#client_max_body_size

The reload note in this snippet also matches the Nginx beginner's guide:

- https://nginx.org/en/docs/beginners_guide.html

## Entry File

- `nginx-client-max-body-size.md`
- The entry file intentionally stays in Markdown so the request flow and the config fragment remain together on the first screen.

## How to Use It

- Open `nginx-client-max-body-size.md` first.
- Paste the example into your active Nginx configuration.
- Replace the sample upstream.
- Adjust the size limit for your real upload path.
- Run `nginx -t`, then reload Nginx.

## Manual Verification

1. Open `nginx-client-max-body-size.md` and confirm the request-flow diagram matches your intended behavior.
2. Paste the config block into a staging Nginx configuration.
3. Replace the sample upstream with your real upload service.
4. Run `nginx -t` on the staging host.
5. Reload Nginx.
6. Upload a file below the configured limit and confirm the request reaches the upstream service.
7. Upload a file above the configured limit and confirm Nginx returns `413`.

## Notes

- Setting `client_max_body_size 0;` disables the size check, but this snippet intentionally keeps an explicit limit.
- Browsers do not always display the `413` error cleanly, which is also noted in the official docs.