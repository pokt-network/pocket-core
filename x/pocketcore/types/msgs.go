package types

import (
	"github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
)

// RouterKey is the module name router key
const RouterKey = ModuleName

// MsgProof defines a SetName message
type MsgProof struct {
	Header                     // header information
	Root        []byte         // root
	TotalRelays int64          // total number of relays
	FromAddress sdk.ValAddress // todo remove use ProofOfRelay -> nodePubKey.Address()
}

// Route should return the name of the module
func (msg MsgProof) Route() string { return RouterKey }

// Type should return the action
func (msg MsgProof) Type() string { return "proof" }

// ValidateBasic runs stateless checks on the message
func (msg MsgProof) ValidateBasic() sdk.Error {
	if msg.Chain == "" {
		return NewEmptyChainError(ModuleName)
	}
	if msg.SessionBlockHeight <= 1 {
		return NewEmptyBlockIDError(ModuleName)
	}
	if msg.ApplicationPubKey == "" {
		return NewEmptyAppPubKeyError(ModuleName)
	}
	if msg.TotalRelays <= 0 {
		return NewEmptyProofsError(ModuleName)
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgProof) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgProof) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.FromAddress)}
}

// ---------------------------------------------------------------------------------------------------------------------

// MsgProof defines a SetName message
type MsgClaimProof struct {
	MerkleProof    // the branch needed to verify the proofs
	LeafNode Proof // the needed to verify the proof
}

// Route should return the name of the module
func (msg MsgClaimProof) Route() string { return RouterKey }

// Type should return the action
func (msg MsgClaimProof) Type() string { return "claim" }

// ValidateBasic runs stateless checks on the message
func (msg MsgClaimProof) ValidateBasic() sdk.Error {
	if len(msg.MerkleProof) == 0 {
		return NewEmtpyBranchError(ModuleName)
	}
	if msg.LeafNode.SessionBlockHeight < 0 {
		return NewInvalidBlockHeightError(ModuleName)
	}
	if err := PubKeyVerification(msg.LeafNode.ServicerPubKey); err != nil {
		return err
	}
	if err := HashVerification(msg.LeafNode.Blockchain); err != nil {
		return err
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgClaimProof) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgClaimProof) GetSigners() []sdk.AccAddress {
	pk, err := crypto.NewPublicKey(msg.LeafNode.ServicerPubKey)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sdk.AccAddress(pk.Address())}
}
