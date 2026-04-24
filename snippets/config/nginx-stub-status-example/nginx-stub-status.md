Use this guide when you want a small Nginx status endpoint that stays reachable only from the local host.

## Request Flow

```mermaid
flowchart LR
    localhost[Local curl or monitor] --> endpoint[/nginx_status]
    endpoint --> status[stub_status metrics]
    remote[Remote client] --> blocked[Denied by Nginx]
```

## Minimal Example

```nginx
server {
    listen 127.0.0.1:8080;
    listen [::1]:8080;
    server_name localhost;

    location = /nginx_status {
        stub_status;
        allow 127.0.0.1;
        allow ::1;
        deny all;
    }
}
```

## Why This Is Correct

- The official stub status module docs show `stub_status;` inside a location block.
- The official docs say the module exposes basic connection and request counters.
- The example keeps the endpoint local-only by binding the server to loopback and by denying every non-loopback address.

## Before You Use It

- Confirm your Nginx build includes `ngx_http_stub_status_module`.
- One quick check is `nginx -V 2>&1 | grep http_stub_status_module`.
- Keep this endpoint off the public Internet unless you have an explicit reason to expose it more broadly.
- Test locally with `curl http://127.0.0.1:8080/nginx_status`.
- Run `nginx -t`, then reload with `nginx -s reload`.

## Official References

- https://nginx.org/en/docs/http/ngx_http_stub_status_module.html
- https://nginx.org/en/docs/http/ngx_http_access_module.html