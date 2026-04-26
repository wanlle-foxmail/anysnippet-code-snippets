chmod 644 app.conf
# Give the owner read-write access and everyone else read-only access.

chmod 755 scripts/deploy.sh
# Make a script executable while keeping read access for everyone.

chmod -R u+rwX,g+rX,o-rwx ./shared
# Apply safer recursive permissions to a shared directory tree.

find ./scripts -type f -name '*.sh' -exec chmod 755 {} \;
# Make every shell script in one folder executable.

chown alice report.txt
# Change the owner of one file.

chown -R alice ./uploads
# Change the owner of one directory tree.