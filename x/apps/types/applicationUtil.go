package types

import (
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

// MUST return the amino encoded version of this application
func MustMarshalApplication(cdc *codec.Codec, application Application) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(application)
}

// MUST decode the app from the bytes
func MustUnmarshalApplication(cdc *codec.Codec, valBytes []byte) Application {
	application, err := UnmarshalApplication(cdc, valBytes)
	if err != nil {
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
	bechConsPubKey, err := sdk.Bech32ifyConsPub(a.ConsPubKey)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf(`AppPubKey
  Address:           		  %s
  AppPubKey Cons Pubkey: 	  %s
  Jailed:                     %v
  Status:                     %s
  Tokens:               	  %s
  Unstakeing Completion Time: %v`,
		a.Address, bechConsPubKey, a.Jailed, a.Status, a.StakedTokens, a.UnstakingCompletionTime,
	)
}

// this is a helper struct used for JSON de- and encoding only
type bechApplication struct {
	Address                 sdk.ValAddress `json:"operator_address" yaml:"operator_address"` // the bech32 address of the application
	ConsPubKey              string         `json:"cons_pubkey" yaml:"cons_pubkey"`           // the bech32 consensus public key of the application
	Jailed                  bool           `json:"jailed" yaml:"jailed"`                     // has the application been jailed from staked status?
	Status                  sdk.BondStatus `json:"status" yaml:"status"`                     // application status (bonded/unbonding/unbonded)
	StakedTokens            sdk.Int        `json:"stakedTokens" yaml:"stakedTokens"`         // how many staked tokens
	UnstakingCompletionTime time.Time      `json:"unstaking_time" yaml:"unstaking_time"`     // if unstaking, min time for the application to complete unstaking
}

// MarshalJSON marshals the application to JSON using Bech32
func (a Application) MarshalJSON() ([]byte, error) {
	bechConsPubKey, err := sdk.Bech32ifyConsPub(a.ConsPubKey)
	if err != nil {
		return nil, err
	}
	return codec.Cdc.MarshalJSON(bechApplication{
		Address:                 a.Address,
		ConsPubKey:              bechConsPubKey,
		Jailed:                  a.Jailed,
		Status:                  a.Status,
		StakedTokens:            a.StakedTokens,
		UnstakingCompletionTime: a.UnstakingCompletionTime,
	})
}

// UnmarshalJSON unmarshals the application from JSON using Bech32
func (a *Application) UnmarshalJSON(data []byte) error {
	bv := &bechApplication{}
	if err := codec.Cdc.UnmarshalJSON(data, bv); err != nil {
		return err
	}
	consPubKey, err := sdk.GetConsPubKeyBech32(bv.ConsPubKey)
	if err != nil {
		return err
	}
	*a = Application{
		Address:                 bv.Address,
		ConsPubKey:              consPubKey,
		Jailed:                  bv.Jailed,
		StakedTokens:            bv.StakedTokens,
		Status:                  bv.Status,
		UnstakingCompletionTime: bv.UnstakingCompletionTime,
	}
	return nil
}
