Probe ports and test simple TCP or UDP connectivity with `nc` commands for reachability checks, listeners, short payloads, and timed connection attempts.

## What This Snippet Covers

- Checking whether a remote TCP port is reachable
- Checking whether a local TCP service is reachable
- Sending one short payload to a listening service
- Starting a simple local listener
- Trying a UDP port probe
- Using a short timeout for a connection attempt

## Before Using

- Replace the example hosts, ports, and payload targets with real values.
- Check which `nc` variant is installed, because flags differ across implementations.
- Run listener examples only on ports that are not already in use.

## Code

```sh
nc -vz example.com 443
# Check whether a remote TCP port is reachable.

nc -vz 127.0.0.1 5432
# Check whether a local service port is reachable.

printf 'hello\n' | nc 127.0.0.1 9000
# Send one short payload to a listening TCP service.

nc -l 9000
# Start a simple TCP listener on a local port.

nc -u -vz 8.8.8.8 53
# Check whether a remote UDP port appears reachable.

nc -w 3 example.com 80 < /dev/null
# Try a TCP connection with a short timeout.
```

## Why These Commands Are Useful

- They cover the lightweight connectivity checks people use before opening a heavier debugging tool.
- They help separate DNS success from actual socket reachability.
- They keep both client-side probes and a simple local listener in one place.

## Limitations

- This snippet stays `Draft` because it depends on real hosts, ports, and network access.
- `nc` flag behavior varies across OpenBSD, GNU, and other netcat builds.
- UDP reachability checks are less definitive than TCP checks.

## Manual Verification

1. Confirm `nc -h` works in your environment.
2. Replace the example hosts and ports with real local or remote targets.
3. Run the commands only against services and ports you are allowed to test.
4. Confirm the connection results or listener behavior match the intended network setup.

## Files

- `src/nc_commands_for_quick_port_checks.sh`
- `snippet.json`