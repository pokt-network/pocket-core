package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pokt-network/pocket-core/codec"
	"strings"
	"time"

	"github.com/pokt-network/pocket-core/crypto"

	sdk "github.com/pokt-network/pocket-core/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

type Validator struct {
	Address                 sdk.Address      `json:"address" yaml:"address"`                         // address of the validator; hex encoded in JSON
	PublicKey               crypto.PublicKey `json:"public_key" yaml:"public_key"`                   // the consensus public key of the validator; hex encoded in JSON
	Jailed                  bool             `json:"jailed" yaml:"jailed"`                           // has the validator been jailed from staked status?
	Status                  sdk.StakeStatus  `json:"status" yaml:"status"`                           // validator status (staked/unstaking/unstaked)
	Chains                  []string         `json:"chains" yaml:"chains"`                           // validator non native blockchains
	ServiceURL              string           `json:"service_url" yaml:"service_url"`                 // url where the pocket service api is hosted
	StakedTokens            sdk.BigInt       `json:"tokens" yaml:"tokens"`                           // tokens staked in the network
	UnstakingCompletionTime time.Time        `json:"unstaking_time" yaml:"unstaking_time"`           // if unstaking, min time for the validator to complete unstaking
	OutputAddress           sdk.Address      `json:"output_address,omitempty" yaml:"output_address"` // the custodial output address of the validator
}

// NewValidator - initialize a new validator
func NewValidator(addr sdk.Address, consPubKey crypto.PublicKey, chains []string, serviceURL string, tokensToStake sdk.BigInt, outputAddress sdk.Address) Validator {
	return Validator{
		Address:                 addr,
		PublicKey:               consPubKey,
		Jailed:                  false,
		Status:                  sdk.Staked,
		Chains:                  chains,
		ServiceURL:              serviceURL,
		StakedTokens:            tokensToStake,
		UnstakingCompletionTime: time.Time{},
		OutputAddress:           outputAddress,
	}
}

// ABCIValidatorUpdate returns an abci.ValidatorUpdate from a staking validator type
// with the full validator power
func (v Validator) ABCIValidatorUpdate() abci.ValidatorUpdate {
	return abci.ValidatorUpdate{
		PubKey: tmtypes.TM2PB.PubKey(v.PublicKey.PubKey()),
		Power:  v.ConsensusPower(),
	}
}

// ABCIValidatorUpdate returns an abci.ValidatorUpdate from a staking validator type
// with the full validator power
func (v Validator) ABCIValidatorZeroUpdate() abci.ValidatorUpdate {
	return abci.ValidatorUpdate{
		PubKey: tmtypes.TM2PB.PubKey(v.PublicKey.PubKey()),
		Power:  0,
	}
}

// get the consensus-engine power
// a reduction of 10^6 from validator tokens is applied
func (v Validator) ConsensusPower() int64 {
	if v.IsStaked() && !v.IsJailed() {
		return sdk.TokensToConsensusPower(v.StakedTokens)
	}
	return 0
}

// RemoveStakedTokens removes tokens from a validator
func (v Validator) RemoveStakedTokens(tokens sdk.BigInt) (Validator, error) {
	if tokens.IsNegative() {
		return Validator{}, fmt.Errorf("should not happen: trying to remove negative tokens: %s from valdiator %s", tokens.String(), v.Address)
	}
	if v.StakedTokens.LT(tokens) {
		return Validator{}, fmt.Errorf("should not happen: only have %v tokens, trying to remove %v", v.StakedTokens, tokens)
	}
	v.StakedTokens = v.StakedTokens.Sub(tokens)
	return v, nil
}

// AddStakedTokens tokens to staked field for a validator
func (v Validator) AddStakedTokens(tokens sdk.BigInt) (Validator, error) {
	if tokens.IsNegative() {
		return Validator{}, fmt.Errorf("should not happen: trying to add negative tokens: %s from valdiator %s", tokens.String(), v.Address)
	}
	v.StakedTokens = v.StakedTokens.Add(tokens)
	return v, nil
}

