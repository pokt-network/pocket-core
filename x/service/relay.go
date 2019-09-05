package service

import (
	"errors"
	"github.com/pokt-network/pocket-core/crypto"
)

type Relay struct {
	Blockchain   ServiceBlockchain `json:"blockchain"`
	Payload      ServicePayload    `json:"payload"`
	ServiceToken ServiceToken      `json:"token"`
}

func (r Relay)IsValid(hostedBlockchains ServiceBlockchains) error{
	if r.Blockchain == nil || len(r.Blockchain) == 0{
		return EmptyBlockchainError
	}
	if !hostedBlockchains.Contains(r.Blockchain) {
		return UnsupportedBlockchainError
	}
	if err := r.ServiceToken.IsValid(); err !=nil{
		return errors.New(InvalidTokenError.error + " : " + err.Error())
	}
	return nil
}
