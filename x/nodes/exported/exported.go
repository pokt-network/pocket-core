package exported

import (
	sdk "github.com/pokt-network/posmint/types"
	"github.com/tendermint/tendermint/crypto"
)

// ValidatorI expected validator functions
type ValidatorI interface {
	IsJailed() bool               // whether the validator is jailed
	GetStatus() sdk.BondStatus    // status of the validator
	IsStaked() bool               // check if has a bonded status
	IsUnstaked() bool             // check if has status unbonded
	IsUnstaking() bool            // check if has status unbonding
	GetChains() []string          // retrieve the staked chains
	GetServiceURL() string        // retrieve the url for pocket core service api
	GetAddress() sdk.ValAddress   // operator address to receive/return validators coins
	GetConsPubKey() crypto.PubKey // validation consensus pubkey
	GetConsAddr() sdk.ConsAddress // validation consensus address
	GetTokens() sdk.Int           // validation tokens
	GetConsensusPower() int64     // validation power in tendermint
}
