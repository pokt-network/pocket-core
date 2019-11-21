package types

import "encoding/hex"

type Blockchain string

func (b Blockchain) Validate() error {
	// todo more validation
	if len(b) == 0 {
		return EmptyNonNativeChainError
	}
	return nil
}

func (b Blockchain) Bytes() ([]byte, error) {
	return hex.DecodeString(string(b))
}
