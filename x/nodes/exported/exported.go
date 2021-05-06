package exported

import (
	"github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
)

// ValidatorI expected validator functions
type ValidatorI interface {
	IsStaked() bool                 // check if has a staked status
	IsUnstaked() bool               // check if has status unstaked
	IsUnstaking() bool              // check if has status unstaking
	IsJailed() bool                 // whether the validator is jailed
	GetStatus() sdk.StakeStatus     // status of the validator
	GetAddress() sdk.Address        // operator address to receive/return validators coins
	GetPublicKey() crypto.PublicKey // validation consensus pubkey
	GetTokens() sdk.BigInt          // validation tokens
	GetConsensusPower() int64       // validation power in tendermint
	GetChains() []string            // get chains staked for validator
}
