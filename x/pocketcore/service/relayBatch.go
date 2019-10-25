package service

import (
	"github.com/pokt-network/pocket-core/types"
	"sync"
)

// todo abstract so that POS module can use
type RelayBatch struct {
	ProofsHeader
	Proofs
}

type RelayBatches types.List

var (
	globalRelayBatches *RelayBatches // map[ProofsHeader]RelayBatch // [ProofsHeader] -> RelayBatch
	relayBatchOnce     sync.Once
)

func GetGlobalRelayBatches() *RelayBatches {
	relayBatchOnce.Do(func() {
		globalRelayBatches = (*RelayBatches)(types.NewList())
	})
	return globalRelayBatches
}

func (rbs *RelayBatches) AddProof(authentication ServiceProof, sessionBlockIDHex string, maxNumberOfRelays int) error {
	(*types.List)(rbs).Mux.Lock()
	defer (*types.List)(rbs).Mux.Unlock()
	rbh := ProofsHeader{
		SessionHash:       sessionBlockIDHex,
		ApplicationPubKey: authentication.ServiceToken.AATMessage.ApplicationPublicKey,
	}
	if relayBatch, contains := rbs.M[rbh]; contains {
		return relayBatch.(RelayBatch).Proofs.AddProof(authentication)
	} else {
		return rbs.NewRelayBatch(authentication, sessionBlockIDHex, maxNumberOfRelays)
	}
}

func (rbs *RelayBatches) NewRelayBatch(authentication ServiceProof, latestSessionBlockHex string, maxNumberOfRelays int) error {
	rb := RelayBatch{
		ProofsHeader: ProofsHeader{
			SessionHash:       latestSessionBlockHex,
			ApplicationPubKey: authentication.ServiceToken.AATMessage.ApplicationPublicKey,
		},
		Proofs: make([]ServiceProof, maxNumberOfRelays),
	}
	err := rb.AddProof(authentication)
	if err != nil {
		return NewRelayBatchCreationError(err)
	}
	return nil
}

func (rbs *RelayBatches) AddBatch(batch RelayBatch) {
	(*types.List)(rbs).Add(batch.ProofsHeader, batch)
}

func (rbs *RelayBatches) Getbatch(relayBatchHeader ProofsHeader) {
	(*types.List)(rbs).Get(relayBatchHeader)
}

func (rbs *RelayBatches) Removebatch(relayBatchHeader ProofsHeader) {
	(*types.List)(rbs).Remove(relayBatchHeader)
}

func (rbs *RelayBatches) Len() int {
	return (*types.List)(rbs).Count()
}

func (rbs *RelayBatches) Contains(relayBatchHeader ProofsHeader) bool {
	return (*types.List)(rbs).Contains(relayBatchHeader)
}

func (rbs *RelayBatches) Clear() {
	(*types.List)(rbs).Clear()
}
