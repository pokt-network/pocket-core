package types

import (
	"encoding/hex"
	"github.com/pokt-network/pocket-core/crypto"
)

const (
	SUPPORTEDTOKENVERSIONS = "0.0.1" // todo
)

type AAT struct {
	Version              string `json:"version"`
	ApplicationSignature string `json:"signature"`
	ApplicationPublicKey string `json:"ApplicaitonAddress"`
	ClientPublicKey      string `json:"ClientPublicKey"`
}

func (a AAT) VersionIsIncluded() bool {
	if a.Version == "" {
		return false
	}
	return true
}

func (a AAT) VersionIsSupported() bool {
	if a.Version == SUPPORTEDTOKENVERSIONS {
		return true
	}
	return false
}

func (a AAT) Validate() error {
	if err := a.ValidateVersion(); err != nil {
		return err
	}
	if err := a.ValidateMessage(); err != nil {
		return err
	}
	if err := a.ValidateSignature(); err != nil {
		return err
	}
	return nil
}

func (a AAT) Hash() []byte {
	return crypto.SHA3FromString(a.ApplicationPublicKey + a.ClientPublicKey) // todo standardize
}

func (a AAT) ValidateVersion() error {
	// check for valid version
	if !a.VersionIsIncluded() {
		return MissingTokenVersionError
	}
	if !a.VersionIsSupported() {
		return UnsupportedTokenVersionError
	}
	return nil
}

func (a AAT) ValidateMessage() error {
	// check for valid application public key
	// todo pub key format verification
	if len(a.ApplicationPublicKey) == 0 {
		return MissingApplicationPublicKeyError
	}
	// todo pub key format verification
	if len(a.ClientPublicKey) == 0 {
		return MissingClientPublicKeyError
	}
	return nil
}

func (a AAT) ValidateSignature() error {
	// check for valid signature
	messageHash := a.Hash()
	publicKeyBytes, err := hex.DecodeString(a.ApplicationPublicKey)
	if err != nil {
		return err
	}
	sig, err := hex.DecodeString(a.ApplicationSignature)
	// todo crypto signature validation
	if !crypto.MockVerifySignature(publicKeyBytes, messageHash, sig) {
		return InvalidTokenSignatureErorr
	}
	return nil
}
