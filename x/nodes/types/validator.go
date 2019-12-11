package types

import (
	"bytes"
	"fmt"
	"time"

	sdk "github.com/pokt-network/posmint/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	tmtypes "github.com/tendermint/tendermint/types"
)

type Validator struct {
	Address                 sdk.ValAddress `json:"address" yaml:"address"`               // address of the validator; bech encoded in JSON
	ConsPubKey              crypto.PubKey  `json:"cons_pubkey" yaml:"cons_pubkey"`       // the consensus public key of the validator; bech encoded in JSON
	Jailed                  bool           `json:"jailed" yaml:"jailed"`                 // has the validator been jailed from bonded status?
	Status                  sdk.BondStatus `json:"status" yaml:"status"`                 // validator status (bonded/unbonding/unbonded)
	Chains                  []string       `json:"chains" yaml:"chains"`                 // validator non native blockchains
	ServiceURL              string         `json:"serviceurl" yaml:"serviceurl"`         // url where the pocket service api is hosted
	StakedTokens            sdk.Int        `json:"Tokens" yaml:"Tokens"`                 // tokens staked in the network
	UnstakingCompletionTime time.Time      `json:"unstaking_time" yaml:"unstaking_time"` // if unstaking, min time for the validator to complete unstaking
}

// NewValidator - initialize a new validator
func NewValidator(addr sdk.ValAddress, consPubKey crypto.PubKey, chains []string, serviceURL string, tokensToStake sdk.Int) Validator {
	return Validator{
		Address:                 addr,
		ConsPubKey:              consPubKey,
		Jailed:                  false,
		Status:                  sdk.Bonded,
		Chains:                  chains,
		StakedTokens:            tokensToStake,
		ServiceURL:              serviceURL,
		UnstakingCompletionTime: time.Unix(0, 0).UTC(), // zero out because status: bonded
	}
}

// ABCIValidatorUpdate returns an abci.ValidatorUpdate from a staking validator type
// with the full validator power
func (v Validator) ABCIValidatorUpdate() abci.ValidatorUpdate {
	return abci.ValidatorUpdate{
		PubKey: tmtypes.TM2PB.PubKey(v.ConsPubKey),
		Power:  v.ConsensusPower(),
	}
}

// ABCIValidatorUpdateZero returns an abci.ValidatorUpdate from a staking validator type
// with zero power used for validator updates.
func (v Validator) ABCIValidatorUpdateZero() abci.ValidatorUpdate {
	return abci.ValidatorUpdate{
		PubKey: tmtypes.TM2PB.PubKey(v.ConsPubKey),
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
	return v.ConsPubKey.Equals(v2.ConsPubKey) &&
		bytes.Equal(v.Address, v2.Address) &&
		v.Status.Equal(v2.Status) &&
		v.StakedTokens.Equal(v2.StakedTokens)
}

// UpdateStatus updates the staking status
func (v Validator) UpdateStatus(newStatus sdk.BondStatus) Validator {
	v.Status = newStatus
	return v
}

// return the TM validator address
func (v Validator) ConsAddress() sdk.ConsAddress { return sdk.ConsAddress(v.ConsPubKey.Address()) }
func (v Validator) GetChains() []string          { return v.Chains }
func (v Validator) GetServiceURL() string        { return v.ServiceURL }
func (v Validator) IsStaked() bool               { return v.GetStatus().Equal(sdk.Bonded) }
func (v Validator) IsUnstaked() bool             { return v.GetStatus().Equal(sdk.Unbonded) }
func (v Validator) IsUnstaking() bool            { return v.GetStatus().Equal(sdk.Unbonding) }
func (v Validator) IsJailed() bool               { return v.Jailed }
func (v Validator) GetStatus() sdk.BondStatus    { return v.Status }
func (v Validator) GetAddress() sdk.ValAddress   { return v.Address }
func (v Validator) GetConsPubKey() crypto.PubKey { return v.ConsPubKey }
func (v Validator) GetConsAddr() sdk.ConsAddress { return sdk.ConsAddress(v.ConsPubKey.Address()) }
func (v Validator) GetTokens() sdk.Int           { return v.StakedTokens }
func (v Validator) GetConsensusPower() int64     { return v.ConsensusPower() }
