package service

import (
	"bytes"
	
	"github.com/pokt-network/pocket-core/const"
)

func ValidationHash(relayAnswer string) []byte {
	hasher := _const.VALIDATEHASHINGALGORITHM.New()
	hasher.Write([]byte(relayAnswer))
	return hasher.Sum(nil)
}

func Validate(relay Relay, hash []byte) (bool, error) {
	// complete relay locally
	relayAnswer, err := RouteRelay(relay)
	if err != nil {
		return false, err
	}
	// hash the local relay answer
	myHash := ValidationHash(relayAnswer)
	// compare your answer with their answer
	if bytes.Compare(myHash, hash) != 0 {
		return false, nil
	}
	// if the same return true
	return true, nil
}
