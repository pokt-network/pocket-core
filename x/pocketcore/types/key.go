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
	InvoiceKey = []byte{0x01} // key for the verified proofs
	ClaimKey   = []byte{0x02} // key for non-verified proofs
)

func KeyForInvoice(ctx sdk.Context, addr sdk.ValAddress, header SessionHeader) []byte {
	if err := header.ValidateHeader(); err != nil {
		panic(err)
	}
	if err := AddressVerification(addr.String()); err != nil {
		panic(err)
	}
	appPubKey, err := hex.DecodeString(header.ApplicationPubKey)
	if err != nil {
		panic(err)
	}
	sessionCtx := ctx.WithBlockHeight(header.SessionBlockHeight)
	sessionBlockHeader := sessionCtx.BlockHeader()
	sessionHash := sessionBlockHeader.GetLastBlockId().Hash
	return append(append(append(InvoiceKey, addr.Bytes()...), appPubKey...), sessionHash...)
}

func KeyForInvoices(addr sdk.ValAddress) []byte {
	if err := AddressVerification(addr.String()); err != nil {
		panic(err)
	}
	return append(InvoiceKey, addr.Bytes()...)
}

func KeyForClaim(ctx sdk.Context, addr sdk.ValAddress, header SessionHeader) []byte {
	if err := header.ValidateHeader(); err != nil {
		panic(err)
	}
	if err := AddressVerification(addr.String()); err != nil {
		panic(err)
	}
	appPubKey, err := hex.DecodeString(header.ApplicationPubKey)
	if err != nil {
		panic(err)
	}
	sessionCtx := ctx.WithBlockHeight(header.SessionBlockHeight)
	sessionBlockHeader := sessionCtx.BlockHeader()
	sessionHash := sessionBlockHeader.GetLastBlockId().Hash
	return append(append(append(ClaimKey, addr.Bytes()...), appPubKey...), sessionHash...)
}

func KeyForClaims(addr sdk.ValAddress) []byte {
	if err := AddressVerification(addr.String()); err != nil {
		panic(err)
	}
	return append(ClaimKey, addr.Bytes()...)
}
