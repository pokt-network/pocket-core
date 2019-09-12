package service

import (
	"encoding/hex"
	"github.com/pokt-network/pocket-core/crypto"
	"github.com/pokt-network/pocket-core/types"
)

type ServiceToken types.AAT

func (st ServiceToken) IsValid() error {
	if err := st.VersionIsValid(); err != nil {
		return err
	}
	if err := st.MessageIsValid(); err != nil {
		return err
	}
	if err := st.SignatureIsValid(); err != nil {
		return err
	}
	// run session algorithm to authenticate service
	// todo 'how should dependencies work here?
	// if (not part of session){
	//    return InvalidSessionError
	// }
	return nil
}

func (st ServiceToken) Hash() []byte {
	// TODO return hash of the amino encoding of service token
	return crypto.Hash([]byte(st.AATMessage.ApplicationPublicKey + st.AATMessage.ClientPublicKey)) // temporary
}

func (st ServiceToken) VersionIsValid() error{
	// check for valid version
	if !st.Version.IsIncluded() {
		return MissingTokenVersionError
	}
	if !st.Version.IsSupported() {
		return UnsupportedTokenVersionError
	}
	return nil
}

func (st ServiceToken) MessageIsValid() error {
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

func (st ServiceToken) SignatureIsValid() error {
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
