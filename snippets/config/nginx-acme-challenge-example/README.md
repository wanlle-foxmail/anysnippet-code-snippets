# Nginx ACME Challenge Example

This snippet uses a Markdown guide with a Mermaid diagram and standard Nginx `location`, `alias`, and `return` directives to explain a minimal ACME HTTP-01 path exception.

## What This Snippet Covers

- One HTTP challenge path that serves token files
- One redirect for all other HTTP traffic
- One isolated directory for ACME challenge files

## Official Correctness Note

This snippet uses standard Nginx directives documented in the core and rewrite modules:

- https://nginx.org/en/docs/http/ngx_http_core_module.html#location
- https://nginx.org/en/docs/http/ngx_http_core_module.html#alias
- https://nginx.org/en/docs/http/ngx_http_rewrite_module.html

## Entry File

- `nginx-acme-challenge.md`
- The entry file intentionally stays in Markdown so the request flow and the config fragment remain together on the first screen.

## How to Use It

- Open `nginx-acme-challenge.md` first.
- Paste the example into the HTTP server that handles port 80.
- Replace the sample host names.
- Point the `alias` directory to the place where your ACME client writes challenge tokens.
- Ensure the Nginx worker user and any host security policy can read that directory.
- Run `nginx -t`, then reload Nginx.

## Manual Verification

1. Open `nginx-acme-challenge.md` and confirm the request-flow diagram matches your intended behavior.
2. Paste the config block into a staging Nginx configuration.
3. Replace the sample host names and ACME directory.
4. Put a test token file into the challenge directory.
5. Run `nginx -t` on the staging host.
6. Reload Nginx.
7. Request the test challenge path and confirm Nginx serves the token file without redirecting it.
8. Request any other HTTP path and confirm Nginx redirects it to HTTPS.

## Notes

- This snippet is about routing only. It does not configure certificate issuance or the HTTPS server itself.
- The `alias` form keeps the on-disk challenge file path shorter than the request URI path.