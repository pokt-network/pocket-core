package types

import (
	"bytes"
	"fmt"
	"time"

	sdk "github.com/pokt-network/posmint/types"
	"github.com/tendermint/tendermint/crypto"
)

type Application struct {
	Address                 sdk.ValAddress      `json:"address" yaml:"address"`               // address of the application; bech encoded in JSON
	ConsPubKey              crypto.PubKey       `json:"cons_pubkey" yaml:"cons_pubkey"`       // the consensus public key of the application; bech encoded in JSON
	Jailed                  bool                `json:"jailed" yaml:"jailed"`                 // has the application been jailed from bonded status?
	Status                  sdk.BondStatus      `json:"status" yaml:"status"`                 // application status (bonded/unbonding/unbonded)
	Chains                  map[string]struct{} `json:"chains" yaml:"chains"`                 // requested chains
	StakedTokens            sdk.Int             `json:"Tokens" yaml:"Tokens"`                 // tokens staked in the network
	MaxRelays               sdk.Int             `json:"max_relays" yaml:"max_relays"`         // maximum number of relays allowed
	UnstakingCompletionTime time.Time           `json:"unstaking_time" yaml:"unstaking_time"` // if unstaking, min time for the application to complete unstaking
}

// NewApplication - initialize a new application
func NewApplication(addr sdk.ValAddress, consPubKey crypto.PubKey, chains map[string]struct{}, tokensToStake sdk.Int) Application {
	return Application{
		Address:                 addr,
		ConsPubKey:              consPubKey,
		Jailed:                  false,
		Status:                  sdk.Bonded,
		Chains:                  chains,
		StakedTokens:            tokensToStake,
		UnstakingCompletionTime: time.Unix(0, 0).UTC(), // zero out because status: bonded
	}
}

// get the consensus-engine power
// a reduction of 10^6 from application tokens is applied
func (a Application) ConsensusPower() int64 {
	if a.IsStaked() {
		return a.PotentialConsensusPower()
	}
	return 0
}

// potential consensus-engine power
func (a Application) PotentialConsensusPower() int64 {
	return sdk.TokensToConsensusPower(a.StakedTokens)
}

// RemoveStakedTokens removes tokens from a application
func (a Application) RemoveStakedTokens(tokens sdk.Int) Application {
	if tokens.IsNegative() {
		panic(fmt.Sprintf("should not happen: trying to remove negative tokens %v", tokens))
	}
	if a.StakedTokens.LT(tokens) {
		panic(fmt.Sprintf("should not happen: only have %v tokens, trying to remove %v", a.StakedTokens, tokens))
	}
	a.StakedTokens = a.StakedTokens.Sub(tokens)
	return a
}

// AddStakedTokens tokens to staked field for a application
func (a Application) AddStakedTokens(tokens sdk.Int) Application {
	if tokens.IsNegative() {
		panic(fmt.Sprintf("should not happen: trying to add negative tokens %v", tokens))
	}
	a.StakedTokens = a.StakedTokens.Add(tokens)
	return a
}

// compares the vital fields of two application structures
func (a Application) Equals(v2 Application) bool {
	return a.ConsPubKey.Equals(v2.ConsPubKey) &&
		bytes.Equal(a.Address, v2.Address) &&
		a.Status.Equal(v2.Status) &&
		a.StakedTokens.Equal(v2.StakedTokens)
}

// UpdateStatus updates the staking status
func (a Application) UpdateStatus(newStatus sdk.BondStatus) Application {
	a.Status = newStatus
	return a
}

// return the TM application address
func (a Application) ConsAddress() sdk.ConsAddress   { return sdk.ConsAddress(a.ConsPubKey.Address()) }
func (a Application) GetChains() map[string]struct{} { return a.Chains }
func (a Application) IsStaked() bool                 { return a.GetStatus().Equal(sdk.Bonded) }
func (a Application) IsUnstaked() bool               { return a.GetStatus().Equal(sdk.Unbonded) }
func (a Application) IsUnstaking() bool              { return a.GetStatus().Equal(sdk.Unbonding) }
func (a Application) IsJailed() bool                 { return a.Jailed }
func (a Application) GetStatus() sdk.BondStatus      { return a.Status }
func (a Application) GetAddress() sdk.ValAddress     { return a.Address }
func (a Application) GetConsPubKey() crypto.PubKey   { return a.ConsPubKey }
func (a Application) GetConsAddr() sdk.ConsAddress   { return sdk.ConsAddress(a.ConsPubKey.Address()) }
func (a Application) GetTokens() sdk.Int             { return a.StakedTokens }
func (a Application) GetConsensusPower() int64       { return a.ConsensusPower() }
func (a Application) GetMaxRelays() sdk.Int          { return a.MaxRelays }
