package session

import (
	"encoding/hex"
)

// extension/wrapper of types.Application specific to the session module
//type SessionAppPubKey types.Application // todo convert to just public key string
type SessionAppPubKey string // todo convert to just public key string

func (s SessionAppPubKey) Bytes() ([]byte, error) {
	return hex.DecodeString(string(s))
}

func (s SessionAppPubKey) Validate() error {
	// todo real key validation
	if b, err := s.Bytes(); err != nil || len(b) == 0 {
		return EmptyAppPubKeyError
	}
	return nil
}
