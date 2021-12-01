package types

import (
	"fmt"
	"github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
)

// unjail legacy

func (m LegacyMsgUnjail) Route() string {
	return RouterKey
}

func (m LegacyMsgUnjail) Type() string {
	return MsgUnjailName
}

func (m LegacyMsgUnjail) ValidateBasic() sdk.Error {
	if m.ValidatorAddr.Empty() {
		return ErrNoValidatorFound(DefaultCodespace)
	}
	return nil
}

func (m LegacyMsgUnjail) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m LegacyMsgUnjail) GetSigners() []sdk.Address {
	return []sdk.Address{m.ValidatorAddr}
}

func (m LegacyMsgUnjail) GetRecipient() sdk.Address {
	return nil
}

func (m LegacyMsgUnjail) GetFee() sdk.BigInt {
	return sdk.NewInt(NodeFeeMap[m.Type()])
}
func (*LegacyMsgUnjail) XXX_MessageName() string {
	return "x.nodes.MsgUnjail"
}

// Unstake Legacy

func (m LegacyMsgBeginUnstake) Route() string {
	return RouterKey
}

func (m LegacyMsgBeginUnstake) Type() string {
	return MsgUnstakeName
}

func (m LegacyMsgBeginUnstake) ValidateBasic() sdk.Error {
	if m.Address.Empty() {
		return ErrNilValidatorAddr(DefaultCodespace)
	}
	return nil
}

func (m LegacyMsgBeginUnstake) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m LegacyMsgBeginUnstake) GetSigners() []sdk.Address {
	return []sdk.Address{m.Address}
}

func (m LegacyMsgBeginUnstake) GetRecipient() sdk.Address {
	return nil
}

func (m LegacyMsgBeginUnstake) GetFee() sdk.BigInt {
	return sdk.NewInt(NodeFeeMap[m.Type()])
}

func (*LegacyMsgBeginUnstake) XXX_MessageName() string {
	return "x.nodes.MsgBeginUnstake"
}

// stake legacy
type LegacyMsgStake struct {
	PublicKey  crypto.PublicKey `json:"public_key" yaml:"public_key"`
	Chains     []string         `json:"chains" yaml:"chains"`
	Value      sdk.BigInt       `json:"value" yaml:"value"`
	ServiceUrl string           `json:"service_url" yaml:"service_url"`
}

func (m LegacyMsgStake) Route() string {
	return RouterKey
}

func (m LegacyMsgStake) Type() string {
	return MsgStakeName
}

func (m LegacyMsgStake) ValidateBasic() sdk.Error {
	if m.PublicKey == nil || m.PublicKey.RawString() == "" {
		return ErrNilValidatorAddr(DefaultCodespace)
	}
	if m.Value.LTE(sdk.ZeroInt()) {
		return ErrBadDelegationAmount(DefaultCodespace)
	}
	if len(m.Chains) == 0 {
		return ErrNoChains(DefaultCodespace)
	}
	for _, chain := range m.Chains {
		err := ValidateNetworkIdentifier(chain)
		if err != nil {
			return err
		}
	}
	if err := ValidateServiceURL(m.ServiceUrl); err != nil {
		return err
	}
	return nil
}

func (m LegacyMsgStake) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m LegacyMsgStake) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(m.PublicKey.Address())}
}

func (m LegacyMsgStake) GetRecipient() sdk.Address {
	return nil
}

func (m LegacyMsgStake) GetFee() sdk.BigInt {
	return sdk.NewInt(NodeFeeMap[m.Type()])
}

func (m *LegacyMsgStake) Reset() {
	*m = LegacyMsgStake{}
}

func (m *LegacyMsgStake) String() string {
	return fmt.Sprintf("Public Key: %s\nChains: %s\nValue: %s\n", m.PublicKey.RawString(), m.Chains, m.Value.String())
}

func (m *LegacyMsgStake) ProtoMessage() {
	a := m.LegacyToProto()
	a.ProtoMessage()
}

// GetFee get fee for msg
func (m *LegacyMsgStake) LegacyToProto() LegacyMsgProtoStake {
	var pkbz []byte
	if m.PublicKey != nil {
		pkbz = m.PublicKey.RawBytes()
	}
	return LegacyMsgProtoStake{
		Publickey:  pkbz,
		Chains:     m.Chains,
		Value:      m.Value,
		ServiceUrl: m.ServiceUrl,
	}
}

func (msg *LegacyMsgStake) Marshal() ([]byte, error) {
	p := msg.LegacyToProto()
	return p.Marshal()
}

func (msg *LegacyMsgStake) MarshalTo(data []byte) (n int, err error) {
	p := msg.LegacyToProto()
	return p.MarshalTo(data)
}

func (msg *LegacyMsgStake) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	p := msg.LegacyToProto()
	return p.MarshalToSizedBuffer(dAtA)
}

func (msg *LegacyMsgStake) Size() int {
	p := msg.LegacyToProto()
	return p.Size()
}

func (msg *LegacyMsgStake) Unmarshal(data []byte) error {
	var m LegacyMsgProtoStake
	err := m.Unmarshal(data)
	if err != nil {
		return err
	}
	pk, err := crypto.NewPublicKeyBz(m.Publickey)
	if err != nil {
		return err
	}
	newMsg := LegacyMsgStake{
		PublicKey:  pk,
		Chains:     m.Chains,
		Value:      m.Value,
		ServiceUrl: m.ServiceUrl,
	}
	*msg = newMsg
	return nil
}

func (msg *LegacyMsgStake) XXX_MessageName() string {
	m := msg.LegacyToProto()
	return m.XXX_MessageName()
}

// GetFee get fee for msg
func (msg MsgStake) LegacyToProto() LegacyMsgProtoStake {
	var pkbz []byte
	if msg.PublicKey != nil {
		pkbz = msg.PublicKey.RawBytes()
	}
	return LegacyMsgProtoStake{
		Publickey:  pkbz,
		Chains:     msg.Chains,
		Value:      msg.Value,
		ServiceUrl: msg.ServiceUrl,
	}
}

func (*LegacyMsgProtoStake) XXX_MessageName() string {
	return "x.nodes.MsgProtoStake"
}
