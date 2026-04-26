Make backups and deployments quicker and safer with `rsync` commands for previewing changes, skipping noisy folders, resuming large transfers, and syncing remote targets.

## What This Snippet Covers

- Local directory backups with metadata preserved
- Dry-run previews before destructive mirror syncs
- Remote directory mirroring over SSH
- Resumable large-file transfers with progress output
- Excluding noisy local-only folders during deployment
- Connecting to remote targets over a custom SSH port

## Before Using

- Replace `user@example.com` and remote paths with real SSH targets.
- Replace local paths with directories or files that exist on your machine.
- Use `--dry-run` first when testing `--delete` on real data.

## Code

```sh
rsync -av ./project/ /Volumes/Backup/project/
# Copy a directory while preserving timestamps, permissions, and symlinks.

rsync -av --delete --dry-run ./public/ user@example.com:/var/www/public/
# Preview a mirror sync before deleting remote files.

rsync -av --delete ./public/ user@example.com:/var/www/public/
# Mirror a local directory to a remote target.

rsync -av --progress --partial ./large-video.mov user@example.com:/srv/uploads/
# Show progress and keep partial data if the transfer is interrupted.

rsync -av --exclude '.git' --exclude 'node_modules' ./app/ user@example.com:/srv/app/
# Skip common local-only folders during deployment.

rsync -av -e 'ssh -p 2222' ./build/ user@example.com:/srv/build/
# Transfer over SSH on a custom port.
```

## Why These Commands Are Useful

- They cover the `rsync` tasks people repeat in backup and deployment workflows.
- They keep the safer preview-first pattern visible when deletions are involved.
- They show practical options that reduce rework on slow or interrupted transfers.

## Limitations

- This snippet stays `Draft` because it depends on placeholder machines, paths, and SSH access.
- `--delete` can remove destination files if the source path is wrong.
- Remote transfers depend on network connectivity, SSH configuration, and remote permissions.

## Manual Verification

1. Confirm `rsync --version` works.
2. Replace the placeholder hosts and paths.
3. Run the `--dry-run` example before any real delete sync.
4. Confirm the copied or mirrored output matches the intended result.

## Files

- `src/rsync_commands_for_fast_sync_and_backup.sh`
- `snippet.json`