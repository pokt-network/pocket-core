package types

import (
	sdk "github.com/pokt-network/pocket-core/types"
)

const (
	ModuleName = "pocketcore"            // name of the module
	StoreKey   = ModuleName              // key for state store
	TStoreKey  = "transient_" + StoreKey // transient key for state store
)

var (
	ClaimKey = []byte{0x02} // key for pending claims
)

var ClaimLen = len(ClaimKey)

// "KeyForClaim" - Generates the key for the claim object for the state store
func KeyForClaim(ctx sdk.Ctx, addr sdk.Address, header SessionHeader, evidenceType EvidenceType) ([]byte, error) {
	// validat the header
	if err := header.ValidateHeader(); err != nil {
		return nil, err
	}
	// validate the address
	if err := AddressVerification(addr.String()); err != nil {
		return nil, err
	}
	// validate the evidence type
	if evidenceType != RelayEvidence && evidenceType != ChallengeEvidence {
		return nil, NewInvalidEvidenceErr(ModuleName)
	}
	et, err := evidenceType.Byte()
	if err != nil {
		return nil, err
	}
	// return the key bz
	b := make([]byte, HashLength+ClaimLen+sdk.AddrLen+1)
	return MultiAppend(b, ClaimKey, addr.Bytes(), header.Hash(), []byte{et}), nil
}

// "KeyForClaims" - Generates the key for the claims object
func KeyForClaims(addr sdk.Address) ([]byte, error) {
	// verify the address
	if err := AddressVerification(addr.String()); err != nil {
		return nil, err
	}
	// return the key bz
	return append(ClaimKey, addr.Bytes()...), nil
}

// "KeyForEvidence" - Generates the key for evidence
func KeyForEvidence(header SessionHeader, evidenceType EvidenceType) ([]byte, error) {
	// validate the evidence type
	if evidenceType != RelayEvidence && evidenceType != ChallengeEvidence {
		return nil, NewInvalidEvidenceErr(ModuleName)
	}
	et, err := evidenceType.Byte()
	if err != nil {
		return nil, err
	}
	return append(header.Hash(), et), nil
}
