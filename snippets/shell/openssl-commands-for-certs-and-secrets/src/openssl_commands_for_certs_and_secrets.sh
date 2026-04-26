openssl rand -base64 32
# Generate a random base64 secret.

openssl dgst -sha256 ./file.tar.gz
# Compute a SHA-256 hash for a file.

openssl x509 -in cert.pem -noout -subject -issuer -dates
# Inspect the subject, issuer, and validity dates of a certificate file.

openssl s_client -connect example.com:443 -servername example.com </dev/null 2>/dev/null | openssl x509 -noout -subject -issuer -dates
# Inspect the certificate presented by a live HTTPS server.

openssl req -new -newkey rsa:2048 -nodes -keyout server.key -out server.csr
# Generate a private key and an interactive certificate signing request.

openssl enc -aes-256-cbc -pbkdf2 -salt -in secrets.txt -out secrets.txt.enc
# Encrypt a file with AES-256-CBC and a passphrase.

openssl enc -d -aes-256-cbc -pbkdf2 -in secrets.txt.enc -out secrets.txt
# Decrypt a file created by the previous command.