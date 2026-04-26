tar -czf backup.tar.gz ./project
# Create a gzipped archive from one directory.

tar -xzf backup.tar.gz
# Extract a gzipped archive into the current directory.

tar -tzf backup.tar.gz
# List the contents of a gzipped archive without extracting it.

tar -czf release.tar.gz README.md src/ dist/
# Archive several files and folders in one bundle.

tar -xzf backup.tar.gz project/config.env
# Extract only one path from an archive.

tar -czf app.tar.gz --exclude='node_modules' ./app
# Create an archive while skipping a noisy folder.