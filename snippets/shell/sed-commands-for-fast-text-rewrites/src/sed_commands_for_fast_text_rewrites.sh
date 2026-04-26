sed 's/http:/https:/' config.txt
# Replace the first matching text on each line.

sed 's/ERROR/FAILED/g' app.log > cleaned.log
# Replace all matching text and save the result to a new file.

sed '/^#/d' config.txt
# Drop comment lines from a config file.

sed 's/[[:space:]]*$//' input.txt
# Trim trailing whitespace from each line.

sed -n '1,20p' server.log
# Print only a chosen line range.

sed -n '/BEGIN CONFIG/,/END CONFIG/p' app.txt
# Print only the lines between two marker patterns.

sed '/^$/d' input.txt
# Remove blank lines from the output.