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
)

// Application represents a pocket network decentralized application. Applications stake in the network for relay throughput.
type Application struct {
	Address                 sdk.Address      `json:"address" yaml:"address"`               // address of the application; hex encoded in JSON
	PublicKey               crypto.PublicKey `json:"public_key" yaml:"public_key"`         // the public key of the application; hex encoded in JSON
	Jailed                  bool             `json:"jailed" yaml:"jailed"`                 // has the application been jailed from staked status?
	Status                  sdk.StakeStatus  `json:"status" yaml:"status"`                 // application status (staked/unstaking/unstaked)
	Chains                  []string         `json:"chains" yaml:"chains"`                 // requested chains
	StakedTokens            sdk.BigInt       `json:"tokens" yaml:"tokens"`                 // tokens staked in the network
	MaxRelays               sdk.BigInt       `json:"max_relays" yaml:"max_relays"`         // maximum number of relays allowed
	UnstakingCompletionTime time.Time        `json:"unstaking_time" yaml:"unstaking_time"` // if unstaking, min time for the application to complete unstaking
}

// NewApplication - initialize a new instance of an application
func NewApplication(addr sdk.Address, publicKey crypto.PublicKey, chains []string, tokensToStake sdk.BigInt) Application {
	return Application{
		Address:                 addr,
		PublicKey:               publicKey,
		Jailed:                  false,
		Status:                  sdk.Staked,
		Chains:                  chains,
		StakedTokens:            tokensToStake,
		UnstakingCompletionTime: time.Time{}, // zero out because status: staked
	}
}

// get the consensus-engine power
// a reduction of 10^6 from application tokens is applied
func (a Application) ConsensusPower() int64 {
	if a.IsStaked() {
		return sdk.TokensToConsensusPower(a.StakedTokens)
	}
	return 0
}

// RemoveStakedTokens removes tokens from a application
func (a Application) RemoveStakedTokens(tokens sdk.BigInt) (Application, error) {
	if tokens.IsNegative() {
		return Application{}, fmt.Errorf("should not happen: trying to remove negative tokens %v", tokens)
	}
	if a.StakedTokens.LT(tokens) {
		return Application{}, fmt.Errorf("should not happen: only have %v tokens, trying to remove %v", a.StakedTokens, tokens)
	}
	a.StakedTokens = a.StakedTokens.Sub(tokens)
	return a, nil
}

// AddStakedTokens tokens to staked field for a application
func (a Application) AddStakedTokens(tokens sdk.BigInt) (Application, error) {
	if tokens.IsNegative() {
		return Application{}, fmt.Errorf("should not happen: trying to remove negative tokens %v", tokens)
	}
	a.StakedTokens = a.StakedTokens.Add(tokens)
	return a, nil
}

// compares the vital fields of two application structures
func (a Application) Equals(v2 Application) bool {
	return a.PublicKey.Equals(v2.PublicKey) &&
		bytes.Equal(a.Address, v2.Address) &&
		a.Status.Equal(v2.Status) &&
		a.StakedTokens.Equal(v2.StakedTokens)
}

// UpdateStatus updates the staking status
func (a Application) UpdateStatus(newStatus sdk.StakeStatus) Application {
	a.Status = newStatus
	return a
}

func (a Application) GetChains() []string            { return a.Chains }
func (a Application) IsStaked() bool                 { return a.GetStatus().Equal(sdk.Staked) }
func (a Application) IsUnstaked() bool               { return a.GetStatus().Equal(sdk.Unstaked) }
func (a Application) IsUnstaking() bool              { return a.GetStatus().Equal(sdk.Unstaking) }
func (a Application) IsJailed() bool                 { return a.Jailed }
func (a Application) GetStatus() sdk.StakeStatus     { return a.Status }
func (a Application) GetAddress() sdk.Address        { return a.Address }
func (a Application) GetPublicKey() crypto.PublicKey { return a.PublicKey }
func (a Application) GetTokens() sdk.BigInt          { return a.StakedTokens }
func (a Application) GetConsensusPower() int64       { return a.ConsensusPower() }
func (a Application) GetMaxRelays() sdk.BigInt       { return a.MaxRelays }

var _ codec.ProtoMarshaler = &Application{}

func (a *Application) Reset() {
	*a = Application{}
}

func (a Application) ProtoMessage() {
	p := a.ToProto()
	p.ProtoMessage()
}

func (a Application) Marshal() ([]byte, error) {
	p := a.ToProto()
	return p.Marshal()
}

func (a Application) MarshalTo(data []byte) (n int, err error) {
	p := a.ToProto()
	return p.MarshalTo(data)
}

