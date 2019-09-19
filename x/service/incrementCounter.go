package service

import (
	"encoding/hex"
	"errors"
	"github.com/pokt-network/pocket-core/crypto"
)

type IncrementCounter struct {
	Counter   int              `json:"counter"`
	Signature crypto.Signature `json:"signature"`
}

func (ic IncrementCounter) Validate(clientPubKey string, messageHash []byte) error {
	// check if counter is valid
	// todo
	if ic.Counter < 0 {
		return NegativeICCounterError
	}
	cpkBytes, err := hex.DecodeString(clientPubKey)
	if err != nil {
		return errors.New(ClientPubKeyDecodeError.Error() + " : " + err.Error())
	}
	if !crypto.MockVerifySignature(cpkBytes, messageHash, ic.Signature) {
		return InvalidICSignatureError
	}
	return nil
}
