package types

import (
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/pokt-network/posmint/types"
)

// ensure Msg interface compliance at compile time
var (
	_ sdk.Msg = &MsgAppStake{}
	_ sdk.Msg = &MsgBeginAppUnstake{}
	_ sdk.Msg = &MsgAppUnjail{}
)

//----------------------------------------------------------------------------------------------------------------------
// MsgAppStake - struct for staking transactions
type MsgAppStake struct {
	Address sdk.ValAddress      `json:"application_address" yaml:"application_address"`
	PubKey  crypto.PubKey       `json:"pubkey" yaml:"pubkey"`
	Chains  map[string]struct{} `json:"chains" yaml:"chains"`
	Value   sdk.Int             `json:"value" yaml:"value"`
}

// Return address(es) that must sign over msg.GetSignBytes()
func (msg MsgAppStake) GetSigners() []sdk.AccAddress {
	addrs := []sdk.AccAddress{sdk.AccAddress(msg.Address)}
	return addrs
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgAppStake) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgAppStake) ValidateBasic() sdk.Error {
	if msg.Address.Empty() {
		return ErrNilApplicationAddr(DefaultCodespace)
	}
	if msg.Value.LTE(sdk.ZeroInt()) {
		return ErrBadStakeAmount(DefaultCodespace)
	}
	if len(msg.Chains) == 0 {
		return ErrNoChains(DefaultCodespace)
	}
	for chain := range msg.Chains {
		if err := types.HashVerification(chain); err != nil {
			return err
		}
	}
	return nil
}

//nolint
func (msg MsgAppStake) Route() string { return RouterKey }
func (msg MsgAppStake) Type() string  { return "stake_application" }

//----------------------------------------------------------------------------------------------------------------------
// MsgBeginAppUnstake - struct for unstaking transaciton
type MsgBeginAppUnstake struct {
	Address sdk.ValAddress `json:"application_address" yaml:"application_address"`
}

func (msg MsgBeginAppUnstake) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Address)}
}

func (msg MsgBeginAppUnstake) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgBeginAppUnstake) ValidateBasic() sdk.Error {
	if msg.Address.Empty() {
		return ErrNilApplicationAddr(DefaultCodespace)
	}
	return nil
}

//nolint
func (msg MsgBeginAppUnstake) Route() string { return RouterKey }
func (msg MsgBeginAppUnstake) Type() string  { return "begin_unstaking_application" }

//----------------------------------------------------------------------------------------------------------------------
// MsgAppUnjail - struct for unjailing jailed application
type MsgAppUnjail struct {
	AppAddr sdk.ValAddress `json:"address" yaml:"address"` // address of the application operator
}

//nolint
func (msg MsgAppUnjail) Route() string { return RouterKey }
func (msg MsgAppUnjail) Type() string  { return "unjail" }
func (msg MsgAppUnjail) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.AppAddr)}
}

func (msg MsgAppUnjail) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgAppUnjail) ValidateBasic() sdk.Error {
	if msg.AppAddr.Empty() {
		return ErrBadApplicationAddr(DefaultCodespace)
	}
	return nil
}
