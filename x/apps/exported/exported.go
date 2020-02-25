package exported

import (
	"github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
)

// ApplicationI expected application functions
type ApplicationI interface {
	IsJailed() bool                 // whether the application is jailed
	GetStatus() sdk.StakeStatus     // status of the application
	IsStaked() bool                 // check if has a staked status
	IsUnstaked() bool               // check if has status unstaked
	IsUnstaking() bool              // check if has status unstaking
	GetChains() []string            // retrieve the staked chains
	GetAddress() sdk.Address        // operator address to receive/return applications coins
	GetPublicKey() crypto.PublicKey // validation consensus pubkey
	GetTokens() sdk.Int             // validation tokens
	GetMaxRelays() sdk.Int          // maximum relays
}
