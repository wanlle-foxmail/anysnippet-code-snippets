Package handoff files more safely with `tar` and `OpenSSL` commands for archive creation, content checks, and checksum verification.

## What This Snippet Covers

- Creating an archive and hashing it immediately
- Excluding dependency folders before hashing an app bundle
- Saving a checksum file beside an archive
- Listing archive contents and printing a digest in one pass
- Extracting one file for spot checks before sharing
- Comparing digests for two archive copies

## Before Using

- Replace the example file names and directories with real local paths.
- Decide whether you want to print hashes to the terminal or save them to checksum files.
- Use disposable archives when testing extraction commands.

## Code

```sh
tar -czf release.tar.gz ./release && openssl dgst -sha256 release.tar.gz
# Create a release archive and print its SHA-256 digest.

tar -czf app.tar.gz --exclude='node_modules' ./app && openssl dgst -sha256 app.tar.gz
# Archive one app directory without dependencies and hash the result.

tar -czf logs.tar.gz ./logs && openssl dgst -sha256 logs.tar.gz > logs.tar.gz.sha256
# Create a logs archive and save its checksum beside the file.

tar -tzf release.tar.gz && openssl dgst -sha256 release.tar.gz
# List archive contents and print the archive hash in one check.

tar -xzf release.tar.gz release/README.md && openssl dgst -sha256 release/README.md
# Extract one file for spot checks and hash the extracted file.

openssl dgst -sha256 release.tar.gz release-copy.tar.gz
# Compare the printed digests of two archive copies.
```

## Why These Commands Are Useful

- They tie packaging and verification together so checksum work is not skipped.
- They support safer artifact handoffs when you need a fast integrity check.
- They cover both archive-level verification and spot checks of extracted files.

## Limitations

- This snippet stays `Draft` because it depends on local files, directories, and archives.
- Hash output comparison is still a manual check unless you add a verification script.
- Extraction examples write files into the working directory.

## Manual Verification

1. Confirm `tar` and `openssl` are available in your shell.
2. Replace the example directories, file names, and exclusions.
3. Create test folders with known contents before running the archive commands.
4. Confirm the printed or saved digests match the archives you created.

## Files

- `src/tar_and_openssl_commands_for_safer_archive_verification.sh`
- `snippet.json`