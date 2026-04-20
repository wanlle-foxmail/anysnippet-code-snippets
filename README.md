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
| Process Large CSV Files in Chunks with pandas | Python | Read a large CSV with pandas in chunks and return one processed result per chunk. | `snippets/python/process-large-csv-files-in-chunks-with-python/` | `unittest` | Verified |
| Read Headerless CSV in Chunks with pandas | Python | Read a headerless CSV file in chunks and yield row dictionaries. | `snippets/python/read-headerless-csv-in-chunks-with-pandas/` | `unittest` | Verified |
| Limit I/O Concurrency in Python | Python | Run I/O-bound tasks with ThreadPoolExecutor and a fixed worker limit. | `snippets/python/limit-concurrent-tasks-in-python/` | `unittest` | Verified |
| Refactor If Elif Dispatch to a Handler Map in Python | Python | Replace a long if and elif dispatch chain with a handler mapping for email, sms, and webhook notifications. | `snippets/python/refactor-large-if-elif-dispatch-in-python/` | `unittest` | Verified |
| Refactor Nested If Else to a Config Map in Python | Python | Replace nested if else decision logic with a config mapping for support ticket routing. | `snippets/python/replace-nested-if-else-with-config-mapping-in-python/` | `unittest` | Verified |
| Download Large File with Python | Python | Stream-download a large file with requests, compute MD5 on the fly, and clean up on failure. | `snippets/python/download-large-file-with-python/` | `unittest` | Verified |
| FastAPI Health Check Endpoint with Python | Python | A /health endpoint that checks database and Redis, then returns aggregated service status. | `snippets/python/fastapi-health-check-endpoint-with-python/` | `unittest` | Verified |
| Echo Health Check Endpoint with Go | Go | A /health endpoint that checks database, Redis, and disk space with uptime reporting. | `snippets/go/echo-health-check-endpoint-with-go/` | `go test` | Verified |
| Redis Baseline Config for New Deployment | Configuration | A small redis.conf baseline with security, persistence, and logging settings for a new Redis deployment. | `snippets/config/recommended-redis-conf-settings-for-new-deployment/` | `manual verification` | Draft |
| Redis systemd Unit for New Deployment | Configuration | A Linux systemd unit for running Redis in the foreground with automatic restarts and an example resource cap. | `snippets/config/recommended-redis-systemd-service-for-new-deployment/` | `manual verification` | Draft |
| Download YouTube Video with yt-dlp | Shell | A commented yt-dlp command for downloading one YouTube video with resumable fragments, explicit format sorting, and optional proxy support. | `snippets/shell/download-youtube-video-with-yt-dlp/` | `manual verification` | Draft |
| Upsert Tenant Contacts via SQLite ON CONFLICT | SQL | Idempotently insert or update tenant-scoped contacts by the unique key tenant_id plus email in SQLite. | `snippets/sql/upsert-rows-on-conflict-with-sql/` | `unittest` | Verified |
| Compute Order and Refund Metrics in SQLite | SQL | Aggregate order counts, paid revenue, and refunds for one reporting window in SQLite. | `snippets/sql/build-dashboard-metrics-with-sql/` | `unittest` | Verified |
| Select Latest Row per Group with ROW_NUMBER() | SQL | Return one newest status row per device with ROW_NUMBER() and a stable tie-break on event_id. | `snippets/sql/select-latest-row-per-group-with-sql/` | `unittest` | Verified |
