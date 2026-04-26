See what is holding a port or file with `lsof` commands for process lookup, listening sockets, open directories, and per-process inspection.

## What This Snippet Covers

- Finding which process is using one port
- Checking one TCP port specifically
- Listing listening TCP sockets without name lookups
- Checking which process holds one file or socket path
- Checking which processes have files open in one directory tree
- Inspecting open files for one process id

## Before Using

- Replace the example port numbers, process id, and paths with real local values.
- Expect directory scans such as `+D` to be slower in large trees.
- Some `lsof` results require elevated privileges to see every process.

## Code

```sh
lsof -i :3000
# Show which process is using port 3000.

lsof -i TCP:443
# Show which process is using TCP port 443.

lsof -iTCP -sTCP:LISTEN -n -P
# List listening TCP sockets without DNS or service-name lookups.

lsof /tmp/app.sock
# Show which process has one file or socket path open.

lsof +D ./logs
# Show which processes have files open inside one directory tree.

lsof -p 12345
# Show the open files and sockets for one process id.
```

## Why These Commands Are Useful

- They answer the first question behind many port-conflict or file-lock issues: who is holding it.
- They help you move from a failing service or locked file to the exact owning process.
- They keep port, path, and process-based views in one compact reference.

## Limitations

- This snippet stays `Draft` because results depend on local processes and open files.
- Some systems restrict visibility into other users' processes.
- Recursive directory scans with `+D` can be expensive in large folders.

## Manual Verification

1. Confirm `lsof -v` works in your environment.
2. Replace the example ports, process id, and paths with real local values.
3. Run the commands while a known process has a port or file open.
4. Confirm the reported process or path matches the expected owner.

## Files

- `src/lsof_commands_for_finding_ports_and_file_locks.sh`
- `snippet.json`