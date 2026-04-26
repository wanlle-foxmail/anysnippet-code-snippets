rsync -av --dry-run ./dist/ deploy@example.com:/srv/app/
# Preview which files would change on the remote server.

ssh deploy@example.com 'df -h /srv/app && ls -lah /srv/app'
# Check disk space and current files before syncing.

rsync -av --delete ./dist/ deploy@example.com:/srv/app/
# Push the local dist folder and remove stale remote files.

ssh deploy@example.com 'systemctl restart my-app && systemctl status my-app --no-pager'
# Restart the service and print its status after deployment.

rsync -av -e 'ssh -p 2222' ./dist/ deploy@example.com:/srv/app/
# Deploy over a non-default SSH port.

ssh -J bastion.example.com deploy@private.example.com 'uname -a && uptime'
# Run a quick post-deploy check through a jump host.