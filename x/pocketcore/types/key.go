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

func KeyForPOR(appPubKey, chain, sessionHeight string) string {
	return appPubKey + chain + sessionHeight
}

func KeyForProofOfRelay(ctx sdk.Context, addr sdk.ValAddress, header PORHeader) []byte {
	appPubKey, err := hex.DecodeString(header.ApplicationPubKey)
	if err != nil {
		panic(err)
	}
	sessionHash := ctx.WithBlockHeight(header.SessionBlockHeight).BlockHeader().GetLastBlockId().Hash
	return append(append(append(ProofSummaryKey, addr.Bytes()...), appPubKey...), sessionHash...)
}

func KeyForProofOfRelays(addr sdk.ValAddress) []byte {
	return append(ProofSummaryKey, addr.Bytes()...)
}

func KeyForProofOfRelaysApp(addr sdk.ValAddress, appPubKeyHex string) []byte {
	appPubKey, err := hex.DecodeString(appPubKeyHex)
	if err != nil {
		panic(err)
	}
	return append(append(ProofSummaryKey, addr.Bytes()...), appPubKey...)
}
