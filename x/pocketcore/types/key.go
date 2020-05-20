package types

import (
	sdk "github.com/pokt-network/posmint/types"
	"time"
)

const (
	ModuleName = "pocketcore"            // name of the module
	StoreKey   = ModuleName              // key for state store
	TStoreKey  = "transient_" + StoreKey // transient key for state store
)

var (
	ReceiptKey = []byte{0x01} // key for the verified and stored evidence
	ClaimKey   = []byte{0x02} // key for pending claims
)

// "KeyForReceipt" - Generates a key for the receipt object for the state store
func KeyForReceipt(ctx sdk.Ctx, addr sdk.Address, header SessionHeader, evidenceType EvidenceType) ([]byte, error) {
	// validate the header
	if err := header.ValidateHeader(); err != nil {
		return nil, err
	}
	// verify the address
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
	return append(append(append(ReceiptKey, addr.Bytes()...), header.Hash()...), et), nil
}

// "KeyForReceipts" - Generates a key for the receips object using an address
func KeyForReceipts(addr sdk.Address) ([]byte, error) {
	// verify the address passed
	if err := AddressVerification(addr.String()); err != nil {
		return nil, err
	}
	// return the key bz
	return append(ReceiptKey, addr.Bytes()...), nil
}

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
	return append(append(append(ClaimKey, addr.Bytes()...), header.Hash()...), et), nil
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
	defer timeTrack(time.Now(), "key for evidence")
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
