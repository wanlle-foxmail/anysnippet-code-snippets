rsync -av ./project/ /Volumes/Backup/project/
# Copy a directory while preserving timestamps, permissions, and symlinks.

rsync -av --delete --dry-run ./public/ user@example.com:/var/www/public/
# Preview a mirror sync before deleting remote files.

rsync -av --delete ./public/ user@example.com:/var/www/public/
# Mirror a local directory to a remote target.

rsync -av --progress --partial ./large-video.mov user@example.com:/srv/uploads/
# Show progress and keep partial data if the transfer is interrupted.

rsync -av --exclude '.git' --exclude 'node_modules' ./app/ user@example.com:/srv/app/
# Skip common local-only folders during deployment.

rsync -av -e 'ssh -p 2222' ./build/ user@example.com:/srv/build/
# Transfer over SSH on a custom port.