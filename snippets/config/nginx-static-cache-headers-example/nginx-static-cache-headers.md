Use this guide when Nginx should send cache headers for versioned static assets.

## Request Flow

```mermaid
flowchart LR
    client[Client] --> nginx[Nginx location /assets/]
    nginx --> headers[Expires plus Cache-Control]
    headers --> cache[Browser or CDN cache]
```

## Minimal Example

```nginx
location /assets/ {
    root /srv/www;
    # Add Expires and Cache-Control max-age.
    expires 7d;
    # Add extra caching hints for versioned assets.
    add_header Cache-Control "public, immutable";
}
```

## Why This Is Correct

- The official `expires` directive adds or modifies both the `Expires` and `Cache-Control` response headers.
- The official docs say a positive `expires` value sets `Cache-Control: max-age=...`.
- The official `add_header` directive can add another response header field at the same location level.

## Before You Use It

- Use this pattern for versioned assets, not for frequently changing HTML pages.
- Replace `/assets/` and `/srv/www` with your real asset path.
- Adjust the cache duration for your deployment.
- Run `nginx -t`, then reload with `nginx -s reload`.

## Official References

- https://nginx.org/en/docs/http/ngx_http_headers_module.html
- https://nginx.org/en/docs/http/ngx_http_core_module.html#root