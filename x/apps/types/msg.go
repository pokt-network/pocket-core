package types

import (
	"github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
)

// ensure Msg interface compliance at compile time
var (
	_ sdk.Msg = &MsgAppStake{}
	_ sdk.Msg = &MsgBeginAppUnstake{}
	_ sdk.Msg = &MsgAppUnjail{}
)

const (
	MsgAppStakeName   = "app_stake"
	MsgAppUnstakeName = "app_begin_unstake"
	MsgAppUnjailName  = "app_unjail"
)

//----------------------------------------------------------------------------------------------------------------------

// MsgAppStake - struct for staking transactions
type MsgAppStake struct {
	PubKey crypto.PublicKey `json:"pubkey" yaml:"pubkey"`
	Chains []string         `json:"chains" yaml:"chains"`
	Value  sdk.Int          `json:"value" yaml:"value"`
}

// GetSigners return address(es) that must sign over msg.GetSignBytes()
func (msg MsgAppStake) GetSigner() sdk.Address {
	return sdk.Address(msg.PubKey.Address())
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgAppStake) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic quick validity check for staking an application
func (msg MsgAppStake) ValidateBasic() sdk.Error {
	if msg.PubKey == nil || msg.PubKey.RawString() == "" {
		return ErrNilApplicationAddr(DefaultCodespace)
	}
	if msg.Value.LTE(sdk.ZeroInt()) {
		return ErrBadStakeAmount(DefaultCodespace)
	}
	if len(msg.Chains) == 0 {
		return ErrNoChains(DefaultCodespace)
	}
	for _, chain := range msg.Chains {
		if err := ValidateNetworkIdentifier(chain); err != nil {
			return err
		}
	}
	return nil
}

// Route provides router key for msg
func (msg MsgAppStake) Route() string { return RouterKey }

// Type provides msg name
func (msg MsgAppStake) Type() string { return MsgAppStakeName }

// GetFee get fee for msg
func (msg MsgAppStake) GetFee() sdk.Int {
	return sdk.NewInt(AppFeeMap[msg.Type()])
}

//----------------------------------------------------------------------------------------------------------------------

// MsgBeginAppUnstake - struct for unstaking transaciton
type MsgBeginAppUnstake struct {
	Address sdk.Address `json:"application_address" yaml:"application_address"`
}

// GetSigners address(es) that must sign over msg.GetSignBytes()
func (msg MsgBeginAppUnstake) GetSigner() sdk.Address {
	return msg.Address
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgBeginAppUnstake) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic quick validity check for staking an application
func (msg MsgBeginAppUnstake) ValidateBasic() sdk.Error {
	if msg.Address.Empty() {
		return ErrNilApplicationAddr(DefaultCodespace)
	}
	return nil
}

// Route provides router key for msg
func (msg MsgBeginAppUnstake) Route() string { return RouterKey }

// Type provides msg name
func (msg MsgBeginAppUnstake) Type() string { return MsgAppUnstakeName }

// GetFee get fee for msg
func (msg MsgBeginAppUnstake) GetFee() sdk.Int {
	return sdk.NewInt(AppFeeMap[msg.Type()])
}

//----------------------------------------------------------------------------------------------------------------------

// MsgAppUnjail - struct for unjailing jailed application
type MsgAppUnjail struct {
	AppAddr sdk.Address `json:"address" yaml:"address"` // address of the application operator
}

// Route provides router key for msg
func (msg MsgAppUnjail) Route() string { return RouterKey }

// Type provides msg name
func (msg MsgAppUnjail) Type() string { return MsgAppUnjailName }

// GetFee get fee for msg
func (msg MsgAppUnjail) GetFee() sdk.Int {
	return sdk.NewInt(AppFeeMap[msg.Type()])
}

// GetSigners return address(es) that must sign over msg.GetSignBytes()
func (msg MsgAppUnjail) GetSigner() sdk.Address {
	return msg.AppAddr
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgAppUnjail) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic quick validity check for staking an application
func (msg MsgAppUnjail) ValidateBasic() sdk.Error {
	if msg.AppAddr.Empty() {
		return ErrBadApplicationAddr(DefaultCodespace)
	}
	return nil
}
