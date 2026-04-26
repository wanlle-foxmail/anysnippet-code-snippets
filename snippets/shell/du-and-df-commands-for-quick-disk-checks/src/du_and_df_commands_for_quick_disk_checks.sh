du -sh .
# Show the total size of the current directory.

du -sh ./*
# Show one human-readable size per item in the current directory.

du -sh ./* | sort -h
# Sort directory sizes from small to large.

df -h
# Show free and used disk space for mounted filesystems.

df -h .
# Show the filesystem usage for the current path.

df -i .
# Show inode usage for the current path.