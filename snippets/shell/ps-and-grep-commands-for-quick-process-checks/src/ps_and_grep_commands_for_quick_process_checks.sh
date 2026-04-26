ps aux | grep '[n]ginx'
# Find running nginx processes without matching the grep command itself.

ps aux | grep '[n]ode .*server.js'
# Check whether one Node.js server process is running.

ps aux | grep -i '[r]edis'
# Search for Redis-related processes without caring about letter case.

ps -ef | grep '[q]ueue-worker'
# Find queue worker processes with the long ps format.

ps aux | grep '[p]ython .*celery'
# Check whether one Celery-style Python worker is active.

ps aux | grep '[p]ostgres'
# Check whether Postgres-related processes are running.