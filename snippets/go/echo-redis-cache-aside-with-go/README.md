Serve a product detail request in Echo with a Redis cache-aside flow and a database fallback.

## What It Does

- Defines `GET /products/:id` in Echo
- Reads Redis first with the key pattern `product:{id}`
- Falls back to a database-style lookup when the cache misses or the cached JSON is invalid
- Writes the database result back to Redis after a successful lookup
- Returns `source: "cache"` or `source: "database"` so the request path is easy to see

## Product Example

```go
type Product struct {
    ID       string  `json:"id"`
    Name     string  `json:"name"`
    Price    float64 `json:"price"`
    Category string  `json:"category"`
}
```

The example uses an in-memory `sampleProducts` map as the database placeholder. Replace `FindProductByID` with your real database query when you wire it into a real service.

## Usage

Run the handler:

```bash
go run redis_cache_handler.go
```

Optional: start Redis locally to observe a real cache hit on the second request.

```bash
docker run --rm -p 6379:6379 redis:7
```

Request a product:

```bash
curl http://localhost:8080/products/p-100
```

First response with Redis running:

```json
{
  "message": "Product loaded from database and cached.",
  "source": "database",
  "product": {
    "id": "p-100",
    "name": "Mechanical Keyboard",
    "price": 129.99,
    "category": "hardware"
  }
}
```

Second response for the same product:

```json
{
  "message": "Product loaded from cache.",
  "source": "cache",
  "product": {
    "id": "p-100",
    "name": "Mechanical Keyboard",
    "price": 129.99,
    "category": "hardware"
  }
}
```

If Redis is not running, the handler still returns the database result and logs the cache read or write failure. That keeps the example runnable while still showing the cache-aside control flow.

## Verification

```bash
go mod tidy
go test -race ./...
```

Automated verification stubs `CacheGet`, `CacheSet`, and `FindProductByID`, so a live Redis server is not required for the test suite.