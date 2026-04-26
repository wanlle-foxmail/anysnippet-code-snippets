Ship remote updates with less friction using `ssh` and `rsync` commands for previewing changes, syncing releases, restarting services, and checking hosts after deploys.

## What This Snippet Covers

- Previewing remote sync changes with a dry run
- Checking remote disk space and target files before deploy
- Syncing a release folder and deleting stale files
- Restarting a service and checking its status after deploy
- Deploying over a non-default SSH port
- Running a quick check through a jump host

## Before Using

- Replace the example hosts, usernames, ports, paths, and service names.
- Use `--dry-run` first before any `--delete` sync.
- Confirm SSH access and permissions before trying remote restarts.
- The `systemctl` example assumes a remote Linux host that uses `systemd`.
- The jump-host example requires a reachable bastion and a private target host.

## Code

```sh
rsync -av --dry-run ./dist/ deploy@example.com:/srv/app/
# Preview which files would change on the remote server.

ssh deploy@example.com 'df -h /srv/app && ls -lah /srv/app'
# Check disk space and current files before syncing.

rsync -av --delete ./dist/ deploy@example.com:/srv/app/
# Push the local dist folder and remove stale remote files.

ssh deploy@example.com 'systemctl restart my-app && systemctl status my-app --no-pager'
# Restart the service and print its status after deployment.

rsync -av -e 'ssh -p 2222' ./dist/ deploy@example.com:/srv/app/
# Deploy over a non-default SSH port.

ssh -J bastion.example.com deploy@private.example.com 'uname -a && uptime'
# Run a quick post-deploy check through a jump host.
```

## Why These Commands Are Useful

- They cover the practical deploy loop: preview, sync, restart, verify.
- They keep common remote options close to the command so repeat deploys are faster.
- They help reduce risky manual steps during small server updates.

## Limitations

- This snippet stays `Draft` because it depends on reachable hosts, SSH access, and remote permissions.
- This snippet was not executed in the repository environment because no reachable SSH target was available.
- `--delete` can remove files on the remote side if the local source is incomplete.
- Service restart commands depend on the remote init system and access level.

## Manual Verification

1. Confirm `ssh` and `rsync` are available in your shell.
2. Replace the example hosts, usernames, paths, ports, and service names.
3. Start with the dry-run command and confirm the proposed file changes.
4. Run the remote checks only against disposable or known-safe environments with reachable SSH access; the `systemctl` example also needs restart permission on a `systemd` host.

## Files

- `src/ssh_and_rsync_commands_for_faster_remote_deploys.sh`
- `snippet.json`