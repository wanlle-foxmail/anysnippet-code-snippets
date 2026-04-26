du -sh ./* 2>/dev/null | sort -h
# Sort one level of directory sizes from small to large.

du -sh ./* 2>/dev/null | sort -hr | head -10
# Show the ten largest items in the current directory.

du -sh ./*/ 2>/dev/null | sort -hr
# Rank only the top-level directories by size.

du -sh ~/Downloads/* 2>/dev/null | sort -hr | head -10
# Find the biggest items in the Downloads folder.

du -sh ./node_modules/* 2>/dev/null | sort -hr | head -10
# Check which dependency folders are taking the most space.

du -sh ./cache/* 2>/dev/null | sort -hr
# Rank cache entries by size before cleanup.