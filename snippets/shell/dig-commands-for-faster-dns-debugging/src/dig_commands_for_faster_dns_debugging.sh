dig example.com
# Show the full default DNS answer for a domain.

dig example.com A +short
# Show only IPv4 answers.

dig example.com AAAA +short
# Show only IPv6 answers.

dig example.com MX +short
# Show the mail exchangers for a domain.

dig @1.1.1.1 example.com TXT +short
# Query one specific resolver for TXT records.

dig +trace example.com
# Trace the DNS lookup from root servers downward.