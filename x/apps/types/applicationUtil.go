package types

import (
	"encoding/json"
	"fmt"
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"strings"
	"time"
)

// Applications is a collection of AppPubKey
type Applications []Application

func (v Applications) String() (out string) {
	for _, val := range v {
		out += val.String() + "\n"
	}
	return strings.TrimSpace(out)
}

func (v Applications) JSON() (out []byte, err error) {
	var result []string
	for _, val := range v {
		r := val.String()
		result = append(result, r)
	}
	return json.Marshal(result)
}

// MUST return the amino encoded version of this application
func MustMarshalApplication(cdc *codec.Codec, application Application) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(application)
}

// MUST decode the app from the bytes
func MustUnmarshalApplication(cdc *codec.Codec, valBytes []byte) Application {
	application, err := UnmarshalApplication(cdc, valBytes)
	if err != nil { // todo test panic!!
		panic(err)
	}
	return application
}

// unmarshal the application
func UnmarshalApplication(cdc *codec.Codec, appBytes []byte) (application Application, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(appBytes, &application)
	return application, err
}

// HashString returns a human readable string representation of a application.
func (a Application) String() string {
	return fmt.Sprintf(`AppPubKey
  Address:           		  %s
  AppPubKey Cons Pubkey: 	  %s
  Jailed:                     %v
  Chains:                     %v
  MaxRelays:                  %d
  Status:                     %s
  Tokens:               	  %s
  Unstakeing Completion Time: %v`,
		a.Address, sdk.HexConsPub(a.ConsPubKey), a.Jailed, a.Chains, a.MaxRelays, a.Status, a.StakedTokens, a.UnstakingCompletionTime,
	)
}

// this is a helper struct used for JSON de- and encoding only
type hexApplication struct {
	Address                 sdk.ValAddress `json:"operator_address" yaml:"operator_address"` // the hex address of the application
	ConsPubKey              string         `json:"cons_pubkey" yaml:"cons_pubkey"`           // the hex consensus public key of the application
	Jailed                  bool           `json:"jailed" yaml:"jailed"`                     // has the application been jailed from staked status?
	Chains                  []string       `json:"chains" yaml:"chains"`
	MaxRelays               sdk.Int
	Status                  sdk.BondStatus `json:"status" yaml:"status"`                 // application status (bonded/unbonding/unbonded)
	StakedTokens            sdk.Int        `json:"stakedTokens" yaml:"stakedTokens"`     // how many staked tokens
	UnstakingCompletionTime time.Time      `json:"unstaking_time" yaml:"unstaking_time"` // if unstaking, min time for the application to complete unstaking
}

// MarshalJSON marshals the application to JSON using Hex
func (a Application) MarshalJSON() ([]byte, error) {
	return codec.Cdc.MarshalJSON(hexApplication{
		Address:                 a.Address,
		ConsPubKey:              sdk.HexConsPub(a.ConsPubKey),
		Jailed:                  a.Jailed,
		Status:                  a.Status,
		Chains:                  a.Chains,
		MaxRelays:               a.MaxRelays,
		StakedTokens:            a.StakedTokens,
		UnstakingCompletionTime: a.UnstakingCompletionTime,
	})
}

// UnmarshalJSON unmarshals the application from JSON using Hex
func (a *Application) UnmarshalJSON(data []byte) error {
	bv := &hexApplication{}
	if err := codec.Cdc.UnmarshalJSON(data, bv); err != nil {
		return err
	}
	consPubKey, err := sdk.GetConsPubKeyHex(bv.ConsPubKey)
	if err != nil {
		return err
	}
	*a = Application{
		Address:                 bv.Address,
		ConsPubKey:              consPubKey,
		Chains:                  bv.Chains,
		MaxRelays:               bv.MaxRelays,
		Jailed:                  bv.Jailed,
		StakedTokens:            bv.StakedTokens,
		Status:                  bv.Status,
		UnstakingCompletionTime: bv.UnstakingCompletionTime,
	}
	return nil
}
