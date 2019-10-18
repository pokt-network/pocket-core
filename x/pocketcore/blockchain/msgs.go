package blockchain

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// todo possibly remove

// RouterKey is the module name router key
const RouterKey = ModuleName // this was defined in your key.go file

// MsgRelayBatch defines a SetName message
type MsgRelayBatch struct{}

// NewMsgSetName is a constructor function for MsgRelayBatch
func NewMsgSetName() MsgRelayBatch {
	return MsgRelayBatch{}
}

// Route should return the name of the module
func (msg MsgRelayBatch) Route() string { return RouterKey }

// Type should return the action
func (msg MsgRelayBatch) Type() string { return "relay_batch" }

// ValidateBasic runs stateless checks on the message
func (msg MsgRelayBatch) ValidateBasic() sdk.Error {
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgRelayBatch) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgRelayBatch) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}
