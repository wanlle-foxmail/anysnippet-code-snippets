cut -d',' -f1 users.csv
# Extract the first comma-separated field.

cut -d',' -f1,3 users.csv
# Keep only selected comma-separated fields.

cut -c1-12 access.log
# Slice a fixed character range from each line.

cut -d'=' -f1 .env
# Keep only the key name from key-value lines.

paste first_names.txt last_names.txt
# Merge two files side by side with tab separators.

paste -d',' ids.txt scores.txt
# Merge two files side by side with a custom delimiter.