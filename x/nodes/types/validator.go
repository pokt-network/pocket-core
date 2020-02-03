package types

import (
	"bytes"
	"fmt"
	"github.com/pokt-network/posmint/crypto"
	"time"

	sdk "github.com/pokt-network/posmint/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

type Validator struct {
	Address                 sdk.Address      `json:"address" yaml:"address"`               // address of the validator; hex encoded in JSON
	PublicKey               crypto.PublicKey `json:"public_key" yaml:"public_key"`         // the consensus public key of the validator; hex encoded in JSON
	Jailed                  bool             `json:"jailed" yaml:"jailed"`                 // has the validator been jailed from staked status?
	Status                  sdk.StakeStatus  `json:"status" yaml:"status"`                 // validator status (staked/unstaking/unstaked)
	Chains                  []string         `json:"chains" yaml:"chains"`                 // validator non native blockchains
	ServiceURL              string           `json:"service_url" yaml:"service_url"`       // url where the pocket service api is hosted
	StakedTokens            sdk.Int          `json:"tokens" yaml:"tokens"`                 // tokens staked in the network
	UnstakingCompletionTime time.Time        `json:"unstaking_time" yaml:"unstaking_time"` // if unstaking, min time for the validator to complete unstaking
}

// NewValidator - initialize a new validator
func NewValidator(addr sdk.Address, consPubKey crypto.PublicKey, chains []string, serviceURL string, tokensToStake sdk.Int) Validator {
	return Validator{
		Address:                 addr,
		PublicKey:               consPubKey,
		Jailed:                  false,
		Status:                  sdk.Staked,
		Chains:                  chains,
		StakedTokens:            tokensToStake,
		ServiceURL:              serviceURL,
		UnstakingCompletionTime: time.Unix(0, 0).UTC(), // zero out because status: staked
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

// ABCIValidatorUpdateZero returns an abci.ValidatorUpdate from a staking validator type
// with zero power used for validator updates.
func (v Validator) ABCIValidatorUpdateZero() abci.ValidatorUpdate {
	return abci.ValidatorUpdate{
		PubKey: tmtypes.TM2PB.PubKey(v.PublicKey.PubKey()),
		Power:  0,
	}
}

// get the consensus-engine power
// a reduction of 10^6 from validator tokens is applied
func (v Validator) ConsensusPower() int64 {
	if v.IsStaked() {
		return v.PotentialConsensusPower()
	}
	return 0
}

// potential consensus-engine power
func (v Validator) PotentialConsensusPower() int64 {
	return sdk.TokensToConsensusPower(v.StakedTokens)
}

// RemoveStakedTokens removes tokens from a validator
func (v Validator) RemoveStakedTokens(tokens sdk.Int) Validator {
	if tokens.IsNegative() {
		panic(fmt.Sprintf("should not happen: trying to remove negative tokens %v", tokens))
	}
	if v.StakedTokens.LT(tokens) {
		panic(fmt.Sprintf("should not happen: only have %v tokens, trying to remove %v", v.StakedTokens, tokens))
	}
	v.StakedTokens = v.StakedTokens.Sub(tokens)
	return v
}

// AddStakedTokens tokens to staked field for a validator
func (v Validator) AddStakedTokens(tokens sdk.Int) Validator {
	if tokens.IsNegative() {
		panic(fmt.Sprintf("should not happen: trying to add negative tokens %v", tokens))
	}
	v.StakedTokens = v.StakedTokens.Add(tokens)
	return v
}

// compares the vital fields of two validator structures
func (v Validator) Equals(v2 Validator) bool {
	return v.PublicKey.Equals(v2.PublicKey) &&
		bytes.Equal(v.Address, v2.Address) &&
		v.Status.Equal(v2.Status) &&
		v.StakedTokens.Equal(v2.StakedTokens)
}

// UpdateStatus updates the staking status
func (v Validator) UpdateStatus(newStatus sdk.StakeStatus) Validator {
	v.Status = newStatus
	return v
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
func (v Validator) GetTokens() sdk.Int             { return v.StakedTokens }
func (v Validator) GetConsensusPower() int64       { return v.ConsensusPower() }
