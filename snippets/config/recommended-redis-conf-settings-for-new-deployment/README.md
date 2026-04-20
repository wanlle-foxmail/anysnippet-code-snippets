# Redis Baseline Config for New Deployment

A small `redis.conf` baseline for a new Redis deployment, using only directive names from the official Redis 7.4 `redis.conf` template.

## What This Snippet Covers

- Restricts listening addresses with `bind`
- Keeps `protected-mode yes`
- Uses `aclfile` for Redis 6+ authentication
- Changes the default `port`
- Enables both RDB snapshots and AOF persistence
- Sets explicit data and log paths
- Keeps `rdbchecksum` and `stop-writes-on-bgsave-error`
- Enables `slowlog`, `tcp-keepalive`, and `disable-thp`

## Source Version

This snippet is based on the official Redis 7.4 `redis.conf` template:

- https://raw.githubusercontent.com/redis/redis/7.4/redis.conf

## How to Use It

- Treat `redis.conf` in this snippet as a baseline fragment, not a full replacement for every deployment.
- Replace the port, ACL file path, data path, and log path before applying it.
- If you keep `aclfile`, create that file before startup. Also ensure the data directory and the parent directory of the log file already exist and have the correct permissions.
- If clients connect from other hosts, replace the loopback-only `bind` value with your private server addresses.

## Authentication Note

For a new deployment, prefer `aclfile` over `requirepass`.

Redis 6+ treats `requirepass` as a compatibility layer on top of ACLs. If you only need a single shared password, you can use `requirepass` instead of `aclfile`, but do not configure both at the same time.

## Additional Directives Worth Reviewing

These directives exist in the official `redis.conf`, but they are not included in this minimal baseline because they depend heavily on topology or workload:

- `tls-port`, `tls-cert-file`, `tls-key-file`, `tls-ca-cert-file`
- `maxmemory`, `maxmemory-policy`
- `latency-monitor-threshold`
- `replicaof`, `masterauth`, and other `repl-*` directives
- `cluster-enabled` and other `cluster-*` directives
- `unixsocket`, `unixsocketperm`

## Manual Verification

1. Compare each directive in `redis.conf` against the official Redis 7.4 template.
2. Replace the placeholder values and file paths for your environment.
3. Create the ACL file, data directory, and log directory required by your chosen paths.
4. Merge the fragment into a staging Redis configuration.
5. Start Redis in staging with the merged configuration and confirm the server accepts the settings.

## Notes

- Changing the port reduces background scan noise, but it is not a substitute for `bind`, `protected-mode`, and proper authentication.
- This snippet is Linux-oriented because it uses `supervised auto`, `disable-thp yes`, `/var/lib/redis`, and `/var/log/redis` paths.