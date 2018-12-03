package simplcrypto

import (
	"encoding/base64"
	"strings"
)

// Base64URLEncode encodes a given byte array to a base64-url-encoded string
func Base64URLEncode(src []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(src), "=")
}

// Base64URLDecode decodes a given base64-URL-encoded string to a byte array
func Base64URLDecode(encoded string) ([]byte, error) {
	if m := len(encoded) % 4; m != 0 {
		encoded += strings.Repeat("=", 4-m)
	}
	return base64.URLEncoding.DecodeString(encoded)
}
