package types

import (
	"encoding/hex"
	sdk "github.com/pokt-network/posmint/types"
	"reflect"
)

// RouterKey is the module name router key
const (
	RouterKey    = ModuleName
	MsgClaimName = "claim"
	MsgProofName = "proof"
)

// MsgClaim claims that you completed `NumOfProofs` and provides the merkle root for data integrity
type MsgClaim struct {
	SessionHeader `json:"header"` // header information for identification
	MerkleRoot    HashSum         `json:"merkle_root"`   // merkle root for data integrity
	TotalProofs   int64           `json:"total_relays"`  // total number of relays
	FromAddress   sdk.Address     `json:"from_address"`  // claimant
	EvidenceType  EvidenceType    `json:"evidence_type"` // relay or challenge?
}

func (msg MsgClaim) Route() string { return RouterKey }
func (msg MsgClaim) Type() string  { return MsgClaimName }
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
	if msg.EvidenceType == 0 {
		return NewNoEvidenceTypeErr(ModuleName)
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgClaim) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.FromAddress)}
}

// ---------------------------------------------------------------------------------------------------------------------

// MsgProof proves the previous claim by providing the merkle Proof and the leaf node
type MsgProof struct {
	MerkleProofs MerkleProofs `json:"merkle_proofs"` // the merkleProof needed to verify the proofs
	Leaf         Proof        `json:"leaf"`          // the needed to verify the Proof
	Cousin       Proof        `json:"cousin"`        // the cousin needed to verify the Proof
}

func (msg MsgProof) Route() string { return RouterKey }
func (msg MsgProof) Type() string  { return MsgProofName }
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
	if err := msg.Leaf.ValidateBasic(); err != nil {
		return err
	}
	if err := msg.Cousin.ValidateBasic(); err != nil {
		return err
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgProof) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgProof) GetSigners() []sdk.Address {
	return msg.Leaf.GetSigners()
}
