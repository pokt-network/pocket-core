package types

import (
	"encoding/hex"
	sdk "github.com/pokt-network/posmint/types"
)

const (
	ModuleName = "pocketcore"
	StoreKey   = ModuleName
	TStoreKey  = "transient_" + StoreKey
)

var (
	EvidenceKey = []byte{0x01} // key for the verified proofs
	ClaimKey   = []byte{0x02} // key for non-verified proofs
)

func KeyForEvidence(ctx sdk.Context, addr sdk.Address, header SessionHeader) ([]byte, error) {
	if err := header.ValidateHeader(); err != nil {
		return nil, err
	}
	if err := AddressVerification(addr.String()); err != nil {
		return nil, err
	}
	appPubKey, err := hex.DecodeString(header.ApplicationPubKey)
	if err != nil {
		return nil, err
	}
	sessionCtx := ctx.MustGetPrevCtx(header.SessionBlockHeight)
	sessionBlockHeader := sessionCtx.BlockHeader()
	sessionHash := sessionBlockHeader.GetLastBlockId().Hash
	return append(append(append(EvidenceKey, addr.Bytes()...), appPubKey...), sessionHash...), nil
}

func KeyForEvidences(addr sdk.Address) ([]byte, error) {
	if err := AddressVerification(addr.String()); err != nil {
		return nil, err
	}
	return append(EvidenceKey, addr.Bytes()...), nil
}

func KeyForClaim(ctx sdk.Context, addr sdk.Address, header SessionHeader) ([]byte, error) {
	if err := header.ValidateHeader(); err != nil {
		return nil, err
	}
	if err := AddressVerification(addr.String()); err != nil {
		return nil, err
	}
	appPubKey, err := hex.DecodeString(header.ApplicationPubKey)
	if err != nil {
		return nil, err
	}
	sessionCtx := ctx.MustGetPrevCtx(header.SessionBlockHeight)
	sessionBlockHeader := sessionCtx.BlockHeader()
	sessionHash := sessionBlockHeader.GetLastBlockId().Hash
	return append(append(append(ClaimKey, addr.Bytes()...), appPubKey...), sessionHash...), nil
}

func KeyForClaims(addr sdk.Address) ([]byte, error) {
	if err := AddressVerification(addr.String()); err != nil {
		return nil, err
	}
	return append(ClaimKey, addr.Bytes()...), nil
}
