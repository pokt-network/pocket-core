package service

import (
	"github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/pocketcore/blockchain"
	"sync"
)

// todo abstract so that POS module can use
type RelayBatch struct {
	EvidenceHeader
	Evidence
}

type RelayBatches types.List

var (
	globalRelayBatches *RelayBatches // map[EvidenceHeader]RelayBatch // [EvidenceHeader] -> RelayBatch
	relayBatchOnce     sync.Once
)

func GetGlobalRelayBatches() *RelayBatches {
	relayBatchOnce.Do(func() {
		globalRelayBatches = (*RelayBatches)(types.NewList())
	})
	return globalRelayBatches
}

func (rbs *RelayBatches) AddEvidence(authentication ServiceCertificate) error {
	(*types.List)(rbs).Mux.Lock()
	defer (*types.List)(rbs).Mux.Unlock()
	rbh := EvidenceHeader{
		SessionHash:       blockchain.GetLatestSessionBlock().HashHex(),
		ApplicationPubKey: authentication.ServiceToken.AATMessage.ApplicationPublicKey,
	}
	if relayBatch, contains := rbs.M[rbh]; contains {
		return relayBatch.(RelayBatch).Evidence.AddEvidence(authentication)
	} else {
		return rbs.NewRelayBatch(authentication)
	}
}

func (rbs *RelayBatches) NewRelayBatch(authentication ServiceCertificate) error {
	rb := RelayBatch{
		EvidenceHeader: EvidenceHeader{
			SessionHash:       blockchain.GetLatestSessionBlock().HashHex(),
			ApplicationPubKey: authentication.ServiceToken.AATMessage.ApplicationPublicKey,
		},
		Evidence: make([]ServiceCertificate, blockchain.GetMaxNumberOfRelaysForApp(authentication.ServiceToken.AATMessage.ApplicationPublicKey)),
	}
	err := rb.AddEvidence(authentication)
	if err != nil {
		return NewRelayBatchCreationError(err)
	}
	return nil
}

func (rbs *RelayBatches) AddBatch(batch RelayBatch) {
	(*types.List)(rbs).Add(batch.EvidenceHeader, batch)
}

func (rbs *RelayBatches) Getbatch(relayBatchHeader EvidenceHeader) {
	(*types.List)(rbs).Get(relayBatchHeader)
}

func (rbs *RelayBatches) Removebatch(relayBatchHeader EvidenceHeader) {
	(*types.List)(rbs).Remove(relayBatchHeader)
}

func (rbs *RelayBatches) Len() int {
	return (*types.List)(rbs).Count()
}

func (rbs *RelayBatches) Contains(relayBatchHeader EvidenceHeader) bool {
	return (*types.List)(rbs).Contains(relayBatchHeader)
}

func (rbs *RelayBatches) Clear() {
	(*types.List)(rbs).Clear()
}
