package service

import (
	"github.com/pokt-network/pocket-core/crypto"
)

type IncrementCounter struct {
	Counter   int              `json:"counter"`
	Signature crypto.Signature `json:"signature"`
}

func (ic IncrementCounter) IsValid(applicationPubKey crypto.PublicKey) bool {
	// todo need signature verification crypto function
	return true
}