// compares the vital fields of two validator structures
func (v Validator) Equals(v2 Validator) bool {
	return v.PublicKey.Equals(v2.PublicKey) &&
		bytes.Equal(v.Address, v2.Address) &&
		v.Status.Equal(v2.Status) &&
		v.StakedTokens.Equal(v2.StakedTokens) &&
		v.OutputAddress.Equals(v2.OutputAddress)
}

// UpdateStatus updates the staking status
func (v Validator) UpdateStatus(newStatus sdk.StakeStatus) Validator {
	v.Status = newStatus
	return v
}

func (v Validator) HasChain(netID string) bool {
	for _, c := range v.Chains {
		if c == netID {
			return true
		}
	}
	return false
}

// return the TM validator address
func (v Validator) GetChains() []string            { return v.Chains }
func (v Validator) GetServiceURL() string          { return v.ServiceURL }
func (v Validator) IsStaked() bool                 { return v.GetStatus().Equal(sdk.Staked) }
func (v Validator) IsUnstaked() bool               { return v.GetStatus().Equal(sdk.Unstaked) }
func (v Validator) IsUnstaking() bool              { return v.GetStatus().Equal(sdk.Unstaking) }
func (v Validator) IsJailed() bool                 { return v.Jailed }
func (v Validator) GetStatus() sdk.StakeStatus     { return v.Status }
func (v Validator) GetAddress() sdk.Address        { return v.Address }
func (v Validator) GetPublicKey() crypto.PublicKey { return v.PublicKey }
func (v Validator) GetTokens() sdk.BigInt          { return v.StakedTokens }
func (v Validator) GetConsensusPower() int64       { return v.ConsensusPower() }
func (v *Validator) Reset()                        { *v = Validator{} }

func (v Validator) ProtoMessage() {
	p := v.ToProto()
	p.ProtoMessage()
}

func (v Validator) Marshal() ([]byte, error) {
	p := v.ToProto()
	return p.Marshal()
}

func (v Validator) MarshalTo(data []byte) (n int, err error) {
	p := v.ToProto()
	return p.MarshalTo(data)
}

func (v Validator) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	p := v.ToProto()
	return p.MarshalToSizedBuffer(dAtA)
}

func (v Validator) Size() int {
	p := v.ToProto()
	return p.Size()
}

func (v *Validator) Unmarshal(data []byte) error {
	var vp ProtoValidator
	err := vp.Unmarshal(data)
	if err != nil {
		return err
	}
	*v, err = vp.FromProto()
	return err
}

// String returns a human readable string representation of a validator.
func (v Validator) String() string {
	outputPubKeyString := ""
	if v.OutputAddress != nil {
		outputPubKeyString = v.OutputAddress.String()
	}
	return fmt.Sprintf("Address:\t\t%s\nPublic Key:\t\t%s\nJailed:\t\t\t%v\nStatus:\t\t\t%s\nTokens:\t\t\t%s\n"+
		"ServiceUrl:\t\t%s\nChains:\t\t\t%v\nUnstaking Completion Time:\t\t%v\nOutput Address:\t\t%s"+
		"\n----\n",
		v.Address, v.PublicKey.RawString(), v.Jailed, v.Status, v.StakedTokens, v.ServiceURL, v.Chains, v.UnstakingCompletionTime, outputPubKeyString,
	)
}

var _ codec.ProtoMarshaler = &Validator{}

// MarshalJSON marshals the validator to JSON using Hex
func (v Validator) MarshalJSON() ([]byte, error) {
	return json.Marshal(JSONValidator{
		Address:                 v.Address,
		PublicKey:               v.PublicKey.RawString(),
		Jailed:                  v.Jailed,
		Status:                  v.Status,
		ServiceURL:              v.ServiceURL,
		Chains:                  v.Chains,
		StakedTokens:            v.StakedTokens,
		UnstakingCompletionTime: v.UnstakingCompletionTime,
		OutputAddress:           v.OutputAddress,
	})
}

