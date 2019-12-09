package types

import (
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/pokt-network/posmint/types"
)

// ensure Msg interface compliance at compile time
var (
	_ sdk.Msg = &MsgStake{}
	_ sdk.Msg = &MsgBeginUnstake{}
	_ sdk.Msg = &MsgUnjail{}
	_ sdk.Msg = &MsgSend{}
)

//----------------------------------------------------------------------------------------------------------------------
// MsgStake - struct for staking transactions
type MsgStake struct {
	Address    sdk.ValAddress      `json:"validator_address" yaml:"validator_address"`
	PubKey     crypto.PubKey       `json:"pubkey" yaml:"pubkey"`
	Chains     map[string]struct{} `json:"chains" yaml:"chains"`
	Value      sdk.Int             `json:"value" yaml:"value"`
	ServiceURL string              `json:"serviceurl" yaml:"serviceurl"`
}

// Return address(es) that must sign over msg.GetSignBytes()
func (msg MsgStake) GetSigners() []sdk.AccAddress {
	addrs := []sdk.AccAddress{sdk.AccAddress(msg.Address)}
	return addrs
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgStake) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgStake) ValidateBasic() sdk.Error {
	if msg.Address.Empty() {
		return ErrNilValidatorAddr(DefaultCodespace)
	}
	if msg.Value.LTE(sdk.ZeroInt()) {
		return ErrBadDelegationAmount(DefaultCodespace)
	}
	if len(msg.Chains) == 0 {
		return ErrNoChains(DefaultCodespace)
	}
	for chain := range msg.Chains {
		if err := types.HashVerification(chain); err != nil {
			return err
		}
	}
	if len(msg.ServiceURL) == 0 {
		return ErrNoServiceURL(DefaultCodespace)
	}
	return nil
}

//nolint
func (msg MsgStake) Route() string { return RouterKey }
func (msg MsgStake) Type() string  { return "stake_validator" }

//----------------------------------------------------------------------------------------------------------------------
// MsgBeginUnstake - struct for unstaking transaciton
type MsgBeginUnstake struct {
	Address sdk.ValAddress `json:"validator_address" yaml:"validator_address"`
}

func (msg MsgBeginUnstake) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Address)}
}

func (msg MsgBeginUnstake) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgBeginUnstake) ValidateBasic() sdk.Error {
	if msg.Address.Empty() {
		return ErrNilValidatorAddr(DefaultCodespace)
	}
	return nil
}

//nolint
func (msg MsgBeginUnstake) Route() string { return RouterKey }
func (msg MsgBeginUnstake) Type() string  { return "begin_unstaking_validator" }

//----------------------------------------------------------------------------------------------------------------------
// MsgUnjail - struct for unjailing jailed validator
type MsgUnjail struct {
	ValidatorAddr sdk.ValAddress `json:"address" yaml:"address"` // address of the validator operator
}

//nolint
func (msg MsgUnjail) Route() string { return RouterKey }
func (msg MsgUnjail) Type() string  { return "unjail" }
func (msg MsgUnjail) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.ValidatorAddr)}
}

func (msg MsgUnjail) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgUnjail) ValidateBasic() sdk.Error {
	if msg.ValidatorAddr.Empty() {
		return ErrBadValidatorAddr(DefaultCodespace)
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------------
// MsgSend structure for sending coins
type MsgSend struct {
	FromAddress sdk.ValAddress
	ToAddress   sdk.ValAddress
	Amount      sdk.Int
}

//nolint
func (msg MsgSend) Route() string { return RouterKey }
func (msg MsgSend) Type() string  { return "send" }
func (msg MsgSend) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.FromAddress)}
}

func (msg MsgSend) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgSend) ValidateBasic() sdk.Error {
	if msg.FromAddress.Empty() {
		return ErrBadValidatorAddr(DefaultCodespace)
	}
	if msg.ToAddress.Empty() {
		return ErrBadValidatorAddr(DefaultCodespace)
	}
	if msg.Amount.LTE(sdk.ZeroInt()) {
		return ErrBadSendAmount(DefaultCodespace)
	}
	return nil
}
