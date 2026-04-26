Inspect running processes faster with `ps` and `grep` commands for service lookup, process-name checks, and lightweight runtime debugging.

## What This Snippet Covers

- Checking for an nginx process
- Checking for a specific Node.js server command
- Searching for Redis-related processes case-insensitively
- Looking for queue worker processes with long ps output
- Checking for a Celery-style Python worker
- Checking for Postgres-related processes

## Before Using

- Replace the example process names with real commands or services from your machine.
- Use bracketed grep patterns like `[n]ginx` when you want to avoid matching the grep process itself.
- Broad patterns like `redis` or `postgres` can return several related process lines.
- Expect process names and arguments to vary across environments.

## Code

```sh
ps aux | grep '[n]ginx'
# Find running nginx processes without matching the grep command itself.

ps aux | grep '[n]ode .*server.js'
# Check whether one Node.js server process is running.

ps aux | grep -i '[r]edis'
# Search for Redis-related processes without caring about letter case.

ps -ef | grep '[q]ueue-worker'
# Find queue worker processes with the long ps format.

ps aux | grep '[p]ython .*celery'
# Check whether one Celery-style Python worker is active.

ps aux | grep '[p]ostgres'
# Check whether Postgres-related processes are running.
```

## Why These Commands Are Useful

- They cover the fast process checks people do before reaching for a deeper debugger.
- They keep the common bracketed-grep trick visible so results stay cleaner.
- They work well for quick service verification during setup or incident triage.

## Limitations

- This snippet stays `Draft` because it depends on the processes running on your local machine.
- Process arguments differ between environments, so example patterns often need adjustment.
- Broad search words can return multiple related processes or unrelated system helpers.
- `ps` output columns vary slightly between platforms.

## Manual Verification

1. Confirm `ps` and `grep` are available in your shell.
2. Replace the example process names or command fragments with real local values.
3. Run the commands while the target services are actually running.
4. Confirm the returned process lines match the intended service or command, and expect multiple matches for broader patterns.

## Files

- `src/ps_and_grep_commands_for_quick_process_checks.sh`
- `snippet.json`