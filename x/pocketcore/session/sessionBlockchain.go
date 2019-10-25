package session

import "encoding/hex"

// extension/wrapper of legacy.Blockchain for session
// TODO non native chains need to be defined by config, for now will be hash
type SessionBlockchain string

func (sbc SessionBlockchain) Validate() error {
	// todo more validation
	if len(sbc) == 0 {
		return EmptyNonNativeChainError
	}
	return nil
}

func (sbc SessionBlockchain) Bytes() ([]byte, error) {
	return hex.DecodeString(string(sbc))
}
