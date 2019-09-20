package crypto

import "encoding/hex"

// wrapper for hex functions
func HexEncodeToString(b []byte) string {
	return hex.EncodeToString(b)
}

func HexDecodeStringToBytes(s string) ([]byte, error) {
	return hex.DecodeString(s)
}
