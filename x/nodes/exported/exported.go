package exported

import (
	"github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
)

// ValidatorI expected validator functions
type ValidatorI interface {
	IsJailed() bool                 // whether the validator is jailed
	GetStatus() sdk.StakeStatus     // status of the validator
	IsStaked() bool                 // check if has a staked status
	IsUnstaked() bool               // check if has status unstaked
	IsUnstaking() bool              // check if has status unstaking
	GetChains() []string            // retrieve the staked chains
	GetServiceURL() string          // retrieve the url for pocket core service api
	GetAddress() sdk.Address        // address to receive/return validators coins
	GetPublicKey() crypto.PublicKey // validator public key
	GetTokens() sdk.Int             // validator tokens
	GetConsensusPower() int64       // validator power in tendermint
}
