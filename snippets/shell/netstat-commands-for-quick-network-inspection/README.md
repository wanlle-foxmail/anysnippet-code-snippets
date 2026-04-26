Inspect connections and listening ports faster with `netstat` commands for socket state checks, routing tables, interface counters, and protocol statistics.

## What This Snippet Covers

- Showing all sockets with numeric addresses
- Filtering for listening sockets
- Filtering for established connections
- Viewing the routing table
- Viewing per-interface counters
- Viewing protocol-level statistics

## Before Using

- Expect output format differences between operating systems.
- Use `grep` filters only after confirming the socket-state words used by your platform.
- Run the commands on the machine whose network state you actually want to inspect.

## Code

```sh
netstat -an
# Show all sockets with numeric addresses.

netstat -an | grep LISTEN
# Show listening sockets only.

netstat -an | grep ESTABLISHED
# Show currently established connections.

netstat -rn
# Show the kernel routing table.

netstat -i
# Show per-interface counters.

netstat -s
# Show protocol-level network statistics.
```

## Why These Commands Are Useful

- They give a quick first-pass view of connection state, routes, and interface activity.
- They help distinguish between a local listener problem and a broader network-path problem.
- They keep the most common socket and routing checks together in one place.

## Limitations

- This snippet stays `Draft` because results depend on the local machine and current network state.
- `netstat` output varies more across platforms than many other shell tools.
- Some modern Linux environments prefer `ss`, but `netstat` remains widely recognized.

## Manual Verification

1. Confirm `netstat -h` or `netstat --help` works in your environment.
2. Run the commands on a machine with at least one active connection or listener.
3. Check the routing table and interface views against expected local network state.
4. Confirm the filtered output matches the socket states reported by your platform.

## Files

- `src/netstat_commands_for_quick_network_inspection.sh`
- `snippet.json`