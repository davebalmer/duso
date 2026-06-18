package runtime

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"

	"github.com/duso-org/duso/pkg/script"
)

// builtinHMAC computes an HMAC using the specified algorithm
// Usage: hmac(algo, data, key)
// Returns: hex-encoded HMAC string
func builtinHMAC(evaluator *Evaluator, args map[string]any) (any, error) {
	// Get algorithm - support both positional (0) and named (algo)
	var algo string
	if a, ok := args["algo"]; ok {
		if aStr, ok := a.(string); ok {
			algo = aStr
		}
	} else if a, ok := args["0"]; ok {
		if aStr, ok := a.(string); ok {
			algo = aStr
		}
	}

	if algo == "" {
		return nil, fmt.Errorf("hmac() requires an algo argument (sha256, sha512, sha1, or md5)")
	}

	// Normalize algorithm
	switch algo {
	case "sha256", "sha512", "sha1", "md5":
		// Valid algorithms
	default:
		return nil, fmt.Errorf("hmac() unsupported algorithm: %s (use sha256, sha512, sha1, or md5)", algo)
	}

	// Get data - support both positional (1) and named (data)
	var dataBytes []byte
	var dataArg any

	if d, ok := args["data"]; ok {
		dataArg = d
	} else if d, ok := args["1"]; ok {
		dataArg = d
	}

	if dataArg == nil {
		return nil, fmt.Errorf("hmac() requires a data argument")
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

	// Get key - support both positional (2) and named (key)
	var keyBytes []byte
	var keyArg any

	if k, ok := args["key"]; ok {
		keyArg = k
	} else if k, ok := args["2"]; ok {
		keyArg = k
	}

	if keyArg == nil {
		return nil, fmt.Errorf("hmac() requires a key argument")
	}

	// Handle binary key
	if val, ok := keyArg.(script.Value); ok && val.IsBinary() {
		binVal := val.AsBinary()
		if binVal != nil && binVal.Data != nil {
			keyBytes = *binVal.Data
		}
	} else if val, ok := keyArg.(*script.ValueRef); ok && val.Val.IsBinary() {
		binVal := val.Val.AsBinary()
		if binVal != nil && binVal.Data != nil {
			keyBytes = *binVal.Data
		}
	} else if str, ok := keyArg.(string); ok {
		keyBytes = []byte(str)
	} else {
		keyBytes = []byte(fmt.Sprintf("%v", keyArg))
	}

	// Select hash function based on algorithm
	var h hash.Hash
	switch algo {
	case "sha256":
		h = hmac.New(sha256.New, keyBytes)
	case "sha512":
		h = hmac.New(sha512.New, keyBytes)
	case "sha1":
		h = hmac.New(sha1.New, keyBytes)
	case "md5":
		h = hmac.New(md5.New, keyBytes)
	}

	// Write data to hash
	h.Write(dataBytes)

	// Get hex-encoded result
	return hex.EncodeToString(h.Sum(nil)), nil
}
