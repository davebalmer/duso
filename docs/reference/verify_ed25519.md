# verify_ed25519()

Verify an Ed25519 signature.

`verify_ed25519(data, signature, public_key_pem)`

## Parameters

- `data` (string | binary) - The data that was signed
- `signature` (string) - Base64url-encoded or hex-encoded Ed25519 signature
- `public_key_pem` (string) - PEM-encoded Ed25519 public key

## Returns

Boolean: `true` if signature is valid, `false` if invalid (never throws on verification failure)

## Examples

Verify a signature with a public key:

```duso
public_key = load("/path/to/ed25519_public_key.pem")
data = "message to verify"
signature = load("message.sig")
is_valid = verify_ed25519(data, signature, public_key)
print("Signature valid: " + tostring(is_valid))
```

Discord webhook verification:

```duso
signature = req.headers["x-signature-ed25519"]
timestamp = req.headers["x-signature-timestamp"]
body = req.body
discord_public_key = "your-discord-public-key"

message = timestamp + body
if verify_ed25519(message, signature, discord_public_key) then
  print("Webhook signature verified")
else
  print("Invalid webhook signature - rejecting")
end
```

Verify signed binary data (e.g., file):

```duso
public_key = load("/path/to/ed25519_public_key.pem")
file_data = load_binary("document.pdf")
file_sig = load("document.pdf.sig")
if verify_ed25519(file_data, file_sig, public_key) then
  print("File signature verified")
else
  print("File signature invalid - data may be tampered")
end
```

## Key Format

Accepts PEM-encoded Ed25519 public keys:

```
-----BEGIN PUBLIC KEY-----
[base64 encoded Ed25519 key]
-----END PUBLIC KEY-----
```

Generate an Ed25519 key pair with OpenSSL:

```bash
openssl genpkey -algorithm ed25519 -out private_key.pem
openssl pkey -in private_key.pem -pubout -out public_key.pem
```

## Security Notes

- Uses Ed25519 (Edwards-curve Digital Signature Algorithm)
- Returns `false` instead of throwing on verification failure
- Safe to use in conditionals without try/catch for invalid signatures
- Only throws on PEM parsing errors or missing/invalid parameters
- Public keys can be safely shared - only private keys need protection
- Ed25519 signatures are deterministic and shorter than RSA/ECDSA

## Common Use Cases

- **Discord webhook verification**: Verify incoming Discord webhook signatures
- **Webhook verification**: Verify signed webhooks from services using Ed25519
- **Digital signatures**: Verify signed documents
- **API authentication**: Verify Ed25519-signed API requests
- **Code signing verification**: Verify released code integrity

## Signature Failure vs Errors

**Returns false** (signature doesn't match):
- Data was modified
- Wrong signature provided
- Wrong public key used
- Original signature was tampered with

**Throws error** (validation error):
- Public key PEM is invalid or unparseable
- Key is not an Ed25519 key
- Signature is not valid base64/hex
- Parameters missing

## See Also

- [sign_ed25519() - Create Ed25519 signatures](/docs/reference/sign_ed25519.md)
- [hmac() - Compute HMAC for message authentication](/docs/reference/hmac.md)
- [verify_ec() - EC signature verification](/docs/reference/verify_ec.md)
- [verify_rsa() - RSA signature verification](/docs/reference/verify_rsa.md)
- [hash() - Compute cryptographic hashes](/docs/reference/hash.md)
