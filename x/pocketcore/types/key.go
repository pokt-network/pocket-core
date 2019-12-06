package types

import (
	"encoding/hex"
	sdk "github.com/pokt-network/posmint/types"
)

const (
	ModuleName = "pocketcore"
	StoreKey   = ModuleName
)

var (
	ProofKey           = []byte{0x01} // key for the verified proofs
	UnverifiedProofKey = []byte{0x02} // key for non-verified proofs
)

func KeyForProof(ctx sdk.Context, addr sdk.ValAddress, header SessionHeader) []byte {
	appPubKey, err := hex.DecodeString(header.ApplicationPubKey)
	if err != nil {
		panic(err)
	}
	sessionHash := ctx.WithBlockHeight(header.SessionBlockHeight).BlockHeader().GetLastBlockId().Hash
	return append(append(append(ProofKey, addr.Bytes()...), appPubKey...), sessionHash...)
}

func KeyForProofs(addr sdk.ValAddress) []byte {
	return append(ProofKey, addr.Bytes()...)
}

func KeyForUnverifiedProof(ctx sdk.Context, addr sdk.ValAddress, header SessionHeader) []byte {
	appPubKey, err := hex.DecodeString(header.ApplicationPubKey)
	if err != nil {
		panic(err)
	}
	sessionHash := ctx.WithBlockHeight(header.SessionBlockHeight).BlockHeader().GetLastBlockId().Hash
	return append(append(append(UnverifiedProofKey, addr.Bytes()...), appPubKey...), sessionHash...)
}

func KeyForUnverifiedProofs(addr sdk.ValAddress) []byte {
	return append(UnverifiedProofKey, addr.Bytes()...)
}