// UnmarshalJSON unmarshals the validator from JSON using Hex
func (v *Validator) UnmarshalJSON(data []byte) error {
	bv := &JSONValidator{}
	if err := json.Unmarshal(data, bv); err != nil {
		return err
	}
	publicKey, err := crypto.NewPublicKey(bv.PublicKey)
	if err != nil {
		return err
	}
	*v = Validator{
		Address:                 bv.Address,
		PublicKey:               publicKey,
		Jailed:                  bv.Jailed,
		Chains:                  bv.Chains,
		ServiceURL:              bv.ServiceURL,
		StakedTokens:            bv.StakedTokens,
		Status:                  bv.Status,
		UnstakingCompletionTime: bv.UnstakingCompletionTime,
		OutputAddress:           bv.OutputAddress,
	}
	return nil
}

// FromProto converts the Protobuf structure to Validator
func (v ProtoValidator) FromProto() (Validator, error) {
	pubkey, err := crypto.NewPublicKeyBz(v.PublicKey)
	if err != nil {
		return Validator{}, err
	}
	return Validator{
		Address:                 v.Address,
		PublicKey:               pubkey,
		Jailed:                  v.Jailed,
		Status:                  sdk.StakeStatus(v.Status),
		ServiceURL:              v.ServiceURL,
		Chains:                  v.Chains,
		StakedTokens:            v.StakedTokens,
		UnstakingCompletionTime: v.UnstakingCompletionTime,
		OutputAddress:           v.OutputAddress,
	}, nil
}

// ToProto converts the validator to Protobuf compatible structure
func (v Validator) ToProto() ProtoValidator {
	return ProtoValidator{
		Address:                 v.Address,
		PublicKey:               v.PublicKey.RawBytes(),
		Jailed:                  v.Jailed,
		Status:                  int32(v.Status),
		Chains:                  v.Chains,
		ServiceURL:              v.ServiceURL,
		StakedTokens:            v.StakedTokens,
		UnstakingCompletionTime: v.UnstakingCompletionTime,
		OutputAddress:           v.OutputAddress,
	}
}

type JSONValidator struct {
	Address                 sdk.Address     `json:"address" yaml:"address"`               // address of the validator; hex encoded in JSON
	PublicKey               string          `json:"public_key" yaml:"public_key"`         // the consensus public key of the validator; hex encoded in JSON
	Jailed                  bool            `json:"jailed" yaml:"jailed"`                 // has the validator been jailed from staked status?
	Status                  sdk.StakeStatus `json:"status" yaml:"status"`                 // validator status (staked/unstaking/unstaked)
	Chains                  []string        `json:"chains" yaml:"chains"`                 // validator non native blockchains
	ServiceURL              string          `json:"service_url" yaml:"service_url"`       // url where the pocket service api is hosted
	StakedTokens            sdk.BigInt      `json:"tokens" yaml:"tokens"`                 // tokens staked in the network
	UnstakingCompletionTime time.Time       `json:"unstaking_time" yaml:"unstaking_time"` // if unstaking, min time for the validator to complete unstaking
	OutputAddress           sdk.Address     `json:"output_address" yaml:"output_address"` // custodial output address of tokens
}

// Validators is a collection of Validator
type Validators []Validator

func (v Validators) String() (out string) {
	for _, val := range v {
		out += val.String() + "\n\n"
	}
	return strings.TrimSpace(out)
}

type ValidatorsPage struct {
	Result Validators `json:"result"`
	Total  int        `json:"total_pages"`
	Page   int        `json:"page"`
}

// String returns a human readable string representation of a validator page
func (vP ValidatorsPage) String() string {
	return fmt.Sprintf("Total:\t\t%d\nPage:\t\t%d\nResult:\t\t\n====\n%s\n====\n", vP.Total, vP.Page, vP.Result.String())
}
