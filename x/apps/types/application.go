package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pokt-network/posmint/crypto"

	sdk "github.com/pokt-network/posmint/types"
)

// Application represents a pocket network decentralized application. Applications stake in the network for relay throughput.
type Application struct {
	Address                 sdk.Address      `json:"address" yaml:"address"`               // address of the application; hex encoded in JSON
	PublicKey               crypto.PublicKey `json:"public_key" yaml:"public_key"`         // the public key of the application; hex encoded in JSON
	Jailed                  bool             `json:"jailed" yaml:"jailed"`                 // has the application been jailed from staked status?
	Status                  sdk.StakeStatus  `json:"status" yaml:"status"`                 // application status (staked/unstaking/unstaked)
	Chains                  []string         `json:"chains" yaml:"chains"`                 // requested chains
	StakedTokens            sdk.Int          `json:"tokens" yaml:"tokens"`                 // tokens staked in the network
	MaxRelays               sdk.Int          `json:"max_relays" yaml:"max_relays"`         // maximum number of relays allowed
	UnstakingCompletionTime time.Time        `json:"unstaking_time" yaml:"unstaking_time"` // if unstaking, min time for the application to complete unstaking
}

type ApplicationsPage struct {
	Result Applications `json:"result"`
	Total  int          `json:"total_pages"`
	Page   int          `json:"page"`
}

// Marshals struct into JSON
func (aP ApplicationsPage) JSON() (out []byte, err error) {
	// each element should be a JSON
	return json.Marshal(aP)
}

// String returns a human readable string representation of a validator page
func (aP ApplicationsPage) String() string {
	return fmt.Sprintf("Total:\t\t%d\nPage:\t\t%d\nResult:\t\t\n====\n%s\n====\n", aP.Total, aP.Page, aP.Result.String())
}

// NewApplication - initialize a new instance of an application
func NewApplication(addr sdk.Address, publicKey crypto.PublicKey, chains []string, tokensToStake sdk.Int) Application {
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
func (a Application) RemoveStakedTokens(tokens sdk.Int) (Application, error) {
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
func (a Application) AddStakedTokens(tokens sdk.Int) (Application, error) {
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
func (a Application) GetTokens() sdk.Int             { return a.StakedTokens }
func (a Application) GetConsensusPower() int64       { return a.ConsensusPower() }
func (a Application) GetMaxRelays() sdk.Int          { return a.MaxRelays }
