package types

import "github.com/pokt-network/pocket-core/x/service"

type RelayBatch struct {
	RelayBatchHeader
	Evidence
}

type RelayBatchHeader struct {
	SessionHash       string
	ApplicationPubKey string
}

type RelayBatches map[RelayBatchHeader]RelayBatch // public key of Application -> RelayBatch

type Evidence []service.ServiceAuthentication
