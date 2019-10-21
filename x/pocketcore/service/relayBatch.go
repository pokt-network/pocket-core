package service

import (
	"github.com/pokt-network/pocket-core/types"
	"sync"
)

// todo abstract so that POS module can use
type RelayBatch struct {
	CertificatesHeader
	Certificates
}

type RelayBatches types.List

var (
	globalRelayBatches *RelayBatches // map[CertificatesHeader]RelayBatch // [CertificatesHeader] -> RelayBatch
	relayBatchOnce     sync.Once
)

func GetGlobalRelayBatches() *RelayBatches {
	relayBatchOnce.Do(func() {
		globalRelayBatches = (*RelayBatches)(types.NewList())
	})
	return globalRelayBatches
}

func (rbs *RelayBatches) AddCertificate(authentication ServiceCertificate, sessionBlockIDHex string, maxNumberOfRelays int) error {
	(*types.List)(rbs).Mux.Lock()
	defer (*types.List)(rbs).Mux.Unlock()
	rbh := CertificatesHeader{
		SessionHash:       sessionBlockIDHex,
		ApplicationPubKey: authentication.ServiceToken.AATMessage.ApplicationPublicKey,
	}
	if relayBatch, contains := rbs.M[rbh]; contains {
		return relayBatch.(RelayBatch).Certificates.AddCertificate(authentication)
	} else {
		return rbs.NewRelayBatch(authentication, sessionBlockIDHex, maxNumberOfRelays)
	}
}

func (rbs *RelayBatches) NewRelayBatch(authentication ServiceCertificate, latestSessionBlockHex string, maxNumberOfRelays int) error {
	rb := RelayBatch{
		CertificatesHeader: CertificatesHeader{
			SessionHash:       latestSessionBlockHex,
			ApplicationPubKey: authentication.ServiceToken.AATMessage.ApplicationPublicKey,
		},
		Certificates: make([]ServiceCertificate, maxNumberOfRelays),
	}
	err := rb.AddCertificate(authentication)
	if err != nil {
		return NewRelayBatchCreationError(err)
	}
	return nil
}

func (rbs *RelayBatches) AddBatch(batch RelayBatch) {
	(*types.List)(rbs).Add(batch.CertificatesHeader, batch)
}

func (rbs *RelayBatches) Getbatch(relayBatchHeader CertificatesHeader) {
	(*types.List)(rbs).Get(relayBatchHeader)
}

func (rbs *RelayBatches) Removebatch(relayBatchHeader CertificatesHeader) {
	(*types.List)(rbs).Remove(relayBatchHeader)
}

func (rbs *RelayBatches) Len() int {
	return (*types.List)(rbs).Count()
}

func (rbs *RelayBatches) Contains(relayBatchHeader CertificatesHeader) bool {
	return (*types.List)(rbs).Contains(relayBatchHeader)
}

func (rbs *RelayBatches) Clear() {
	(*types.List)(rbs).Clear()
}
