package types

import (
	"github.com/pokt-network/pocket-core/crypto"
)

// "Account" is the base structure for actors in the Pocket Network
type Account struct {
	Address     Address          `json:"address"`
	PubKey      crypto.PublicKey `json:"publicKey"`
	Balance     Coins            `json:"balance"`
	StakeAmount Coins            `json:"stakeamount"`
}
