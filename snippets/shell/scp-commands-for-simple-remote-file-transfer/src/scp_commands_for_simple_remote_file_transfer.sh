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