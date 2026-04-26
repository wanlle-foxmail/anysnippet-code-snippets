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