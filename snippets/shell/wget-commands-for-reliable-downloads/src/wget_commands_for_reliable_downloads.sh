wget -c https://example.com/archive.tar.gz
# Resume a partially downloaded file.

wget -O package.zip https://example.com/releases/package.zip
# Save the download with a custom file name.

wget --timestamping https://example.com/app-config.json
# Download only when the remote file is newer.

wget --retry-connrefused --waitretry=2 --tries=20 https://example.com/large-file.iso
# Retry cleanly when the server is temporarily unavailable.

wget --limit-rate=500k https://example.com/backup.tar.gz
# Throttle download speed to avoid saturating the network.

wget --mirror --convert-links --adjust-extension --page-requisites --no-parent https://docs.example.com/
# Mirror a small static site for offline browsing.