package types

import (
	"fmt"

	"github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
)

// ensure Msg interface compliance at compile time
var (
	_ sdk.Msg = &MsgApplicationStake{}
	_ sdk.Msg = &MsgBeginAppUnstake{}
	_ sdk.Msg = &MsgAppUnjail{}
)

const (
	MsgAppStakeName   = "app_stake"
	MsgAppUnstakeName = "app_begin_unstake"
	MsgAppUnjailName  = "app_unjail"
)

//----------------------------------------------------------------------------------------------------------------------

// GetSigners return address(es) that must sign over msg.GetSignBytes()
func (msg MsgApplicationStake) GetSigner() sdk.Address {
	res, err := crypto.NewPublicKey(msg.PubKey)
	if err != nil {
		return nil
	}
	return sdk.Address(res.Address())
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgApplicationStake) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic quick validity check for staking an application
func (msg MsgApplicationStake) ValidateBasic() sdk.Error {
	if msg.PubKey == "" {
		return ErrNilApplicationAddr(DefaultCodespace)
	}
	pk, err := crypto.NewPublicKey(msg.PubKey)
	if err != nil {
		return sdk.ErrInvalidPubKey(err.Error())
	}
	if pk == nil || pk.RawString() == "" {
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
func (msg MsgApplicationStake) Route() string { return RouterKey }

// Type provides msg name
func (msg MsgApplicationStake) Type() string { return MsgAppStakeName }

// GetFee get fee for msg
func (msg MsgApplicationStake) GetFee() sdk.Int {
	return sdk.NewInt(AppFeeMap[msg.Type()])
}

//----------------------------------------------------------------------------------------------------------------------

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

// Legacy Apps Amino Msg below
// ---------------------------------------------------------------------------------------------------------------------
// MsgApplicationStake - struct for staking transactions
type MsgAppStake struct {
	PubKey crypto.PublicKey `json:"pubkey" yaml:"pubkey"`
	Chains []string         `json:"chains" yaml:"chains"`
	Value  sdk.Int          `json:"value" yaml:"value"`
}

func (msg MsgAppStake) Reset() {
	panic("amino only msg")
}

func (msg MsgAppStake) String() string {
	return fmt.Sprintf("Public Key: %s\nChains: %s\nValue: %s\n", msg.PubKey.RawString(), msg.Chains, msg.Value.String())
}

func (msg MsgAppStake) ProtoMessage() {
	panic("amino only msg")
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

func (msg MsgAppStake) ToProto() MsgApplicationStake {
	return MsgApplicationStake{
		PubKey: msg.PubKey.RawString(),
		Chains: msg.Chains,
		Value:  msg.Value,
	}
}
