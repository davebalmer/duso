package runtime

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/duso-org/duso/pkg/script"
)

// builtinVerifyEd25519 verifies an Ed25519 signature
// Usage: verify_ed25519(data, signature, public_key_pem)
// Returns: true if signature is valid, false otherwise
func builtinVerifyEd25519(evaluator *Evaluator, args map[string]any) (any, error) {
	// Get data - support both positional (0) and named (data)
	var dataBytes []byte
	var dataArg any

	if d, ok := args["data"]; ok {
		dataArg = d
	} else if d, ok := args["0"]; ok {
		dataArg = d
	}

	if dataArg == nil {
		return nil, fmt.Errorf("verify_ed25519() requires a data argument")
	}

	// Handle binary data
	if val, ok := dataArg.(script.Value); ok && val.IsBinary() {
		binVal := val.AsBinary()
		if binVal != nil && binVal.Data != nil {
			dataBytes = *binVal.Data
		}
	} else if val, ok := dataArg.(*script.ValueRef); ok && val.Val.IsBinary() {
		binVal := val.Val.AsBinary()
		if binVal != nil && binVal.Data != nil {
			dataBytes = *binVal.Data
		}
	} else if str, ok := dataArg.(string); ok {
		dataBytes = []byte(str)
	} else {
		dataBytes = []byte(fmt.Sprintf("%v", dataArg))
	}

	if len(dataBytes) == 0 {
		return nil, fmt.Errorf("verify_ed25519() requires non-empty data")
	}

	// Get signature - support both positional (1) and named (signature)
	var signatureStr string
	if sig, ok := args["signature"]; ok {
		if sigStr, ok := sig.(string); ok {
			signatureStr = sigStr
		}
	} else if sig, ok := args["1"]; ok {
		if sigStr, ok := sig.(string); ok {
			signatureStr = sigStr
		}
	}

	if signatureStr == "" {
		return nil, fmt.Errorf("verify_ed25519() requires a signature string argument")
	}

	// Try to decode signature as base64url first (standard for Ed25519), then base64, then hex
	var signatureBytes []byte
	var err error

	signatureBytes, err = base64.RawURLEncoding.DecodeString(signatureStr)
	if err != nil {
		// Try standard base64
		signatureBytes, err = base64.StdEncoding.DecodeString(signatureStr)
		if err != nil {
			// Try hex
			signatureBytes, err = hex.DecodeString(signatureStr)
			if err != nil {
				return nil, fmt.Errorf("verify_ed25519() failed to decode signature: %v", err)
			}
		}
	}

	// Ed25519 signatures are 64 bytes
	if len(signatureBytes) != 64 {
		return false, nil // Invalid signature length
	}

	// Get public key PEM - support both positional (2) and named (public_key_pem)
	var keyPEM string
	if key, ok := args["public_key_pem"]; ok {
		if keyStr, ok := key.(string); ok {
			keyPEM = keyStr
		}
	} else if key, ok := args["2"]; ok {
		if keyStr, ok := key.(string); ok {
			keyPEM = keyStr
		}
	}

	if keyPEM == "" {
		return nil, fmt.Errorf("verify_ed25519() requires a public_key_pem string argument")
	}

	// Parse the public key - handle PEM-encoded format
	var pubKeyBytes []byte

	if strings.Contains(keyPEM, "-----BEGIN") {
		// Parse PEM format
		block, _ := pem.Decode([]byte(keyPEM))
		if block == nil {
			return nil, fmt.Errorf("verify_ed25519() failed to parse PEM block")
		}
		pubKeyBytes = block.Bytes
	} else {
		// Try to parse as raw base64-encoded DER
		decoded, err := base64.StdEncoding.DecodeString(keyPEM)
		if err != nil {
			// Try base64url as fallback
			decoded, err = base64.RawURLEncoding.DecodeString(keyPEM)
			if err != nil {
				return nil, fmt.Errorf("verify_ed25519() failed to decode public key: %v", err)
			}
		}
		pubKeyBytes = decoded
	}

	// Parse the public key
	var publicKey ed25519.PublicKey

	// Try to parse as a certificate first (some services send certificates)
	cert, certErr := x509.ParseCertificate(pubKeyBytes)
	if certErr == nil {
		// Successfully parsed as certificate
		var ok bool
		publicKey, ok = cert.PublicKey.(ed25519.PublicKey)
		if !ok {
			return nil, fmt.Errorf("verify_ed25519() certificate public key is not an Ed25519 key")
		}
	} else {
		// Try to parse as a bare public key (PKIX format)
		pubKeyInterface, err := x509.ParsePKIXPublicKey(pubKeyBytes)
		if err != nil {
			// Might be raw Ed25519 bytes (32 bytes)
			if len(pubKeyBytes) == 32 {
				publicKey = ed25519.PublicKey(pubKeyBytes)
			} else {
				return nil, fmt.Errorf("verify_ed25519() failed to parse public key: %v", err)
			}
		} else {
			var ok bool
			publicKey, ok = pubKeyInterface.(ed25519.PublicKey)
			if !ok {
				return nil, fmt.Errorf("verify_ed25519() public key is not an Ed25519 key")
			}
		}
	}

	// Verify the signature
	if ed25519.Verify(publicKey, dataBytes, signatureBytes) {
		return true, nil
	}

	return false, nil
}
