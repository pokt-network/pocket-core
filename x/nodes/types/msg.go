package types

import (
	"github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
)

// ensure Msg interface compliance at compile time
var (
	_ sdk.Msg = &MsgStake{}
	_ sdk.Msg = &MsgBeginUnstake{}
	_ sdk.Msg = &MsgUnjail{}
	_ sdk.Msg = &MsgSend{}
)

const (
	MsgStakeName   = "stake_validator"
	MsgUnstakeName = "begin_unstake_validator"
	MsgUnjailName  = "unjail_validator"
	MsgSendName    = "send"
)

//----------------------------------------------------------------------------------------------------------------------
// MsgStake - struct for staking transactions
type MsgStake struct {
	PublicKey  crypto.PublicKey `json:"public_key" yaml:"public_key"`
	Chains     []string         `json:"chains" yaml:"chains"`
	Value      sdk.Int          `json:"value" yaml:"value"`
	ServiceURL string           `json:"service_url" yaml:"service_url"`
}

// Return address(es) that must sign over msg.GetSignBytes()
func (msg MsgStake) GetSigners() []sdk.Address {
	addrs := []sdk.Address{sdk.Address(msg.PublicKey.Address())}
	return addrs
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgStake) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// Quick validity check, stateless
func (msg MsgStake) ValidateBasic() sdk.Error {
	if msg.PublicKey == nil || msg.PublicKey.RawString() == "" {
		return ErrNilValidatorAddr(DefaultCodespace)
	}
	if msg.Value.LTE(sdk.ZeroInt()) {
		return ErrBadDelegationAmount(DefaultCodespace)
	}
	if len(msg.Chains) == 0 {
		return ErrNoChains(DefaultCodespace)
	}
	for _, chain := range msg.Chains {
		if len(chain) == 0 {
			return ErrNoChains(DefaultCodespace)
		}
	}
	if len(msg.ServiceURL) == 0 {
		return ErrNoServiceURL(DefaultCodespace)
	}
	return nil
}

func (msg MsgStake) Route() string { return RouterKey }
func (msg MsgStake) Type() string  { return MsgStakeName }

//----------------------------------------------------------------------------------------------------------------------
// MsgBeginUnstake - struct for unstaking transaciton
type MsgBeginUnstake struct {
	Address sdk.Address `json:"validator_address" yaml:"validator_address"`
}

// Return address(es) that must sign over msg.GetSignBytes()
func (msg MsgBeginUnstake) GetSigners() []sdk.Address {
	return []sdk.Address{msg.Address}
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgBeginUnstake) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// Quick validity check, stateless
func (msg MsgBeginUnstake) ValidateBasic() sdk.Error {
	if msg.Address.Empty() {
		return ErrNilValidatorAddr(DefaultCodespace)
	}
	return nil
}

func (msg MsgBeginUnstake) Route() string { return RouterKey }
func (msg MsgBeginUnstake) Type() string  { return MsgUnstakeName }

//----------------------------------------------------------------------------------------------------------------------
// MsgUnjail - struct for unjailing jailed validator
type MsgUnjail struct {
	ValidatorAddr sdk.Address `json:"address" yaml:"address"` // address of the validator operator
}

// Return address(es) that must sign over msg.GetSignBytes()
func (msg MsgUnjail) GetSigners() []sdk.Address {
	return []sdk.Address{msg.ValidatorAddr}
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgUnjail) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// Quick validity check, stateless
func (msg MsgUnjail) ValidateBasic() sdk.Error {
	if msg.ValidatorAddr.Empty() {
		return ErrNoValidatorFound(DefaultCodespace)
	}
	return nil
}

func (msg MsgUnjail) Route() string { return RouterKey }
func (msg MsgUnjail) Type() string  { return MsgUnjailName }

//----------------------------------------------------------------------------------------------------------------------
// MsgSend structure for sending coins
type MsgSend struct {
	FromAddress sdk.Address
	ToAddress   sdk.Address
	Amount      sdk.Int
}

// Return address(es) that must sign over msg.GetSignBytes()
func (msg MsgSend) GetSigners() []sdk.Address {
	return []sdk.Address{msg.FromAddress}
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgSend) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// Quick validity check, stateless
func (msg MsgSend) ValidateBasic() sdk.Error {
	if msg.FromAddress.Empty() {
		return ErrNoValidatorFound(DefaultCodespace)
	}
	if msg.ToAddress.Empty() {
		return ErrNoValidatorFound(DefaultCodespace)
	}
	if msg.Amount.LTE(sdk.ZeroInt()) {
		return ErrBadSendAmount(DefaultCodespace)
	}
	return nil
}

func (msg MsgSend) Route() string { return RouterKey }
func (msg MsgSend) Type() string  { return MsgSendName }
