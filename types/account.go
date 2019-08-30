package types

import (
	"github.com/pokt-network/pocket-core/crypto"
)

// "Account" is the base structure for actors in the Pocket Network
type Account struct {
	Address     Address          `json:"address"`
	PubKey      crypto.PublicKey `json:"publicKey"`
	Balance     POKT             `json:"balance"`
	StakeAmount POKT             `json:"stakeamount"`
}
