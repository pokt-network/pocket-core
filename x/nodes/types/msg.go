package types

import (
	"fmt"
	"github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
)

// ensure Msg interface compliance at compile time
var (
	_ sdk.Msg = &MsgNodeStake{}
	_ sdk.Msg = &MsgBeginUnstake{}
	_ sdk.Msg = &MsgUnjail{}
	_ sdk.Msg = &MsgSend{}
	_ sdk.Msg = &MsgStake{}
)

const (
	MsgStakeName   = "stake_validator"
	MsgUnstakeName = "begin_unstake_validator"
	MsgUnjailName  = "unjail_validator"
	MsgSendName    = "send"
)

//----------------------------------------------------------------------------------------------------------------------

// GetSigners retrun address(es) that must sign over msg.GetSignBytes()
func (msg MsgNodeStake) GetSigner() sdk.Address {
	pubkey, err := crypto.NewPublicKey(msg.Publickey)
	if err != nil {
		return sdk.Address{}
	}
	return sdk.Address(pubkey.Address())
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgNodeStake) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic quick validity check, stateless
func (msg MsgNodeStake) ValidateBasic() sdk.Error {
	if msg.Publickey == "" {
		return ErrNilValidatorAddr(DefaultCodespace)
	}
	if msg.Value.LTE(sdk.ZeroInt()) {
		return ErrBadDelegationAmount(DefaultCodespace)
	}
	if len(msg.Chains) == 0 {
		return ErrNoChains(DefaultCodespace)
	}
	for _, chain := range msg.Chains {
		err := ValidateNetworkIdentifier(chain)
		if err != nil {
			return err
		}
	}
	if err := ValidateServiceURL(msg.ServiceUrl); err != nil {
		return err
	}
	return nil
}

// Route provides router key for msg
func (msg MsgNodeStake) Route() string { return RouterKey }

// Type provides msg name
func (msg MsgNodeStake) Type() string { return MsgStakeName }

// GetFee get fee for msg
func (msg MsgNodeStake) GetFee() sdk.Int {
	return sdk.NewInt(NodeFeeMap[msg.Type()])
}

//----------------------------------------------------------------------------------------------------------------------

// GetSigners return address(es) that must sign over msg.GetSignBytes()
func (msg MsgBeginUnstake) GetSigner() sdk.Address {
	return msg.Address
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgBeginUnstake) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic quick validity check, stateless
func (msg MsgBeginUnstake) ValidateBasic() sdk.Error {
	if msg.Address.Empty() {
		return ErrNilValidatorAddr(DefaultCodespace)
	}
	return nil
}

// Route provides router key for msg
func (msg MsgBeginUnstake) Route() string { return RouterKey }

// Type provides msg name
func (msg MsgBeginUnstake) Type() string { return MsgUnstakeName }

// GetFee get fee for msg
func (msg MsgBeginUnstake) GetFee() sdk.Int {
	return sdk.NewInt(NodeFeeMap[msg.Type()])
}

//----------------------------------------------------------------------------------------------------------------------

// GetSigners return address(es) that must sign over msg.GetSignBytes()
func (msg MsgUnjail) GetSigner() sdk.Address {
	return msg.ValidatorAddr
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgUnjail) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic quick validity check, stateless
func (msg MsgUnjail) ValidateBasic() sdk.Error {
	if msg.ValidatorAddr.Empty() {
		return ErrNoValidatorFound(DefaultCodespace)
	}
	return nil
}

// Route provides router key for msg
func (msg MsgUnjail) Route() string { return RouterKey }

// Type provides msg name
func (msg MsgUnjail) Type() string { return MsgUnjailName }

// GetFee get fee for msg
func (msg MsgUnjail) GetFee() sdk.Int {
	return sdk.NewInt(NodeFeeMap[msg.Type()])
}

//----------------------------------------------------------------------------------------------------------------------

// GetSigners return address(es) that must sign over msg.GetSignBytes()
func (msg MsgSend) GetSigner() sdk.Address {
	return msg.FromAddress
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgSend) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic quick validity check, stateless
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

// Route provides router key for msg
func (msg MsgSend) Route() string { return RouterKey }

// Type provides msg name
func (msg MsgSend) Type() string { return MsgSendName }

// GetFee get fee for msg
func (msg MsgSend) GetFee() sdk.Int {
	return sdk.NewInt(NodeFeeMap[msg.Type()])
}

//----------------------------------------------------------------------------------------------------------------------

// MsgStake - struct for staking transactions
type MsgStake struct {
	PublicKey  crypto.PublicKey `json:"public_key" yaml:"public_key"`
	Chains     []string         `json:"chains" yaml:"chains"`
	Value      sdk.Int          `json:"value" yaml:"value"`
	ServiceURL string           `json:"service_url" yaml:"service_url"`
} // GetSigners retrun address(es) that must sign over msg.GetSignBytes()

func (msg MsgStake) GetSigner() sdk.Address {
	return sdk.Address(msg.PublicKey.Address())
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgStake) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic quick validity check, stateless
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
		err := ValidateNetworkIdentifier(chain)
		if err != nil {
			return err
		}
	}
	if err := ValidateServiceURL(msg.ServiceURL); err != nil {
		return err
	}
	return nil
}

// Route provides router key for msg
func (msg MsgStake) Route() string { return RouterKey }

// Type provides msg name
func (msg MsgStake) Type() string { return MsgStakeName }

// GetFee get fee for msg
func (msg MsgStake) GetFee() sdk.Int {
	return sdk.NewInt(NodeFeeMap[msg.Type()])
}
func (msg MsgStake) Reset() {
	panic("amino only msg")
}

func (msg MsgStake) String() string {
	return fmt.Sprintf("Public Key: %s\nChains: %s\nValue: %s\n", msg.PublicKey.RawString(), msg.Chains, msg.Value.String())
}

func (msg MsgStake) ProtoMessage() {
	panic("amino only msg")
}

// GetFee get fee for msg
func (msg MsgStake) ToProto() MsgNodeStake {
	return MsgNodeStake{
		Publickey:  msg.PublicKey.RawString(),
		Chains:     msg.Chains,
		Value:      msg.Value,
		ServiceUrl: msg.ServiceURL,
	}
}

//----------------------------------------------------------------------------------------------------------------------
