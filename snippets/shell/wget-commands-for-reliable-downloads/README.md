Keep downloads moving with practical `wget` commands for resume support, custom file names, timestamp checks, retries, bandwidth limits, and small site mirrors.

## What This Snippet Covers

- Resuming interrupted downloads
- Saving output with a predictable file name
- Skipping unchanged remote files
- Retrying on temporary connection failures
- Limiting bandwidth for shared networks
- Mirroring a small static site for offline use

## Before Using

- Replace `https://example.com/...` URLs with real downloadable targets.
- Run mirror or file-saving commands from a writable directory.
- Be careful with `--mirror` on large sites, because it can create many files.

## Code

```sh
wget -c https://example.com/archive.tar.gz
# Resume a partially downloaded file.

wget -O package.zip https://example.com/releases/package.zip
# Save the download with a custom file name.

wget --timestamping https://example.com/app-config.json
# Download only when the remote file is newer.

wget --retry-connrefused --waitretry=2 --tries=20 https://example.com/large-file.iso
# Retry cleanly when the server is temporarily unavailable.

wget --limit-rate=500k https://example.com/backup.tar.gz
# Throttle download speed to avoid saturating the network.

wget --mirror --convert-links --adjust-extension --page-requisites --no-parent https://docs.example.com/
# Mirror a small static site for offline browsing.
```

## Why These Commands Are Useful

- They cover the most common download workflows without requiring a separate script.
- They show the small option changes that make `wget` safer on unstable networks.
- They give you a quick starting point for both single-file and small-site downloads.

## Limitations

- This snippet stays `Draft` because it depends on placeholder URLs and live network access.
- Mirror commands can pull more content than expected if the target site is not tightly scoped.
- Retry behavior depends on the remote server and the installed `wget` version.

## Manual Verification

1. Confirm `wget --version` works.
2. Replace the placeholder URLs with real download targets.
3. Run each command from a writable directory.
4. Confirm resumed files, saved names, retries, limits, or mirrored output behave as expected.

## Files

- `src/wget_commands_for_reliable_downloads.sh`
- `snippet.json`