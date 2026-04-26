awk '{print $1}' access.log
# Print the first whitespace-separated field from each line.

awk -F',' '{print $1 "," $3}' users.csv
# Keep selected comma-separated columns.

awk '$5 >= 500 {print}' access.log
# Keep only lines whose fifth field meets a threshold.

awk '{sum += $3} END {print sum}' sales.txt
# Add the values in the third field.

awk 'NR > 1 {print $2}' report.csv
# Skip a header row and print the second field.

awk '{count[$1]++} END {for (key in count) print key, count[key]}' events.txt
# Count how many times each first-field value appears.

awk 'BEGIN {FS=","} NR > 1 {total += $4} END {print total}' orders.csv
# Sum one comma-separated column after skipping the header.