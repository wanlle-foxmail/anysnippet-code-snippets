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