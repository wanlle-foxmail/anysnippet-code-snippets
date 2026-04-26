Fix common file permission problems faster with `chmod` and `chown` commands for scripts, shared directories, recursive updates, and ownership changes.

## What This Snippet Covers

- Setting a typical config-file mode
- Making one script executable
- Applying safer recursive permissions to a directory tree
- Making every shell script in one folder executable
- Changing the owner of one file
- Changing the owner of one directory tree

## Before Using

- Replace `alice` and the example paths with real users and files on your machine.
- Double-check recursive permission and ownership commands before running them.
- Expect `chown` to require elevated privileges in many environments.

## Code

```sh
chmod 644 app.conf
# Give the owner read-write access and everyone else read-only access.

chmod 755 scripts/deploy.sh
# Make a script executable while keeping read access for everyone.

chmod -R u+rwX,g+rX,o-rwx ./shared
# Apply safer recursive permissions to a shared directory tree.

find ./scripts -type f -name '*.sh' -exec chmod 755 {} \;
# Make every shell script in one folder executable.

chown alice report.txt
# Change the owner of one file.

chown -R alice ./uploads
# Change the owner of one directory tree.
```

## Why These Commands Are Useful

- They cover the permission fixes people repeatedly make during setup, deployment, and cleanup.
- They keep common file and script modes easy to copy correctly.
- They show where recursive permission changes end and ownership changes begin.

## Limitations

- This snippet stays `Draft` because it depends on local users, files, and privileges.
- Recursive commands can affect more files than expected if the path is wrong.
- Ownership changes often need `sudo` or another elevated execution context.

## Manual Verification

1. Confirm `chmod` and `chown` are available in your shell.
2. Replace the example user names and paths with real local values.
3. Run the commands only on disposable or clearly understood files first.
4. Confirm the resulting permissions or owners match the intended state.

## Files

- `src/chmod_and_chown_commands_for_everyday_permission_fixes.sh`
- `snippet.json`