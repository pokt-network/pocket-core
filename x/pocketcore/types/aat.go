package types

import (
	"encoding/hex"
	"encoding/json"
)

const (
	SUPPORTEDTOKENVERSION = "0.0.1" // todo
)

type AAT struct {
	Version              string `json:"version"`
	ApplicationPublicKey string `json:"ApplicaitonAddress"`
	ClientPublicKey      string `json:"ClientPublicKey"`
	ApplicationSignature string `json:"signature"`
}

func (a AAT) VersionIsIncluded() bool {
	if a.Version == "" {
		return false
	}
	return true
}

func (a AAT) VersionIsSupported() bool {
	if a.Version == SUPPORTEDTOKENVERSION {
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
	type aat struct {
		ApplicationSignature string
		AppPubKey            string
		ClientPubKey         string
		Version              string
	}
	r, err := json.Marshal(aat{
		ApplicationSignature: "",
		AppPubKey:            a.ApplicationPublicKey,
		ClientPubKey:         a.ClientPublicKey,
		Version:              a.Version,
	})
	if err != nil {
		panic(err)
	}
	return SHA3FromBytes(r)
}

func (a AAT) HashString() string {
	return hex.EncodeToString(a.Hash())
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
	if len(a.ApplicationPublicKey) == 0 {
		return MissingApplicationPublicKeyError
	}
	if err := PubKeyVerification(a.ApplicationPublicKey); err != nil {
		return err
	}
	if len(a.ClientPublicKey) == 0 {
		return MissingClientPublicKeyError
	}
	if err := PubKeyVerification(a.ClientPublicKey); err != nil {
		return err
	}
	return nil
}

func (a AAT) ValidateSignature() error {
	// check for valid signature
	messageHash := a.HashString()
	if err := SignatureVerification(a.ApplicationPublicKey, messageHash, a.ApplicationSignature); err != nil {
		return InvalidTokenSignatureErorr
	}
	return nil
}
