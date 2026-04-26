Spot disk pressure faster with `du` and `df` commands for directory size checks, human-readable output, filesystem usage, and quick inode visibility.

## What This Snippet Covers

- Checking the total size of the current directory
- Listing item sizes one level down
- Sorting directory sizes for quick triage
- Checking overall filesystem capacity
- Checking the filesystem behind one path
- Checking inode usage for one path

## Before Using

- Run these commands in the directory or filesystem you actually want to inspect.
- Expect size totals to take longer in very large folders.
- Use inode checks when a disk looks full even though byte usage seems low.

## Code

```sh
du -sh .
# Show the total size of the current directory.

du -sh ./*
# Show one human-readable size per item in the current directory.

du -sh ./* | sort -h
# Sort directory sizes from small to large.

df -h
# Show free and used disk space for mounted filesystems.

df -h .
# Show the filesystem usage for the current path.

df -i .
# Show inode usage for the current path.
```

## Why These Commands Are Useful

- They answer the first questions people ask when a machine starts running out of space.
- They separate directory growth from filesystem capacity so the next action becomes clearer.
- They keep both byte usage and inode usage in one quick reference.

## Limitations

- This snippet stays `Draft` because results depend on the local filesystem state.
- Directory totals can take time to compute in large trees.
- Some filesystems report inode usage differently or do not make it equally meaningful.

## Manual Verification

1. Confirm `du` and `df` are available in your shell.
2. Run the commands in a directory with known files or nested folders.
3. Compare the reported sizes and filesystem usage against expected local state.
4. Confirm inode usage appears when your platform supports it.

## Files

- `src/du_and_df_commands_for_quick_disk_checks.sh`
- `snippet.json`