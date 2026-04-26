Work on remote machines more efficiently with `ssh` commands for aliases, one-off checks, key files, custom ports, local forwarding, and jump hosts.

## What This Snippet Covers

- Connecting through an SSH config alias
- Running a one-off command on a remote machine
- Using a specific private key file
- Connecting over a custom port
- Forwarding a remote service to a local port
- Reaching a private machine through a jump host

## Before Using

- Replace the example hosts, users, ports, and key paths with real values.
- Make sure your SSH config or bastion host actually exists before using alias or jump examples.
- Keep local port forwards on unused local ports.

## Code

```sh
ssh app-server
# Connect to a host alias from your SSH config.

ssh user@example.com 'uptime && df -h .'
# Run a quick one-off command on a remote machine.

ssh -i ~/.ssh/deploy_key user@example.com
# Connect with a specific private key file.

ssh -p 2222 user@example.com
# Connect to a host over a non-default SSH port.

ssh -L 8080:127.0.0.1:5432 user@example.com
# Forward a remote service to a local port.

ssh -J bastion.example.com user@private.example.com
# Reach a private host through a jump host.
```

## Why These Commands Are Useful

- They cover the SSH flows people repeat during debugging, access, and local tunneling.
- They keep the highest-value connection variants close at hand.
- They show both direct access and multi-hop access without burying the reader in config details.

## Limitations

- This snippet stays `Draft` because it depends on real hosts, keys, ports, and network access.
- Port-forward commands keep running until the session ends.
- Jump-host access depends on SSH support for `-J` and working credentials on both hops.

## Manual Verification

1. Confirm `ssh -V` works in your environment.
2. Replace the example hosts, users, ports, and key paths.
3. Run the commands only against hosts you are allowed to access.
4. Confirm the connection, remote output, or forwarded port behaves as intended.

## Files

- `src/ssh_commands_for_faster_remote_access.sh`
- `snippet.json`