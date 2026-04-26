Handle common TLS and secret tasks faster with `openssl` commands for random keys, file hashes, certificate inspection, live HTTPS checks, CSRs, and file encryption.

## What This Snippet Covers

- Generating a random secret string
- Computing a SHA-256 file hash
- Inspecting certificate subject, issuer, and dates
- Checking the certificate presented by a live HTTPS server
- Creating a private key and CSR
- Encrypting a file with a passphrase
- Decrypting a previously encrypted file

## Before Using

- Replace `example.com` with a real HTTPS host for live certificate checks.
- Replace the example file names with local files that exist on your machine.
- Run file-writing commands from a writable directory.

## Code

```sh
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
```

## Why These Commands Are Useful

- They cover the certificate and secret tasks people often need in operations or debugging.
- They keep the command shapes easy to copy without forcing a long TLS tutorial first.
- They reduce the need to remember the exact `openssl` subcommand for each common task.

## Limitations

- This snippet stays `Draft` because it depends on placeholder hosts and local files.
- CSR generation is interactive unless you add more subject flags yourself.
- Encryption and decryption commands prompt for a passphrase during execution.

## Manual Verification

1. Confirm `openssl version` works.
2. Replace the placeholder host and file names.
3. Run each command from a writable directory.
4. Confirm the generated output, certificate details, or encrypted files match the intended task.

## Files

- `src/openssl_commands_for_certs_and_secrets.sh`
- `snippet.json`