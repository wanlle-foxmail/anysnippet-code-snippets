lsof -i :3000
# Show which process is using port 3000.

lsof -i TCP:443
# Show which process is using TCP port 443.

lsof -iTCP -sTCP:LISTEN -n -P
# List listening TCP sockets without DNS or service-name lookups.

lsof /tmp/app.sock
# Show which process has one file or socket path open.

lsof +D ./logs
# Show which processes have files open inside one directory tree.

lsof -p 12345
# Show the open files and sockets for one process id.