package types

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/pokt-network/pocket-core/types"
)

// RouterKey is the module name router key
const (
	RouterKey    = ModuleName // router name is module name
	MsgClaimName = "claim"    // name for the claim message
	MsgProofName = "proof"    // name for the proof message
)

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
	if msg.SessionHeader.Chain == "" {
		return NewEmptyChainError(ModuleName)
	}
	// basic validation on the session block height
	if msg.SessionHeader.SessionBlockHeight < 1 {
		return NewEmptyBlockIDError(ModuleName)
	}
	// validate greater than 5 relays (need 5 for the tree structure)
	if msg.TotalProofs < 5 {
		return NewEmptyProofsError(ModuleName)
	}
	// validate the public key format
	if err := PubKeyVerification(msg.SessionHeader.ApplicationPubKey); err != nil {
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
	// ensure non zero root upper range
	if !msg.MerkleRoot.isValidRange() {
		return NewInvalidMerkleRangeError(ModuleName)
	}
	// ensure zero root lower range
	if msg.MerkleRoot.Range.Lower != 0 {
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

// "GetFee" - Returns the fee (sdk.Int) of the messgae type
func (msg MsgProtoProof) GetFee() sdk.Int {
	return sdk.NewInt(PocketFeeMap[msg.Type()])
}

// "Route" - Returns module router key
func (msg MsgProtoProof) Route() string { return RouterKey }

// "Type" - Returns message name
func (msg MsgProtoProof) Type() string { return MsgProofName }

// "ValidateBasic" - Storeless validity check for proof message
func (msg MsgProtoProof) ValidateBasic() sdk.Error {
	// verify valid number of levels for merkle proofs
	if len(msg.MerkleProof.HashRanges) < 3 {
		return NewInvalidLeafCousinProofsComboError(ModuleName)
	}
	// validate the target range
	if !msg.MerkleProof.Target.isValidRange() {
		return NewInvalidMerkleRangeError(ModuleName)
	}
	// validate the leaf
	if err := msg.GetLeaf().ValidateBasic(); err != nil {
		return err
	}
	if _, err := msg.EvidenceType.Byte(); err != nil {
		return NewInvalidEvidenceErr(ModuleName)
	}
	return nil
}

// "GetSignBytes" - Encodes the message for signing
func (msg MsgProtoProof) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners defines whose signature is required
func (msg MsgProtoProof) GetSigner() sdk.Address {
	return msg.GetLeaf().GetSigner()
}

func (msg MsgProtoProof) GetLeaf() Proof {
	return msg.Leaf.FromProto()
}

// Legacy Amino Msg Below
//----------------------------------------------------------------------------------------------------------------------

// "MsgProof" - Proves the previous claim by providing the merkle Proof and the leaf node
type MsgProof struct {
	MerkleProof  MerkleProof  `json:"merkle_proofs"` // the merkleProof needed to verify the proofs
	Leaf         Proof        `json:"leaf"`          // the needed to verify the Proof
	EvidenceType EvidenceType `json:"evidence_type"` // the type of evidence
}

func (msg MsgProof) Reset() {
	panic("amino only msg")
}

func (msg MsgProof) String() string {
	return fmt.Sprintf("MerkleProof: %s\nLeaf: %v\nEvidenceType: %d\n", msg.MerkleProof.String(), msg.Leaf, msg.EvidenceType)
}

func (msg MsgProof) ProtoMessage() {
	panic("amino only msg")
}

func (msg MsgProof) ToProto() MsgProtoProof {
	return MsgProtoProof{
		MerkleProof:  msg.MerkleProof,
		Leaf:         msg.Leaf.ToProto(),
		EvidenceType: msg.EvidenceType,
	}
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
	if len(msg.MerkleProof.HashRanges) < 3 {
		return NewInvalidLeafCousinProofsComboError(ModuleName)
	}
	// validate the target range
	if !msg.MerkleProof.Target.isValidRange() {
		return NewInvalidMerkleRangeError(ModuleName)
	}
	// validate the leaf
	if err := msg.Leaf.ValidateBasic(); err != nil {
		return err
	}
	if _, err := msg.EvidenceType.Byte(); err != nil {
		return NewInvalidEvidenceErr(ModuleName)
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
