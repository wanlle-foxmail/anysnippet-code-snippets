Move files between local and remote machines with `scp` commands for uploads, downloads, recursive copies, custom ports, and key-based transfers.

## What This Snippet Covers

- Uploading one local file to a remote directory
- Downloading one remote file locally
- Uploading a whole directory tree recursively
- Using a custom SSH port for transfer
- Using a specific private key file
- Downloading a remote directory tree recursively

## Before Using

- Replace the example hosts, users, ports, and key paths with real values.
- Replace the example file paths with real local or remote files.
- Make sure the destination paths already exist or are writable.

## Code

```sh
scp report.txt user@example.com:/tmp/
# Upload one local file to a remote directory.

scp user@example.com:/var/log/app.log ./
# Download one remote file into the current directory.

scp -r ./build user@example.com:/srv/www/
# Upload a whole directory tree recursively.

scp -P 2222 report.txt user@example.com:/tmp/
# Transfer a file over a non-default SSH port.

scp -i ~/.ssh/deploy_key release.tar.gz user@example.com:/srv/releases/
# Transfer a file with a specific private key.

scp -r user@example.com:/srv/backups ./local-backups
# Download a whole remote directory tree recursively.
```

## Why These Commands Are Useful

- They cover the most common copy directions without jumping to a more complex sync tool.
- They keep upload, download, recursive copy, and key-based variations easy to copy.
- They match the transfer shapes developers often use during quick maintenance or deployment work.

## Limitations

- This snippet stays `Draft` because it depends on real hosts, paths, keys, and network access.
- Large directory copies can be slow or restart from scratch if interrupted.
- Transfers fail if the remote path does not exist or your account lacks permissions.

## Manual Verification

1. Confirm `scp -V` or `scp -h` works in your environment.
2. Replace the example hosts, users, ports, key paths, and file paths.
3. Run the commands only against hosts and paths you are allowed to access.
4. Confirm the uploaded or downloaded files appear in the intended destination.

## Files

- `src/scp_commands_for_simple_remote_file_transfer.sh`
- `snippet.json`