package types

import (
	"fmt"
	"github.com/pokt-network/pocket-core/codec"
	"github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
	"time"
)

var _ codec.ProtoMarshaler = &LegacyValidator{}

type LegacyValidator struct {
	Address                 sdk.Address      `json:"address" yaml:"address"`               // address of the validator; hex encoded in JSON
	PublicKey               crypto.PublicKey `json:"public_key" yaml:"public_key"`         // the consensus public key of the validator; hex encoded in JSON
	Jailed                  bool             `json:"jailed" yaml:"jailed"`                 // has the validator been jailed from staked status?
	Status                  sdk.StakeStatus  `json:"status" yaml:"status"`                 // validator status (staked/unstaking/unstaked)
	Chains                  []string         `json:"chains" yaml:"chains"`                 // validator non native blockchains
	ServiceURL              string           `json:"service_url" yaml:"service_url"`       // url where the pocket service api is hosted
	StakedTokens            sdk.BigInt       `json:"tokens" yaml:"tokens"`                 // tokens staked in the network
	UnstakingCompletionTime time.Time        `json:"unstaking_time" yaml:"unstaking_time"` // if unstaking, min time for the validator to complete unstaking
}

func (v *LegacyValidator) Marshal() ([]byte, error) {
	a := v.ToProto()
	return a.Marshal()
}

func (v *LegacyValidator) MarshalTo(data []byte) (n int, err error) {
	a := v.ToProto()
	return a.MarshalTo(data)
}

func (v *LegacyValidator) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	a := v.ToProto()
	return a.MarshalToSizedBuffer(dAtA)
}

func (v *LegacyValidator) Size() int {
	a := v.ToProto()
	return a.Size()
}

func (v *LegacyValidator) Unmarshal(data []byte) error {
	var vp LegacyProtoValidator
	err := vp.Unmarshal(data)
	if err != nil {
		return err
	}
	*v, err = vp.FromProto()
	return err
}

func (v *LegacyValidator) IsStaked() bool {
	val := v.ToValidator()
	return val.IsStaked()
}

func (v *LegacyValidator) IsUnstaked() bool {
	val := v.ToValidator()
	return val.IsUnstaked()
}

func (v *LegacyValidator) IsUnstaking() bool {
	val := v.ToValidator()
	return val.IsUnstaking()
}

func (v *LegacyValidator) IsJailed() bool {
	val := v.ToValidator()
	return val.IsJailed()
}

func (v *LegacyValidator) GetStatus() sdk.StakeStatus {
	val := v.ToValidator()
	return val.GetStatus()
}

func (v *LegacyValidator) GetAddress() sdk.Address {
	val := v.ToValidator()
	return val.GetAddress()
}

func (v *LegacyValidator) GetPublicKey() crypto.PublicKey {
	val := v.ToValidator()
	return val.GetPublicKey()
}

func (v *LegacyValidator) GetTokens() sdk.BigInt {
	val := v.ToValidator()
	return val.GetTokens()
}

func (v *LegacyValidator) GetConsensusPower() int64 {
	val := v.ToValidator()
	return val.GetConsensusPower()
}

func (v *LegacyValidator) GetChains() []string {
	val := v.ToValidator()
	return val.GetChains()
}

func (v *LegacyValidator) Reset() {
	*v = LegacyValidator{}
}

func (v LegacyValidator) String() string {
	return fmt.Sprintf("Address:\t\t%s\nPublic Key:\t\t%s\nJailed:\t\t\t%v\nStatus:\t\t\t%s\nTokens:\t\t\t%s\n"+
		"ServiceUrl:\t\t%s\nChains:\t\t\t%v\nUnstaking Completion Time:\t\t%v\n"+
		"\n----\n",
		v.Address, v.PublicKey.RawString(), v.Jailed, v.Status, v.StakedTokens, v.ServiceURL, v.Chains, v.UnstakingCompletionTime,
	)
}

func (v LegacyValidator) ProtoMessage() {
	val := v.ToValidator()
	val.ProtoMessage()
}

func (v LegacyValidator) ToValidator() Validator {
	return Validator{
		Address:                 v.Address,
		PublicKey:               v.PublicKey,
		Jailed:                  v.Jailed,
		Status:                  v.Status,
		Chains:                  v.Chains,
		ServiceURL:              v.ServiceURL,
		StakedTokens:            v.StakedTokens,
		UnstakingCompletionTime: v.UnstakingCompletionTime,
		OutputAddress:           nil,
	}
}

func (v Validator) ToLegacy() LegacyValidator {
	return LegacyValidator{
		Address:                 v.Address,
		PublicKey:               v.PublicKey,
		Jailed:                  v.Jailed,
		Status:                  v.Status,
		Chains:                  v.Chains,
		ServiceURL:              v.ServiceURL,
		StakedTokens:            v.StakedTokens,
		UnstakingCompletionTime: v.UnstakingCompletionTime,
	}
}

// FromProto converts the Protobuf structure to Validator
func (v LegacyProtoValidator) FromProto() (LegacyValidator, error) {
	pubkey, err := crypto.NewPublicKeyBz(v.PublicKey)
	if err != nil {
		return LegacyValidator{}, err
	}
	return LegacyValidator{
		Address:                 v.Address,
		PublicKey:               pubkey,
		Jailed:                  v.Jailed,
		Status:                  sdk.StakeStatus(v.Status),
		ServiceURL:              v.ServiceURL,
		Chains:                  v.Chains,
		StakedTokens:            v.StakedTokens,
		UnstakingCompletionTime: v.UnstakingCompletionTime,
	}, nil
}

// ToProto converts the validator to Protobuf compatible structure
func (v LegacyValidator) ToProto() LegacyProtoValidator {
	return LegacyProtoValidator{
		Address:                 v.Address,
		PublicKey:               v.PublicKey.RawBytes(),
		Jailed:                  v.Jailed,
		Status:                  int32(v.Status),
		Chains:                  v.Chains,
		ServiceURL:              v.ServiceURL,
		StakedTokens:            v.StakedTokens,
		UnstakingCompletionTime: v.UnstakingCompletionTime,
	}
}
