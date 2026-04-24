# Nginx Load Balancer Baseline

This snippet uses a Markdown guide with a Mermaid diagram and a minimal Nginx config example to explain the simplest round-robin load balancer setup.

## What This Snippet Covers

- One Nginx load balancer in front of two backend servers
- Default round-robin balancing through one `upstream` block
- Inline comments that explain the core directives

## Entry File

- `nginx-load-balancer.md`
- The entry file intentionally stays in Markdown so the diagram and the minimal config example remain together on the first screen.

## How to Use It

- Open `nginx-load-balancer.md` first.
- Replace the example backend addresses and port.
- Copy the config block into your active Nginx `http {}` context or a matching include file.
- Run `nginx -t` before reloading Nginx.

## Manual Verification

1. Open `nginx-load-balancer.md` and confirm the request-flow diagram matches your intended topology.
2. Replace the sample backend addresses with real upstream servers.
3. Paste the config block into a staging Nginx configuration.
4. Run `nginx -t` on the staging host.
5. Reload Nginx and send multiple requests to confirm both backends receive traffic.

## Notes

- Nginx uses round-robin by default when multiple `server` lines exist in one `upstream` block and no alternative balancing method is declared.
- This snippet intentionally omits forwarded headers, TLS, health checks, and sticky sessions so the first read stays small and clear.