func (a Application) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	p := a.ToProto()
	return p.MarshalToSizedBuffer(dAtA)
}

func (a Application) Size() int {
	p := a.ToProto()
	return p.Size()
}

func (a *Application) Unmarshal(data []byte) (err error) {
	var pa ProtoApplication
	err = pa.Unmarshal(data)
	if err != nil {
		return err
	}
	*a, err = pa.FromProto()
	return
}

func (a Application) ToProto() ProtoApplication {
	return ProtoApplication{
		Address:                 a.Address,
		PublicKey:               a.PublicKey.RawBytes(),
		Jailed:                  a.Jailed,
		Status:                  a.Status,
		Chains:                  a.Chains,
		StakedTokens:            a.StakedTokens,
		MaxRelays:               a.MaxRelays,
		UnstakingCompletionTime: a.UnstakingCompletionTime,
	}
}

func (ae ProtoApplication) FromProto() (Application, error) {
	pk, err := crypto.NewPublicKeyBz(ae.PublicKey)
	if err != nil {
		return Application{}, err
	}
	return Application{
		Address:                 ae.Address,
		PublicKey:               pk,
		Jailed:                  ae.Jailed,
		Status:                  ae.Status,
		Chains:                  ae.Chains,
		StakedTokens:            ae.StakedTokens,
		MaxRelays:               ae.MaxRelays,
		UnstakingCompletionTime: ae.UnstakingCompletionTime,
	}, nil
}

// Applications is a slice of type application.
type Applications []Application

func (a Applications) String() (out string) {
	for _, val := range a {
		out += val.String() + "\n\n"
	}
	return strings.TrimSpace(out)
}

// String returns a human readable string representation of a application.
func (a Application) String() string {
	return fmt.Sprintf("Address:\t\t%s\nPublic Key:\t\t%s\nJailed:\t\t\t%v\nChains:\t\t\t%v\nMaxRelays:\t\t%v\nStatus:\t\t\t%s\nTokens:\t\t\t%s\nUnstaking Time:\t%v\n----\n",
		a.Address, a.PublicKey.RawString(), a.Jailed, a.Chains, a.MaxRelays, a.Status, a.StakedTokens, a.UnstakingCompletionTime,
	)
}

// this is a helper struct used for JSON de- and encoding only
type JSONApplication struct {
	Address                 sdk.Address     `json:"address" yaml:"address"`               // the hex address of the application
	PublicKey               string          `json:"public_key" yaml:"public_key"`         // the hex consensus public key of the application
	Jailed                  bool            `json:"jailed" yaml:"jailed"`                 // has the application been jailed from staked status?
	Chains                  []string        `json:"chains" yaml:"chains"`                 // non native (external) blockchains needed for the application
	MaxRelays               sdk.BigInt      `json:"max_relays" yaml:"max_relays"`         // maximum number of relays allowed for the application
	Status                  sdk.StakeStatus `json:"status" yaml:"status"`                 // application status (staked/unstaking/unstaked)
	StakedTokens            sdk.BigInt      `json:"staked_tokens" yaml:"staked_tokens"`   // how many staked tokens
	UnstakingCompletionTime time.Time       `json:"unstaking_time" yaml:"unstaking_time"` // if unstaking, min time for the application to complete unstaking
}

// marshal structure into JSON encoding
func (a Applications) JSON() (out []byte, err error) {
	return json.Marshal(a)
}

// MarshalJSON marshals the application to JSON using raw Hex for the public key
func (a Application) MarshalJSON() ([]byte, error) {
	return ModuleCdc.MarshalJSON(JSONApplication{
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
	bv := &JSONApplication{}
	if err := ModuleCdc.UnmarshalJSON(data, bv); err != nil {
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

// unmarshal the application
func MarshalApplication(cdc *codec.Codec, ctx sdk.Ctx, application Application) (result []byte, err error) {
	return cdc.MarshalBinaryLengthPrefixed(&application, ctx.BlockHeight())
}

// unmarshal the application
func UnmarshalApplication(cdc *codec.Codec, ctx sdk.Ctx, appBytes []byte) (application Application, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(appBytes, &application, ctx.BlockHeight())
	return
}

type ApplicationsPage struct {
	Result Applications `json:"result"`
	Total  int          `json:"total_pages"`
	Page   int          `json:"page"`
}

// String returns a human readable string representation of a validator page
func (aP ApplicationsPage) String() string {
	return fmt.Sprintf("Total:\t\t%d\nPage:\t\t%d\nResult:\t\t\n====\n%s\n====\n", aP.Total, aP.Page, aP.Result.String())
}
