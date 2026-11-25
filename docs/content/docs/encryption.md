---
title: "End-to-End Encryption"
weight: 4
---

`pd` provides built-in End-to-End Encryption (E2EE) using the modern [age](https://github.com/FiloSottile/age) encryption tool or your existing SSH public keys.

All encryption and decryption happen strictly on your machine before uploading or after downloading. Pixeldrain only ever sees encrypted bytes.

---

## Encrypting Uploads

### Using an age Public Key or SSH Public Key

Encrypt a file for a specific recipient:

```shell
# Encrypt with an age recipient key
pd upload --recipient "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p" confidential.pdf

# Encrypt with an SSH public key
pd upload --recipient "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAI... user@example.com" confidential.pdf

# Encrypt to all supported SSH keys published by a GitHub user
pd upload --recipient "github:user-id" confidential.pdf
```

GitHub recipients are fetched from `https://github.com/<user-id>.keys` over HTTPS.
Supported Ed25519 and RSA keys are all added as recipients, so anyone holding a
matching private key can decrypt the upload. Unsupported key types are skipped
with a warning; the command fails if there are no supported keys. Verify the
GitHub account and its public keys through a trusted channel before uploading
sensitive data.

Uploaded files will automatically be saved with the `.age` extension appended.

### Multiple Recipients

Encrypt a file for recipients listed in a file (one age or SSH public key per line):

```shell
pd upload --recipient-file recipients.txt confidential.pdf
pd upload --recipient-file ~/.ssh/id_ed25519.pub confidential.pdf
```

Both options also accept hybrid age public keys (`age1pq1...`). Hybrid keys
cannot be combined with non-post-quantum recipients in the same upload.

---

## Decrypting Downloads

To download and decrypt the file automatically, pass your private identity file with `--identity`:

```shell
pd download --identity ~/.ssh/id_ed25519 https://pixeldrain.com/u/xxxxxx
```

SSH Ed25519 and RSA private keys are supported. If the key is protected by a
passphrase, `pd` prompts for it in a terminal.
