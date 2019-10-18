package service

import (
	"encoding/hex"
	"github.com/pokt-network/pocket-core/crypto"
	"github.com/pokt-network/pocket-core/types"
)

type ServiceToken types.AAT

func (st ServiceToken) Validate() error {
	if err := st.ValidateVersion(); err != nil {
		return err
	}
	if err := st.ValidateMessage(); err != nil {
		return err
	}
	if err := st.ValidateSignature(); err != nil {
		return err
	}
	return nil
}

func (st ServiceToken) Hash() []byte {
	// TODO possibly return hash of the amino encoding of service token
	return crypto.SHA3FromString(st.AATMessage.ApplicationPublicKey + st.AATMessage.ClientPublicKey) // temporary
}

func (st ServiceToken) ValidateVersion() error {
	// check for valid version
	if !st.Version.IsIncluded() {
		return MissingTokenVersionError
	}
	if !st.Version.IsSupported() {
		return UnsupportedTokenVersionError
	}
	return nil
}

func (st ServiceToken) ValidateMessage() error {
	// check for valid application public key
	// todo pub key format verification
	if len(st.AATMessage.ApplicationPublicKey) == 0 {
		return MissingApplicationPublicKeyError
	}
	// todo pub key format verification
	if len(st.AATMessage.ClientPublicKey) == 0 {
		return MissingClientPublicKeyError
	}
	return nil
}

func (st ServiceToken) ValidateSignature() error {
	// check for valid signature
	messageHash := st.Hash()
	publicKeyBytes, err := hex.DecodeString(st.AATMessage.ApplicationPublicKey)
	if err != nil {
		return err
	}
	if !crypto.MockVerifySignature(publicKeyBytes, messageHash, st.Signature) {
		return InvalidTokenSignatureErorr
	}
	return nil
}
