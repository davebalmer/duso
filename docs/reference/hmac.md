# hmac()

Compute a Hash-based Message Authentication Code for message authentication and integrity verification.

`hmac(algo, data, key)`

## Parameters

- `algo` (string) - The hash algorithm to use: `"sha256"`, `"sha512"`, `"sha1"`, or `"md5"`
- `data` (string | binary) - The string or binary data to authenticate
- `key` (string | binary) - The secret key used for HMAC computation

## Returns

Hex-encoded HMAC string

## Examples

Compute HMAC with SHA256 (recommended):

```duso
data = "hello world"
secret = "my-secret-key"
mac = hmac("sha256", data, secret)
print(mac)  // Hex-encoded HMAC-SHA256
```

Slack webhook signature verification:

```duso
signature_header = req.headers["x-slack-signature"]
timestamp = req.headers["x-slack-request-timestamp"]
body = req.body
signing_secret = "your-slack-signing-secret"

message = "v0:" + timestamp + ":" + body
computed = hmac("sha256", message, signing_secret)
expected = signature_header.replace("v0=", "")

if computed == expected then
  print("Webhook signature verified")
else
  print("Invalid webhook signature - rejecting")
end
```

API request signing:

```duso
// Sign an outbound API request
method = "POST"
path = "/api/v1/resource"
body = format_json({name = "Alice", email = "alice@example.com"})
timestamp = now()
api_secret = "your-api-secret"

message = method + path + tostring(timestamp) + body
signature = hmac("sha256", message, api_secret)

// Include signature in request header
response = fetch("https://api.example.com" + path, {
  method = "POST",
  body = body,
  headers = {
    "X-Timestamp" = tostring(timestamp),
    "X-Signature" = signature
  }
})
```

HMAC with SHA512:

```duso
long_mac = hmac("sha512", "password123", "secret")
print(long_mac)  // Full 128-character hex string
```

Using named arguments:

```duso
result = hmac(algo = "sha256", data = "some data", key = "secret")
print(result)
```

HMAC with binary data:

```duso
image = load_binary("photo.png")
key = "binary-secret-key"
image_mac = hmac("sha256", image, key)
print("Image HMAC: " + image_mac)
```

## Algorithm Notes

- **sha256**: 64-character hex string (256 bits) - Good for general use, recommended
- **sha512**: 128-character hex string (512 bits) - Larger HMAC, slower but more secure
- **sha1**: 40-character hex string (160 bits) - Legacy, not recommended for security-critical use
- **md5**: 32-character hex string (128 bits) - Legacy, cryptographically broken, don't use for security

## Key Format

The key can be:
- A string (e.g., API secret, signing secret)
- Binary data (raw bytes)
- Any length (typically 32+ bytes for sha256)

For best security:
- Use SHA256 or SHA512
- Use a key with sufficient entropy (at least 256 bits / 32 bytes)
- Store keys securely (environment variables, key management services)
- Rotate keys periodically

## Performance

- All algorithms are fast and suitable for high-volume authentication
- HMAC computation is deterministic (same inputs always produce same output)
- Same input produces identical HMAC across all calls

## Common Use Cases

- **Webhook verification**: Verify incoming webhook signatures from Slack, GitHub, etc.
- **API request signing**: Sign outbound API requests with shared secret
- **Message authentication**: Verify message integrity and authenticity
- **Token generation**: Create authenticated tokens with shared secret
- **Cookie signing**: Sign session cookies to prevent tampering

## Differences from hash()

- `hash()` computes a one-way hash (no key)
- `hmac()` uses a secret key to produce an authenticated hash
- HMAC is suitable for message authentication, hash() is for integrity checking
- Use `hmac()` when you want to verify authenticity (both parties know the secret)
- Use `hash()` when you only care about integrity (one-way verification)

## Security Notes

- HMAC provides both integrity and authenticity (with shared secret)
- Always use HMAC in constant-time comparison (to prevent timing attacks)
- Keep the secret key confidential - whoever has it can forge HMACs
- Use strong, random keys (avoid passwords as keys)
- SHA256 is recommended for new applications

## See Also

- [hash() - Compute cryptographic hashes](/docs/reference/hash.md)
- [verify_ed25519() - Ed25519 signature verification](/docs/reference/verify_ed25519.md)
- [verify_ec() - EC signature verification](/docs/reference/verify_ec.md)
- [verify_rsa() - RSA signature verification](/docs/reference/verify_rsa.md)
- [encode_base64() - Base64 encoding](/docs/reference/encode_base64.md)
