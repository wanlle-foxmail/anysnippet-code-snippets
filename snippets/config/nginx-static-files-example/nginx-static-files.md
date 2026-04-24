Use this guide when Nginx should serve files directly from disk instead of proxying the request to an app server.

## Request Flow

```mermaid
flowchart LR
    client[Client] --> nginx[Nginx location /assets/]
    nginx --> disk[Disk path /srv/www/assets/...]
    disk --> file[Static file response]
```

## Minimal Example

```nginx
server {
    listen 80;
    server_name _;

    location /assets/ {
        # Nginx appends the request URI to this root path.
        root /srv/www;
    }
}
```

## Why This Is Correct

- The official docs say `root` builds a file path by adding the request URI to the configured root path.
- The Nginx beginner's guide uses the same `location` plus `root` shape for serving static content.
- A request such as `/assets/logo.png` maps to `/srv/www/assets/logo.png` in this example.

## Before You Use It

- Put the real files under the directory path implied by the request URI.
- Replace `/assets/` and `/srv/www` with your real path layout.
- Use `alias` instead if you need to strip the location prefix instead of preserving it.
- Run `nginx -t`, then reload with `nginx -s reload`.

## Official References

- https://nginx.org/en/docs/beginners_guide.html
- https://nginx.org/en/docs/http/ngx_http_core_module.html#root