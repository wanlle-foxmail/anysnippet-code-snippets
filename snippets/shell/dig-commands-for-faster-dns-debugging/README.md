Check DNS records faster with `dig` commands for full answers, short A and AAAA results, MX and TXT lookups, resolver checks, and trace-style debugging.

## What This Snippet Covers

- Viewing the default DNS answer for a domain
- Checking IPv4 answers only
- Checking IPv6 answers only
- Looking up mail exchangers
- Querying a specific resolver for TXT data
- Tracing a lookup from root servers downward

## Before Using

- Replace `example.com` with a real domain.
- Replace `1.1.1.1` if you want to test a different recursive resolver.
- Expect `+trace` to take longer than a normal lookup.

## Code

```sh
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
```

## Why These Commands Are Useful

- They cover the record types and lookup styles people use most often during DNS debugging.
- They keep short-answer checks and deeper trace-style debugging in one compact reference.
- They make it easy to separate domain issues from resolver issues.

## Limitations

- This snippet stays `Draft` because it depends on live DNS records and network access.
- Resolver output changes over time as records, TTLs, or nameservers change.
- `+trace` can fail in restricted networks that block parts of the recursive path.

## Manual Verification

1. Confirm `dig -v` works in your environment.
2. Replace the example domain or resolver if needed.
3. Run the commands on real domains with known records.
4. Confirm the returned answers or trace path match the expected DNS setup.

## Files

- `src/dig_commands_for_faster_dns_debugging.sh`
- `snippet.json`