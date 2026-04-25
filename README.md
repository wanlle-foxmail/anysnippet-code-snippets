# AnySnippet Code Snippets

![License: MIT](https://img.shields.io/badge/license-MIT-green.svg)
![Verification Required](https://img.shields.io/badge/verification-required-blue.svg)
![Verified Snippets](https://img.shields.io/badge/snippets-verified-orange.svg)
![Languages](https://img.shields.io/badge/languages-Python%20%7C%20SQL%20%7C%20Go%20%7C%20Config%20%7C%20Shell-blue.svg)

Official open-source code snippets, examples, and verification notes for AnySnippet Market.

This repository contains reusable code snippets for AnySnippet Market. It currently includes Python, Go, SQL, configuration, and shell snippets and is designed to expand to TypeScript, Rust, and more.

Every snippet is published with source code, metadata, and verification notes so developers can inspect quality before using it. Most verified snippets also include automated tests.

If you found this repository while searching for Python code snippets, SQL query patterns, Go backend snippets, Redis configuration snippets, Rust web examples, or reusable developer utilities, start with the snippet index below.

## Browse AnySnippet

Want to do more than copy and paste?

These snippets are published here as open-source references. If you want to browse, save, organize, and sync them across devices, use AnySnippet Market and import the published version into your workspace.


- [Browse in AnySnippet Market](https://www.anysnippet.com/market)
- [Learn more about AnySnippet](https://www.anysnippet.com/market)
- [Open a contribution proposal](https://github.com/wanlle-foxmail/anysnippet-code-snippets/issues/new)
- [Report an issue](https://github.com/wanlle-foxmail/anysnippet-code-snippets/issues)

![AnySnippet](/images/AnySnippet.png)

## Why this repository exists

- Show exactly what ships in AnySnippet Market
- Make snippet quality transparent through code, tests, and verification
- Give developers a trusted source of reusable code examples
- Create a clean path for community contributions
- Turn useful snippets into discoverable open-source entry points for AnySnippet

## What makes these snippets different

- Official snippets curated for AnySnippet Market
- Focused on real developer jobs, not filler examples
- Tests or deterministic verification are required
- Metadata is reviewable and structured
- Boundaries and known limitations are documented
- Contributions happen in public

## Popular languages and categories

### Languages

- Python code snippets for automation, file operations, scripting, and utilities
- Configuration snippets for Redis and deployment baselines
- Shell snippets for CLI workflows and reusable command recipes
- SQL snippets for query patterns, reporting techniques, and database workflows
- Go backend snippets for service handlers and infrastructure patterns
- Planned expansion for TypeScript, Rust, and other practical developer workflows

### Categories

- Backend
- Frontend
- Automation
- DevOps
- Testing
- Data processing
- Developer productivity

## Snippet Index

The table below is a human-friendly index for discovery. Each snippet should have its own folder with source code, tests or deterministic verification, metadata, and a local README.

| Snippet | Language | Description | Path | Verification | Status |
| --- | --- | --- | --- | --- | --- |
| Calculate Large File Hashes with Python | Python | Hash a large file in chunks with SHA-256 or MD5. | `snippets/python/calculate-large-file-hashes-with-python/` | `unittest` | Verified |
| Find Files by Extension with Python | Python | Walk through a directory and collect files whose extensions match a target list. | `snippets/python/find-files-by-extension-with-python/` | `unittest` | Verified |
| Calculate Directory Size with Python | Python | Walk a directory tree and return total bytes, file count, and subdirectory count. | `snippets/python/calculate-directory-size-with-python/` | `unittest` | Verified |
| Batch Records by Total Byte Size with Python | Python | Group UTF-8 string records into batches whose total encoded byte size stays under a limit. | `snippets/python/batch-records-by-total-byte-size-with-python/` | `unittest` | Verified |
| Compress Text for Storage with Python | Python | Compress structured text into a base64-encoded gzip string so it uses less space in text-only storage. | `snippets/python/compress-text-for-storage-with-python/` | `unittest` | Verified |
| Deduplicate Large JSONL Records by Key with Python | Python | Read a large JSONL file and write only the first record for each unique top-level string key. | `snippets/python/deduplicate-large-jsonl-records-by-key-with-python/` | `unittest` | Verified |
| Process Large CSV Files in Chunks with pandas | Python | Read a large CSV with pandas in chunks and return one processed result per chunk. | `snippets/python/process-large-csv-files-in-chunks-with-python/` | `unittest` | Verified |
| Read Headerless CSV in Chunks with pandas | Python | Read a headerless CSV file in chunks and yield row dictionaries. | `snippets/python/read-headerless-csv-in-chunks-with-pandas/` | `unittest` | Verified |
| Read Parquet Files as Records with pandas | Python | Load a Parquet file with pandas and return Python row dictionaries. | `snippets/python/read-parquet-files-as-records-with-pandas/` | `unittest` | Verified |
| Limit I/O Concurrency in Python | Python | Run I/O-bound tasks with ThreadPoolExecutor and a fixed worker limit. | `snippets/python/limit-concurrent-tasks-in-python/` | `unittest` | Verified |
| Graceful Worker Shutdown in Python | Python | Run a queue-based worker that finishes accepted work, rejects new submissions during shutdown, and runs cleanup once. | `snippets/python/graceful-shutdown-worker-with-python/` | `unittest` | Verified |
| Merge Sorted Files Lazily with Python | Python | Lazily merge multiple sorted UTF-8 text files into one sorted output stream. | `snippets/python/merge-sorted-files-lazily-with-python/` | `unittest` | Verified |
| Merge Overlapping Time Ranges with Python | Python | Sort numeric time ranges by start time and merge overlapping intervals. | `snippets/python/merge-overlapping-time-ranges-with-python/` | `unittest` | Verified |
| Refactor If Elif Dispatch to a Handler Map in Python | Python | Replace a long if and elif dispatch chain with a handler mapping for email, sms, and webhook notifications. | `snippets/python/refactor-large-if-elif-dispatch-in-python/` | `unittest` | Verified |
| Refactor Nested If Else to a Config Map in Python | Python | Replace nested if else decision logic with a config mapping for support ticket routing. | `snippets/python/replace-nested-if-else-with-config-mapping-in-python/` | `unittest` | Verified |
| Sample Large Streams with Reservoir Sampling in Python | Python | Sample a fixed number of items from a large stream without loading it all into memory. | `snippets/python/sample-large-streams-with-reservoir-sampling-in-python/` | `unittest` | Verified |
| Download Large File with Python | Python | Stream-download a large file with requests, compute MD5 on the fly, and clean up on failure. | `snippets/python/download-large-file-with-python/` | `unittest` | Verified |
| Split Large CSV Files by Row Count with Python | Python | Split a large CSV file into smaller UTF-8 CSV files with a fixed number of data rows per part. | `snippets/python/split-large-csv-files-by-row-count-with-python/` | `unittest` | Verified |
| Write Files Atomically with Python | Python | Write UTF-8 text to a temporary file in the target directory, then replace the target file atomically. | `snippets/python/write-files-atomically-with-python/` | `unittest` | Verified |
| Load Typed Environment Settings with Python | Python | Load APP_ENV, PORT, and DEBUG from environment variables with string, integer, and boolean parsing. | `snippets/python/load-typed-env-settings-with-python/` | `unittest` | Verified |
| Cache Function Results with TTL in Python | Python | Cache successful function results in memory until a short TTL expires. | `snippets/python/cache-function-results-with-ttl-in-python/` | `unittest` | Verified |
| Retry Failed Operations with Python | Python | Retry one callable on retryable errors with a fixed delay between attempts. | `snippets/python/retry-failed-operations-with-python/` | `unittest` | Verified |
| TDD Red Green Refactor with Python | Python | Practice a simple RED GREEN REFACTOR loop in Python with a tiny shipping fee function and executable tests. | `snippets/python/tdd-red-green-refactor-with-python/` | `unittest` | Verified |
| Read Large JSONL Files with Python | Python | Read a JSON Lines file lazily and yield one parsed value per non-empty line. | `snippets/python/read-large-jsonl-files-with-python/` | `unittest` | Verified |
| Read Gzipped JSONL Files with Python | Python | Read a gzipped JSON Lines file lazily and yield one parsed value per non-empty line. | `snippets/python/read-gzipped-jsonl-files-with-python/` | `unittest` | Verified |
| Read Large JSON Arrays with ijson | Python | Stream a top-level JSON array with ijson and yield one parsed item at a time. | `snippets/python/read-large-json-arrays-with-ijson/` | `unittest` | Verified |
| Write JSON Lines with Python | Python | Write Python values to a UTF-8 JSON Lines file and replace the output file atomically. | `snippets/python/write-json-lines-with-python/` | `unittest` | Verified |
| Deduplicate POST Requests with Idempotency Keys in FastAPI | Python | Deduplicate one FastAPI POST endpoint with an Idempotency-Key header and a single-process in-memory store. | `snippets/python/deduplicate-post-requests-with-idempotency-keys-in-fastapi-and-python/` | `unittest` | Verified |
| Pick Allowed Update Fields from a PATCH Body with Python | Python | Keep only allowed top-level fields from a PATCH-style request body in Python. | `snippets/python/pick-allowed-update-fields-from-a-patch-body-with-python/` | `unittest` | Verified |
| Parse Pagination Query Params Safely with Python | Python | Parse page and page_size query parameters with defaults, validation, and a fixed maximum in Python. | `snippets/python/parse-pagination-query-params-safely-with-python/` | `unittest` | Verified |
| Parse Sort Query Params with an Allowlist in Python | Python | Parse a comma-separated sort query parameter against a fixed allowlist of field names in Python. | `snippets/python/parse-sort-query-params-with-an-allowlist-in-python/` | `unittest` | Verified |
| Sort Tasks by Dependency Order with Python | Python | Return a stable task order that satisfies declared task dependencies. | `snippets/python/sort-tasks-by-dependency-order-with-python/` | `unittest` | Verified |
| Distribute URLs Evenly by Domain with Python | Python | Reorder crawl URLs so requests stay spread across domains instead of clustering on one site. | `snippets/python/distribute-urls-evenly-by-domain-with-python/` | `unittest` | Verified |
| FastAPI Health Check Endpoint with Python | Python | A /health endpoint that checks database and Redis, then returns aggregated service status. | `snippets/python/fastapi-health-check-endpoint-with-python/` | `unittest` | Verified |
| Readiness Check Endpoint with Dependency Gates in FastAPI | Python | Gate a FastAPI readiness endpoint on required dependencies and return 503 until the service is ready. | `snippets/python/readiness-check-endpoint-with-dependency-gates-in-fastapi-and-python/` | `unittest` | Verified |
| Validate Multipart Uploads with FastAPI | Python | Validate one multipart upload in FastAPI with filename, content type, and file size checks. | `snippets/python/validate-multipart-uploads-with-fastapi-and-python/` | `unittest` | Verified |
| Verify JWT Bearer Tokens with FastAPI | Python | Protect one FastAPI route by validating HS256 bearer tokens and exposing claims to the handler. | `snippets/python/verify-jwt-bearer-tokens-with-fastapi-and-python/` | `unittest` | Verified |
| Echo Health Check Endpoint with Go | Go | A /health endpoint that checks database, Redis, and disk space with uptime reporting. | `snippets/go/echo-health-check-endpoint-with-go/` | `go test` | Verified |
| Echo Redis Cache-Aside with Go | Go | Serve a product detail request in Echo with a Redis cache-aside flow and a database fallback. | `snippets/go/echo-redis-cache-aside-with-go/` | `go test` | Verified |
| Echo POST Validation with Go | Go | Validate a Book JSON body in Echo with go-playground/validator and return friendly field messages. | `snippets/go/validate-post-body-with-echo-and-validator/` | `go test` | Verified |
| Load Typed Environment Settings with Go | Go | Load APP_ENV, PORT, and DEBUG from environment variables with string, integer, and boolean parsing in Go. | `snippets/go/load-typed-environment-settings-with-go/` | `go test` | Verified |
| Reject Unknown JSON Fields with Go | Go | Decode one JSON body in Go while rejecting fields that are not declared on the target struct. | `snippets/go/reject-unknown-json-fields-with-go/` | `go test` | Verified |
| RBAC Middleware with Role-to-Route Mapping in Echo | Go | Protect Echo routes with a fixed role-to-route mapping. | `snippets/go/rbac-middleware-with-role-to-route-mapping-in-echo/` | `go test` | Verified |
| Validate Multipart Uploads with Echo | Go | Validate one multipart upload in Echo with filename, content type, and file size checks. | `snippets/go/validate-multipart-uploads-with-echo/` | `go test` | Verified |
| Request ID Middleware with Echo | Go | Ensure every Echo request has an X-Request-ID header and a matching context value. | `snippets/go/request-id-middleware-with-echo/` | `go test` | Verified |
| Limit Request Rate with Token Bucket in Echo | Go | Limit one Echo route per client IP with an in-memory token bucket and 429 responses. | `snippets/go/limit-request-rate-with-token-bucket-in-echo/` | `go test -race` | Verified |
| Verify JWT Bearer Tokens with Echo | Go | Protect one Echo route by validating HS256 bearer tokens and exposing claims to the handler. | `snippets/go/verify-jwt-bearer-tokens-with-echo/` | `go test -race` | Verified |
| Enforce Request Timeouts with Echo | Go | Cap one Echo request with a fixed deadline and return 503 when the handler runs too long. | `snippets/go/enforce-request-timeouts-with-echo/` | `go test -race` | Verified |
| Log HTTP Requests with Request ID in Echo | Go | Log one Echo request with method, path, status, duration, client IP, and request ID. | `snippets/go/log-http-requests-with-request-id-in-echo/` | `go test -race` | Verified |
| Deduplicate POST Requests with Idempotency Keys in Echo | Go | Deduplicate one Echo POST endpoint with an Idempotency-Key header and a single-process in-memory store. | `snippets/go/deduplicate-post-requests-with-idempotency-keys-in-echo/` | `go test -race` | Verified |
| Return Problem JSON Errors with Echo | Go | Return application/problem+json errors from Echo with one small error handler and one reusable writer helper. | `snippets/go/return-problem-json-errors-with-echo/` | `go test -race` | Verified |
| Graceful HTTP Server Shutdown with Go | Go | Run an HTTP server until a context is canceled, then finish in-flight requests within a shutdown timeout. | `snippets/go/graceful-http-server-shutdown-with-go/` | `go test -race` | Verified |
| Retry HTTP GET Requests with Go | Go | Retry an HTTP GET request on transport errors, 429 responses, and 5xx responses with a fixed delay. | `snippets/go/retry-http-requests-with-go/` | `go test -race` | Verified |
| Cache HTTP Responses with ETag in Go | Go | Cache one JSON HTTP response in Go with an ETag header and 304 Not Modified handling. | `snippets/go/cache-http-responses-with-etag-in-go/` | `go test` | Verified |
| Acquire Redis Locks with Go | Go | Acquire and release a Redis lock in Go with SET NX and a compare-and-delete release step. | `snippets/go/acquire-redis-locks-with-go/` | `go test -race` | Verified |
| Stream Server-Sent Events with Go | Go | Stream text/event-stream responses from Go with one message channel and per-message flushes. | `snippets/go/stream-server-sent-events-with-go/` | `go test -race` | Verified |
| Stream JSONL HTTP Responses with Go | Go | Stream JSONL HTTP responses from Go with one record channel and per-record flushes. | `snippets/go/stream-jsonl-http-responses-with-go/` | `go test -race` | Verified |
| Verify Webhook Signatures with Go | Go | Verify a webhook body against a sha256 HMAC signature header in Go. | `snippets/go/verify-webhook-signatures-with-go/` | `go test -race` | Verified |
| Inject a Clock Interface for Testable Time Logic in Go | Go | Inject a tiny clock interface into Go time checks so tests can control the current time without sleeping. | `snippets/go/inject-clock-interface-for-testable-time-logic-in-go/` | `go test -race` | Verified |
| Fan-In Channels with select in Go | Go | Combine two Go channels into one output stream with select and optional context cancellation. | `snippets/go/fan-in-channels-with-select-in-go/` | `go test -race` | Verified |
| Limit Concurrent Work with a Semaphore Channel in Go | Go | Limit concurrent Go work with a buffered channel semaphore and explicit acquire and release helpers. | `snippets/go/limit-concurrent-work-with-semaphore-channel-in-go/` | `go test -race` | Verified |
| Guard a Map with RWMutex in Go | Go | Wrap a Go map with sync.RWMutex so concurrent reads and writes stay race-free. | `snippets/go/guard-map-with-rwmutex-in-go/` | `go test -race` | Verified |
| Cancel Long-Running Work with Context in Go | Go | Run batch work in Go while checking context cancellation between units of work. | `snippets/go/cancel-long-running-work-with-context-in-go/` | `go test -race` | Verified |
| Reset Idle Timeouts with time.Timer in Go | Go | Reset a time.Timer on activity so Go code can detect idle timeouts without recreating timers. | `snippets/go/reset-idle-timeouts-with-time-timer-in-go/` | `go test -race` | Verified |
| Coalesce Bursty Events with select in Go | Go | Merge rapid Go events into fewer downstream batches with select, one timer, and a quiet window. | `snippets/go/coalesce-bursty-events-with-select-in-go/` | `go test -race` | Verified |
| Send Heartbeats with Ticker and select in Go | Go | Send periodic heartbeat timestamps from Go with one time.Ticker, one output channel, and select. | `snippets/go/send-heartbeats-with-ticker-and-select-in-go/` | `go test -race` | Verified |
| Cancel HTTP Client Requests with Context Timeout in Go | Go | Send an outbound HTTP GET request with a per-request context timeout so slow servers are canceled cleanly. | `snippets/go/cancel-http-client-requests-with-context-timeout-in-go/` | `go test -race` | Verified |
| Graceful Background Worker Shutdown with Go | Go | Run a buffered background worker in Go that drains accepted jobs, rejects new submissions during shutdown, and runs cleanup once. | `snippets/go/graceful-background-worker-shutdown-with-go/` | `go test -race` | Verified |
| Deduplicate In-Flight Work with singleflight in Go | Go | Deduplicate concurrent Go loads for the same key with singleflight.Group while still rerunning work after the first request finishes. | `snippets/go/deduplicate-in-flight-work-with-singleflight-in-go/` | `go test -race` | Verified |
| Handle Signals with signal.NotifyContext in Go | Go | Bridge OS signals into a Go context with signal.NotifyContext so shutdown code can wait on ctx.Done(). | `snippets/go/handle-signals-with-signal-notify-context-in-go/` | `go test -race` | Verified |
| Cancel External Commands with Context in Go | Go | Run an external command from Go with exec.CommandContext so timeouts and cancellations stop the process cleanly. | `snippets/go/cancel-external-commands-with-context-in-go/` | `go test -race` | Verified |
| Nginx Basic Auth Example | Configuration | A simple Nginx basic auth example. | `snippets/config/nginx-basic-auth-example/` | `manual verification` | Draft |
| Nginx Client Max Body Size Example | Configuration | Increase the allowed request body size in Nginx. | `snippets/config/nginx-client-max-body-size-example/` | `manual verification` | Draft |
| Nginx ACME Challenge Example | Configuration | Expose the ACME HTTP-01 challenge path in Nginx. | `snippets/config/nginx-acme-challenge-example/` | `manual verification` | Draft |
| Nginx HTTPS Redirect Example | Configuration | Redirect HTTP traffic to HTTPS in Nginx. | `snippets/config/nginx-https-redirect-example/` | `manual verification` | Draft |
| Nginx Load Balancer Baseline | Configuration | A simple Nginx load balancer example. | `snippets/config/nginx-load-balancer-baseline/` | `manual verification` | Draft |
| Nginx IP Rate Limit Example | Configuration | A simple Nginx IP rate limit example. | `snippets/config/nginx-ip-rate-limit-example/` | `manual verification` | Draft |
| Nginx Maintenance Mode Example | Configuration | Serve a maintenance page from Nginx. | `snippets/config/nginx-maintenance-mode-example/` | `manual verification` | Draft |
| Nginx SPA History Fallback Example | Configuration | Serve a single-page app with history fallback in Nginx. | `snippets/config/nginx-spa-history-fallback-example/` | `manual verification` | Draft |
| Nginx Static Cache Headers Example | Configuration | Set cache headers for static files in Nginx. | `snippets/config/nginx-static-cache-headers-example/` | `manual verification` | Draft |
| Nginx Static Files Example | Configuration | Serve static files directly from Nginx. | `snippets/config/nginx-static-files-example/` | `manual verification` | Draft |
| Nginx Real IP Behind Proxy Example | Configuration | Restore client IP addresses behind a trusted proxy in Nginx. | `snippets/config/nginx-real-ip-behind-proxy-example/` | `manual verification` | Draft |
| Nginx Reverse Proxy Headers Example | Configuration | Forward standard reverse proxy headers in Nginx. | `snippets/config/nginx-reverse-proxy-headers-example/` | `manual verification` | Draft |
| Nginx WebSocket Proxy Example | Configuration | Proxy WebSocket connections through Nginx. | `snippets/config/nginx-websocket-proxy-example/` | `manual verification` | Draft |
| Nginx IP Allow Deny Example | Configuration | Restrict a Nginx location by client IP ranges. | `snippets/config/nginx-ip-allow-deny-example/` | `manual verification` | Draft |
| Nginx Stub Status Example | Configuration | Expose a local-only Nginx stub_status endpoint. | `snippets/config/nginx-stub-status-example/` | `manual verification` | Draft |
| Nginx CORS Preflight Example | Configuration | Handle CORS preflight requests in Nginx for one API origin. | `snippets/config/nginx-cors-preflight-example/` | `manual verification` | Draft |
| Redis Baseline redis.conf | Configuration | A small redis.conf baseline with security, persistence, and logging settings for a new Redis deployment. | `snippets/config/redis-conf-baseline/` | `manual verification` | Draft |
| Redis systemd Service Unit | Configuration | A Linux systemd unit for running Redis in the foreground with automatic restarts and an example resource cap. | `snippets/config/redis-systemd-service/` | `manual verification` | Draft |
| Squid Multi-Egress Proxy Example | Configuration | Map one local Squid port per source IP so each port exits through a different public IP. | `snippets/config/squid-multi-egress-proxy-example/` | `manual verification` | Draft |
| Download YouTube Video with yt-dlp | Shell | A commented yt-dlp command for downloading one YouTube video with resumable fragments, explicit format sorting, and optional proxy support. | `snippets/shell/download-youtube-video-with-yt-dlp/` | `manual verification` | Draft |
| Paginate Rows with LIMIT and OFFSET in MySQL | SQL | Return one stable page of rows in MySQL with ORDER BY, LIMIT, and OFFSET. | `snippets/sql/paginate-rows-with-sql/` | `unittest` | Verified |
| Paginate Rows with Cursor Tokens in MySQL | SQL | Return the next stable page of rows in MySQL with keyset pagination on created_at and item_id. | `snippets/sql/paginate-rows-with-cursor-tokens-in-sql/` | `unittest` | Verified |
| Claim Queue Jobs with SKIP LOCKED in MySQL | SQL | Claim a small batch of queued jobs in MySQL with FOR UPDATE SKIP LOCKED inside one transaction. | `snippets/sql/claim-queue-jobs-with-skip-locked-in-sql/` | `unittest` | Verified |
| Count Rows per Group in MySQL | SQL | Return one row count per group value in MySQL with GROUP BY and COUNT(*). | `snippets/sql/count-rows-per-group-with-sql/` | `unittest` | Verified |
| Find Duplicate Values in MySQL | SQL | Return repeated non-null values in one MySQL column with GROUP BY and HAVING. | `snippets/sql/find-duplicate-values-with-sql/` | `unittest` | Verified |
| Update Rows with JOIN in MySQL | SQL | Fill missing row values in MySQL from a lookup table with UPDATE and JOIN. | `snippets/sql/update-rows-with-join-with-sql/` | `unittest` | Verified |
| Find Missing Related Rows with LEFT JOIN in MySQL | SQL | Return rows in MySQL that still have no related row with LEFT JOIN and IS NULL. | `snippets/sql/find-missing-related-rows-with-left-join-with-sql/` | `unittest` | Verified |
| Upsert Tenant Contacts via SQLite ON CONFLICT | SQL | Idempotently insert or update tenant-scoped contacts by the unique key tenant_id plus email in SQLite. | `snippets/sql/upsert-rows-on-conflict-with-sql/` | `unittest` | Verified |
| Compute Order and Refund Metrics in SQLite | SQL | Aggregate order counts, paid revenue, and refunds for one reporting window in SQLite. | `snippets/sql/build-dashboard-metrics-with-sql/` | `unittest` | Verified |
| Select Latest Row per Group in SQLite | SQL | Return one newest status row per device with ROW_NUMBER() and a stable tie-break on event_id. | `snippets/sql/select-latest-row-per-group-with-sql/` | `unittest` | Verified |
| Calculate Running Totals in SQLite | SQL | Return one running total per row in SQLite with SUM() OVER (...) and a stable tie-break on transaction_id. | `snippets/sql/calculate-running-totals-with-sql/` | `unittest` | Verified |
