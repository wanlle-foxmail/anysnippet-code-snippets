sed 's/[[:space:]]*$//' messy.txt | awk 'NF'
# Trim trailing whitespace and drop blank lines.

sed '/^#/d' .env | awk -F'=' '{print $1}'
# Keep only active environment variable names.

sed 's/,/ /g' users.csv | awk '{print $1, $3}'
# Replace commas with spaces and print selected fields.

sed -n '1,100p' access.log | awk '$9 >= 500 {print $1, $7, $9}'
# Inspect the first 100 log lines and keep only HTTP error rows.

sed '/^$/d' input.txt | awk '{print NR ":" $0}'
# Remove blank lines and number the remaining lines.

sed 's/"//g' report.csv | awk -F',' '{print $1 "," $4}'
# Strip double quotes before selecting two CSV-like fields.