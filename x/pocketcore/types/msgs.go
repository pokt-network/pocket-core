package types

import (
	"encoding/hex"
	"github.com/pokt-network/merkle"
	"github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
)

// RouterKey is the module name router key
const RouterKey = ModuleName

// MsgClaim claims that you completed `TotalRelays` and provides the merkle root for data integrity
type MsgClaim struct {
	SessionHeader `json:"header"` // header information for identification
	Root          []byte          `json:"root"`         // merkle root for data integrity
	TotalRelays   int64           `json:"total_relays"` // total number of relays
	FromAddress   sdk.ValAddress  `json:"from_address"` // claimant
}

func (msg MsgClaim) Route() string { return RouterKey }
func (msg MsgClaim) Type() string  { return "claim" }
func (msg MsgClaim) ValidateBasic() sdk.Error {
	// validate a non empty chain
	if msg.Chain == "" {
		return NewEmptyChainError(ModuleName)
	}
	// basic validation on the session block height
	if msg.SessionBlockHeight <= 1 {
		return NewEmptyBlockIDError(ModuleName)
	}
	// validate non negative total relays
	if msg.TotalRelays <= 0 {
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
	if err := HashVerification(hex.EncodeToString(msg.Root)); err != nil {
		return err
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgClaim) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.FromAddress)}
}

// ---------------------------------------------------------------------------------------------------------------------

// MsgProof proves the previous claim by providing the merkle proof and the leaf node
type MsgProof struct {
	merkle.Proof `json:"proof"` // the branch needed to verify the proofs
	LeafNode     Proof          `json:"leaf"` // the needed to verify the proof
}

func (msg MsgProof) Route() string { return RouterKey }
func (msg MsgProof) Type() string  { return "claim" }
func (msg MsgProof) ValidateBasic() sdk.Error {
	// verify non empty merkle proof
	if len(msg.Proof.Hashes) == 0 {
		return NewEmtpyBranchError(ModuleName)
	}
	// verify the session block height is positive
	if msg.LeafNode.SessionBlockHeight < 0 {
		return NewInvalidBlockHeightError(ModuleName)
	}
	// verify the public key format
	if err := PubKeyVerification(msg.LeafNode.ServicerPubKey); err != nil {
		return err
	}
	// verify the blockchain hash format
	if err := HashVerification(msg.LeafNode.Blockchain); err != nil {
		return err
	}
	// verify non negative index
	if msg.LeafNode.Index < 0 {
		return NewInvalidIncrementCounterError(ModuleName)
	}
	// verify a valid token
	if err := msg.LeafNode.Token.Validate(); err != nil {
		return NewInvalidTokenError(ModuleName, err)
	}
	// verify the client signature on the proof
	if err := SignatureVerification(msg.LeafNode.Signature, msg.LeafNode.HashString(), msg.LeafNode.Signature); err != nil {
		return err
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgProof) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgProof) GetSigners() []sdk.AccAddress {
	pk, err := crypto.NewPublicKey(msg.LeafNode.ServicerPubKey)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sdk.AccAddress(pk.Address())}
}
