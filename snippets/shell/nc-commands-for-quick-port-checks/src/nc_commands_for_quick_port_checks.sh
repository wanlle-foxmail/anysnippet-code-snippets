nc -vz example.com 443
# Check whether a remote TCP port is reachable.

nc -vz 127.0.0.1 5432
# Check whether a local service port is reachable.

printf 'hello\n' | nc 127.0.0.1 9000
# Send one short payload to a listening TCP service.

nc -l 9000
# Start a simple TCP listener on a local port.

nc -u -vz 8.8.8.8 53
# Check whether a remote UDP port appears reachable.

nc -w 3 example.com 80 < /dev/null
# Try a TCP connection with a short timeout.