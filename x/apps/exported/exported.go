package exported

import (
	sdk "github.com/pokt-network/posmint/types"
	"github.com/tendermint/tendermint/crypto"
)

// ApplicationI expected application functions
type ApplicationI interface {
	IsJailed() bool                 // whether the application is jailed
	GetStatus() sdk.BondStatus      // status of the application
	IsStaked() bool                 // check if has a bonded status
	IsUnstaked() bool               // check if has status unbonded
	IsUnstaking() bool              // check if has status unbonding
	GetChains() map[string]struct{} // retrieve the staked chains
	GetAddress() sdk.ValAddress     // operator address to receive/return applications coins
	GetConsPubKey() crypto.PubKey   // validation consensus pubkey
	GetConsAddr() sdk.ConsAddress   // validation consensus address
	GetTokens() sdk.Int             // validation tokens
}
