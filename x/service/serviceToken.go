package service

import (
	"encoding/hex"
	"github.com/pokt-network/pocket-core/crypto"
	"github.com/pokt-network/pocket-core/types"
)

type ServiceToken types.AAT

func (st ServiceToken) IsValid() error{
	publicKeyBytes, err := hex.DecodeString(st.AATMessage.ApplicationPublicKey)
	messageHash := st.Hash()
	if err != nil {
		return err
	}
	if !crypto.MockVerifySignature(publicKeyBytes, messageHash, st.Signature) {
		return InvalidSignatureError
	}
}

func (st ServiceToken) Hash() []byte {
	// TODO return hash of the amino encoding of service token
	return crypto.Hash([]byte(st.AATMessage.ApplicationPublicKey + st.AATMessage.ClientAddress)) // temporary
}
