grep -n 'ERROR' app.log
# Show matching lines with line numbers.

grep -rn 'TODO' ./src
# Search recursively through a source directory.

grep -i 'warning' app.log
# Match text without caring about letter case.

grep -v '^#' .env.example
# Hide comment lines and keep active settings visible.

grep -C 2 'panic' server.log
# Show two lines of context before and after a match.

grep -E 'timeout|refused|reset' app.log
# Match any of several error words with one command.