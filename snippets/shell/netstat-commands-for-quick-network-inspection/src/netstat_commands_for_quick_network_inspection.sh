netstat -an
# Show all sockets with numeric addresses.

netstat -an | grep LISTEN
# Show listening sockets only.

netstat -an | grep ESTABLISHED
# Show currently established connections.

netstat -rn
# Show the kernel routing table.

netstat -i
# Show per-interface counters.

netstat -s
# Show protocol-level network statistics.