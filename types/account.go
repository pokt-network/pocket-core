package model

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

// "Account" is the base structure for actors in the Pocket Network
type Account struct {
	Address     sdk.ValAddress `json:"address"`
	PubKey      crypto.PubKey  `json:"publicKey"`
	Balance     sdk.Coins      `json:"balance"`
	StakeAmount sdk.Coins      `json:"stakeamount"`
}
