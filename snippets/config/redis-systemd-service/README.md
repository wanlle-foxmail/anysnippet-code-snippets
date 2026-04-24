# Redis systemd Service Unit

A Linux `systemd` unit for running Redis in the foreground with automatic restarts, file descriptor limits, and an example 8 GB cgroup memory cap.

## Why This Unit Uses `Type=notify`

This snippet follows the official Redis `systemd` example and runs Redis in the foreground:

- `Type=notify`
- `--supervised systemd`
- `--daemonize no`

This is the preferred setup for a new `systemd` deployment because `systemd` should supervise the main Redis process directly instead of tracking a forked background process.

## Source Reference

This snippet is adapted from the official Redis 7.4 example `systemd` unit:

- https://raw.githubusercontent.com/redis/redis/7.4/utils/systemd-redis_server.service

It also depends on Redis `redis.conf` support for:

- `supervised systemd`
- `daemonize no`

## What This Unit Covers

- Starts Redis with `systemd` notification support
- Restarts on unexpected failure with `Restart=on-failure`
- Sets `LimitNOFILE=65535`
- Shows an example `MemoryMax=8G` cap
- Creates runtime, state, and log directories with `systemd`
- Runs Redis as the `redis` user and group

## Assumptions

- Redis is installed at `/usr/bin/redis-server`
- Your config file is at `/etc/redis/redis.conf`
- The host uses Linux with `systemd`
- The `redis` user and group already exist
- Your Redis config is compatible with foreground `systemd` supervision

## Important Redis Config Assumptions

Make sure your Redis config does not fight the service unit:

- Use `supervised systemd` or let the command line override it
- Keep `daemonize no`
- Do not rely on a `pidfile`-driven `Type=forking` model

If you also set `maxmemory` in `redis.conf`, do not set it equal to `MemoryMax=8G`.
Leave headroom for allocator overhead, replication buffers, and persistence work.

## How to Use It

1. Adjust the Redis binary path if your distribution installs it elsewhere.
2. Adjust the config path if your Redis config lives outside `/etc/redis/redis.conf`.
3. Adjust `MemoryMax` to match host RAM and your Redis `maxmemory` plan.
4. Save the unit as `/etc/systemd/system/redis-server.service`.
5. Run `systemctl daemon-reload`.
6. Run `systemctl enable --now redis-server`.

## Manual Verification

1. Run `systemd-analyze verify /etc/systemd/system/redis-server.service`.
2. Run `systemctl daemon-reload`.
3. Start the service with `systemctl start redis-server`.
4. Confirm `systemctl status redis-server` reports `active (running)`.
5. Check logs with `journalctl -u redis-server -b`.
6. Run `redis-cli` with the authentication or TLS flags required by your Redis config and confirm `PING` returns `PONG`.
7. Run `systemctl restart redis-server` and confirm the service comes back cleanly.

## Notes

- `Restart=on-failure` avoids automatic restarts after an intentional `systemctl stop`.
- `StateDirectory`, `LogsDirectory`, and `RuntimeDirectory` depend on modern `systemd` versions. On older distributions, create those paths manually instead.
- If your Redis build does not support `systemd` notifications, fall back to `Type=simple` with the same foreground `ExecStart` pattern instead of switching to `Type=forking`.