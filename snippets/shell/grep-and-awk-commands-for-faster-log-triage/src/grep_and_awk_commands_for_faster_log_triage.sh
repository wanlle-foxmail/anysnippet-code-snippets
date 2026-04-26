grep 'ERROR' app.log | awk '{print $1, $2, $NF}'
# Show the timestamp prefix and final field for error lines.

grep -i 'timeout' app.log | awk '{print $1, $2, $NF}'
# Keep a short timeout view with the timestamp prefix and final field.

grep ' 500 ' access.log | awk '{print $7}'
# Pull only the failing request path from HTTP 500 lines.

grep -E 'ERROR|WARN' app.log | awk '{count[$3]++} END {for (key in count) print key, count[key]}'
# Count how many log lines appear for each level field.

grep 'user_id=' app.log | awk -F'user_id=' '{split($2, parts, /[ \t]/); print parts[1]}'
# Extract only the user_id value from key-value style logs.

grep 'completed in' worker.log | awk '{print $(NF-1), $NF}'
# Pull the duration fields from job completion lines.