package types

import (
	"encoding/json"
	"fmt"
	"github.com/pokt-network/posmint/codec"
	"github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
	"strings"
	"time"
)

// Validators is a collection of Validator
type Validators []Validator

func (v Validators) String() (out string) {
	for _, val := range v {
		out += val.String() + "\n"
	}
	return strings.TrimSpace(out)
}

func (v Validators) JSON() (out []byte, err error) {
	// each element should be a JSON
	return json.Marshal(v)
}

// MUST return the amino encoded version of this validator
func MustMarshalValidator(cdc *codec.Codec, validator Validator) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(validator)
}

// MUST decode the validator from the bytes
func MustUnmarshalValidator(cdc *codec.Codec, valBytes []byte) Validator {
	validator, err := UnmarshalValidator(cdc, valBytes)
	if err != nil {
		panic(err)
	}
	return validator
}

// unmarshal the validator
func UnmarshalValidator(cdc *codec.Codec, valBytes []byte) (validator Validator, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(valBytes, &validator)
	return validator, err
}

// HashString returns a human readable string representation of a validator.
func (v Validator) String() string {
	return fmt.Sprintf(`Validator
  Address:           		  %s
  Validator Cons Pubkey:      %s
  Jailed:                     %v
  Status:                     %s
  Tokens:               	  %s
  ServiceURL:                 %s
  Chains:                     %v
  Unstaking Completion Time:  %v`,
		v.Address, v.PublicKey.RawString(), v.Jailed, v.Status, v.StakedTokens, v.ServiceURL, v.Chains, v.UnstakingCompletionTime,
	)
}

// this is a helper struct used for JSON de- and encoding only
type hexValidator struct {
	Address                 sdk.Address     `json:"address" yaml:"address"`       // the hex address of the validator
	PublicKey               string          `json:"public_key" yaml:"public_key"` // the hex consensus public key of the validator
	Jailed                  bool            `json:"jailed" yaml:"jailed"`         // has the validator been jailed from staked status?
	Status                  sdk.StakeStatus `json:"status" yaml:"status"`         // validator status (staked/unstaking/unstaked)
	StakedTokens            sdk.Int         `json:"tokens" yaml:"tokens"`         // how many staked tokens
	ServiceURL              string          `json:"service_url" yaml:"service_uRL"`
	Chains                  []string        `json:"chains" yaml:"chains"`
	UnstakingCompletionTime time.Time       `json:"unstaking_time" yaml:"unstaking_time"` // if unstaking, min time for the validator to complete unstaking
}

// MarshalJSON marshals the validator to JSON using Hex
func (v Validator) MarshalJSON() ([]byte, error) {
	return codec.Cdc.MarshalJSON(hexValidator{
		Address:                 v.Address,
		PublicKey:               v.PublicKey.RawString(),
		Jailed:                  v.Jailed,
		Status:                  v.Status,
		ServiceURL:              v.ServiceURL,
		Chains:                  v.Chains,
		StakedTokens:            v.StakedTokens,
		UnstakingCompletionTime: v.UnstakingCompletionTime,
	})
}

// UnmarshalJSON unmarshals the validator from JSON using Hex
func (v *Validator) UnmarshalJSON(data []byte) error {
	bv := &hexValidator{}
	if err := codec.Cdc.UnmarshalJSON(data, bv); err != nil {
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
	}
	return nil
}
