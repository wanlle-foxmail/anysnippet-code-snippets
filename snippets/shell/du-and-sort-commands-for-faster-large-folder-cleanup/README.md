Find space-heavy folders faster with `du` and `sort` commands for directory ranking, top-N cleanup candidates, and quick storage triage.

## What This Snippet Covers

- Sorting one level of file and folder sizes
- Showing the biggest items in one directory
- Ranking only top-level directories by size
- Checking the biggest items in Downloads
- Inspecting large dependency folders
- Ranking cache entries before deleting them

## Before Using

- Replace the example paths with real local directories.
- Expect the commands to take longer in very large folders.
- Redirect stderr as shown when you want to ignore permission or missing-path noise.

## Code

```sh
du -sh ./* 2>/dev/null | sort -h
# Sort one level of directory sizes from small to large.

du -sh ./* 2>/dev/null | sort -hr | head -10
# Show the ten largest items in the current directory.

du -sh ./*/ 2>/dev/null | sort -hr
# Rank only the top-level directories by size.

du -sh ~/Downloads/* 2>/dev/null | sort -hr | head -10
# Find the biggest items in the Downloads folder.

du -sh ./node_modules/* 2>/dev/null | sort -hr | head -10
# Check which dependency folders are taking the most space.

du -sh ./cache/* 2>/dev/null | sort -hr
# Rank cache entries by size before cleanup.
```

## Why These Commands Are Useful

- They help you identify cleanup targets before you start deleting anything.
- They turn raw size output into ranked lists you can act on quickly.
- They work well for local triage in developer folders such as Downloads, caches, and dependencies.

## Limitations

- This snippet stays `Draft` because it depends on local folders and disk usage state.
- Some example paths may not exist on your machine.
- Human-readable sorting depends on `sort -h` support in your environment.

## Manual Verification

1. Confirm `du` and `sort` are available in your shell.
2. Replace the example paths with real local directories.
3. Run the commands in folders with known large files or subdirectories.
4. Confirm the ranked output matches the actual largest items.

## Files

- `src/du_and_sort_commands_for_faster_large_folder_cleanup.sh`
- `snippet.json`