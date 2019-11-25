package types

import (
	sdk "github.com/pokt-network/posmint/types"
)

// RouterKey is the module name router key
const RouterKey = ModuleName

// MsgProofOfRelays defines a SetName message
type MsgProofOfRelays struct {
	ProofSummary
	ProofBatch
}

// NewMsgSetName is a constructor function for MsgProofOfRelays
func NewMsgProofBatch(pb ProofBatch, ps ProofSummary) MsgProofOfRelays {
	return MsgProofOfRelays{
		ProofSummary: ps,
		ProofBatch:   pb,
	}
}

// Route should return the name of the module
func (msg MsgProofOfRelays) Route() string { return RouterKey }

// Type should return the action
func (msg MsgProofOfRelays) Type() string { return "relay_batch" }

// ValidateBasic runs stateless checks on the message
func (msg MsgProofOfRelays) ValidateBasic() sdk.Error {
	if msg.ProofsHeader.Chain == "" {
		return NewEmptyChainError(ModuleName)
	}
	if msg.ProofsHeader.SessionBlockHash == "" {
		return NewEmptyBlockIDError(ModuleName)
	}
	if msg.ProofsHeader.ApplicationPubKey == "" {
		return NewEmptyAppPubKeyError(ModuleName)
	}
	if len(msg.Proofs) == 0 {
		return NewEmptyProofsError(ModuleName)
	}
	if msg.RelaysCompleted < 1 {
		return NewInvalidRelaysCompletedError(ModuleName)
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgProofOfRelays) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgProofOfRelays) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}
