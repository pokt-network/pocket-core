package types

import (
	"encoding/hex"
	sdk "github.com/pokt-network/posmint/types"
)

const (
	// ModuleName is the name of the module
	ModuleName = "pocketcore"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName
)

var (
	ProofSummaryKey = []byte{0x01} // key for the proofSummary
)

func KeyForNodeProofSummary(addr sdk.ValAddress, header ProofsHeader) []byte {
	appPubKey, err := hex.DecodeString(header.ApplicationPubKey)
	if err != nil {
		panic(err)
	}
	sessionHash, err := hex.DecodeString(header.SessionBlockHash)
	if err != nil {
		panic(err)
	}
	return append(append(append(ProofSummaryKey, addr.Bytes()...), appPubKey...), sessionHash...)
}

func KeyForNodeProofSummaries(addr sdk.ValAddress) []byte {
	return append(ProofSummaryKey, addr.Bytes()...)
}

func KeyForNodeProofSummariesForApp(addr sdk.ValAddress, appPubKeyHex string) []byte {
	appPubKey, err := hex.DecodeString(appPubKeyHex)
	if err != nil {
		panic(err)
	}
	return append(append(ProofSummaryKey, addr.Bytes()...), appPubKey...)
}
