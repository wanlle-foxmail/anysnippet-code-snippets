sort names.txt
# Sort lines alphabetically.

sort -u emails.txt
# Sort and drop duplicate lines in one step.

sort events.txt | uniq
# Collapse adjacent duplicates after sorting.

sort events.txt | uniq -d
# Show only lines that appear more than once.

sort events.txt | uniq -c | sort -nr
# Count repeated lines and rank the most common ones first.

sort -t',' -k2,2 report.csv
# Sort a CSV-like file by the second field.