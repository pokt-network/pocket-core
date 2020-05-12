package types

import (
	"encoding/hex"
	"reflect"

	sdk "github.com/pokt-network/posmint/types"
)

// RouterKey is the module name router key
const (
	RouterKey    = ModuleName // router name is module name
	MsgClaimName = "claim"    // name for the claim message
	MsgProofName = "proof"    // name for the proof message
)

// "MsgClaim" - claims that you completed `NumOfProofs` for relay or challenge and provides the merkle root for data integrity
type MsgClaim struct {
	SessionHeader    `json:"header"` // header information for identification
	MerkleRoot       HashSum         `json:"merkle_root"`   // merkle root for data integrity
	TotalProofs      int64           `json:"total_relays"`  // total number of relays
	FromAddress      sdk.Address     `json:"from_address"`  // claimant's address
	EvidenceType     EvidenceType    `json:"evidence_type"` // relay or challenge?
	ExpirationHeight int64           `json:"expiration_height"`
}

// "GetFee" - Returns the fee (sdk.Int) of the messgae type
func (msg MsgClaim) GetFee() sdk.Int {
	return sdk.NewInt(PocketFeeMap[msg.Type()])
}

// "Route" - Returns module router key
func (msg MsgClaim) Route() string { return RouterKey }

// "Type" - Returns message name
func (msg MsgClaim) Type() string { return MsgClaimName }

// "ValidateBasic" - Storeless validity check for claim message
func (msg MsgClaim) ValidateBasic() sdk.Error {
	// validate a non empty chain
	if msg.Chain == "" {
		return NewEmptyChainError(ModuleName)
	}
	// basic validation on the session block height
	if msg.SessionBlockHeight < 1 {
		return NewEmptyBlockIDError(ModuleName)
	}
	// validate greater than 5 relays (need 5 for the tree structure)
	if msg.TotalProofs < 5 {
		return NewEmptyProofsError(ModuleName)
	}
	// validate the public key format
	if err := PubKeyVerification(msg.ApplicationPubKey); err != nil {
		return NewPubKeyError(ModuleName, err)
	}
	// validate the address format
	if err := AddressVerification(msg.FromAddress.String()); err != nil {
		return NewInvalidHashError(ModuleName, err)
	}
	// validate the root format
	if err := HashVerification(hex.EncodeToString(msg.MerkleRoot.Hash)); err != nil {
		return err
	}
	// ensure non zero root sum
	if msg.MerkleRoot.Sum == 0 {
		return NewInvalidRootError(ModuleName)
	}
	// ensure non zero evidence
	if msg.EvidenceType == 0 {
		return NewNoEvidenceTypeErr(ModuleName)
	}
	if msg.EvidenceType != RelayEvidence && msg.EvidenceType != ChallengeEvidence {
		return NewInvalidEvidenceErr(ModuleName)
	}
	if msg.ExpirationHeight != 0 {
		return NewInvalidExpirationHeightErr(ModuleName)
	}
	return nil
}

// "GetSignBytes" - Encodes the message for signing
func (msg MsgClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// "GetSigners" - Defines whose signature is required
func (msg MsgClaim) GetSigner() sdk.Address {
	return msg.FromAddress
}

// "IsEmpty" - Returns true if the EvidenceType == 0, this should only happen on initialization and MsgClaim{} calls
func (msg MsgClaim) IsEmpty() bool {
	return msg.EvidenceType == 0
}

// ---------------------------------------------------------------------------------------------------------------------

// "MsgProof" - Proves the previous claim by providing the merkle Proof and the leaf node
type MsgProof struct {
	MerkleProofs MerkleProofs `json:"merkle_proofs"` // the merkleProof needed to verify the proofs
	Leaf         Proof        `json:"leaf"`          // the needed to verify the Proof
	Cousin       Proof        `json:"cousin"`        // the cousin needed to verify the Proof
}

// "GetFee" - Returns the fee (sdk.Int) of the messgae type
func (msg MsgProof) GetFee() sdk.Int {
	return sdk.NewInt(PocketFeeMap[msg.Type()])
}

// "Route" - Returns module router key
func (msg MsgProof) Route() string { return RouterKey }

// "Type" - Returns message name
func (msg MsgProof) Type() string { return MsgProofName }

// "ValidateBasic" - Storeless validity check for proof message
func (msg MsgProof) ValidateBasic() sdk.Error {
	// verify valid number of levels for merkle proofs
	if len(msg.MerkleProofs[0].HashSums) < 3 || len(msg.MerkleProofs[0].HashSums) != len(msg.MerkleProofs[1].HashSums) {
		return NewInvalidLeafCousinProofsComboError(ModuleName)
	}
	// ensure the two indices are not equal
	if msg.MerkleProofs[0].Index == msg.MerkleProofs[1].Index {
		return NewInvalidLeafCousinProofsComboError(ModuleName)
	}
	// ensure leaf does not equal cousin
	if reflect.DeepEqual(msg.Leaf, msg.Cousin) {
		return NewCousinLeafEquivalentError(ModuleName)
	}
	// ensure leaf relayProof does not equal cousin relayProof
	if reflect.DeepEqual(msg.MerkleProofs[0].HashSums, msg.MerkleProofs[1].HashSums) {
		return NewCousinLeafEquivalentError(ModuleName)
	}
	// validate the leaf
	if err := msg.Leaf.ValidateBasic(); err != nil {
		return err
	}
	// validate the cousin
	if err := msg.Cousin.ValidateBasic(); err != nil {
		return err
	}
	return nil
}

// "GetSignBytes" - Encodes the message for signing
func (msg MsgProof) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgProof) GetSigner() sdk.Address {
	return msg.Leaf.GetSigner()
}
