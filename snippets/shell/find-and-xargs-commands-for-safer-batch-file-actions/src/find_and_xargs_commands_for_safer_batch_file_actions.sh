find . -type f -name '*.tmp' -print0 | xargs -0 -n1 echo rm -f
# Preview one delete command per matched file before you run a real delete.

find . -type f -name '*.tmp' -print0 | xargs -0 rm -f
# Delete matched temporary files safely with null-delimited paths.

find . -type f -name '*.jpg' -print0 | xargs -0 -I{} mv '{}' ./images-backup/
# Move matched image files into a backup directory.

find . -type f -name '*.sh' -print0 | xargs -0 chmod 755
# Make each matched shell script executable.

find . -type f -size +100M -print0 | xargs -0 ls -lh
# Inspect large files with human-readable sizes.

find . -type f -name '*.csv' -print0 | xargs -0 gzip
# Compress every matched CSV file in one batch operation.