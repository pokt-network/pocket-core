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

// Applications is a slice of type application.
type Applications []Application

func (a Applications) String() (out string) {
	for _, val := range a {
		out += val.String() + "\n\n"
	}
	return strings.TrimSpace(out)
}

// HashString returns a human readable string representation of a application.
func (a Application) String() string {
	return fmt.Sprintf("Address:\t\t%s\nPublic Key:\t\t%s\nJailed:\t\t\t%v\nChains:\t\t\t%v\nMaxRelays:\t\t%d\nStatus:\t\t\t%s\nTokens:\t\t\t%s\nUnstaking Time:\t%v",
		a.Address, a.PublicKey.RawString(), a.Jailed, a.Chains, a.MaxRelays.Int64(), a.Status, a.StakedTokens, a.UnstakingCompletionTime,
	)
}

// this is a helper struct used for JSON de- and encoding only
type hexApplication struct {
	Address                 sdk.Address     `json:"address" yaml:"address"`               // the hex address of the application
	PublicKey               string          `json:"public_key" yaml:"public_key"`         // the hex consensus public key of the application
	Jailed                  bool            `json:"jailed" yaml:"jailed"`                 // has the application been jailed from staked status?
	Chains                  []string        `json:"chains" yaml:"chains"`                 // non native (external) blockchains needed for the application
	MaxRelays               sdk.Int         `json:"max_relays" yaml:"max_relays"`         // maximum number of relays allowed for the application
	Status                  sdk.StakeStatus `json:"status" yaml:"status"`                 // application status (staked/unstaking/unstaked)
	StakedTokens            sdk.Int         `json:"staked_tokens" yaml:"staked_tokens"`   // how many staked tokens
	UnstakingCompletionTime time.Time       `json:"unstaking_time" yaml:"unstaking_time"` // if unstaking, min time for the application to complete unstaking
}

// marshal structure into JSON encoding
func (a Applications) JSON() (out []byte, err error) {
	return json.Marshal(a)
}

// MarshalJSON marshals the application to JSON using raw Hex for the public key
func (a Application) MarshalJSON() ([]byte, error) {
	return codec.Cdc.MarshalJSON(hexApplication{
		Address:                 a.Address,
		PublicKey:               a.PublicKey.RawString(),
		Jailed:                  a.Jailed,
		Status:                  a.Status,
		Chains:                  a.Chains,
		MaxRelays:               a.MaxRelays,
		StakedTokens:            a.StakedTokens,
		UnstakingCompletionTime: a.UnstakingCompletionTime,
	})
}

// UnmarshalJSON unmarshals the application from JSON using raw hex for the public key
func (a *Application) UnmarshalJSON(data []byte) error {
	bv := &hexApplication{}
	if err := codec.Cdc.UnmarshalJSON(data, bv); err != nil {
		return err
	}
	consPubKey, err := crypto.NewPublicKey(bv.PublicKey)
	if err != nil {
		return err
	}
	*a = Application{
		Address:                 bv.Address,
		PublicKey:               consPubKey,
		Chains:                  bv.Chains,
		MaxRelays:               bv.MaxRelays,
		Jailed:                  bv.Jailed,
		StakedTokens:            bv.StakedTokens,
		Status:                  bv.Status,
		UnstakingCompletionTime: bv.UnstakingCompletionTime,
	}
	return nil
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
