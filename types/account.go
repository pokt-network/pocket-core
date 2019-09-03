package types

import "encoding/hex"

// "Account" is the base structure for actors in the Pocket Network
type Account struct {
	Address     Address          `json:"address"`
	PubKey      AccountPublicKey `json:"publicKey"`
	Balance     POKT             `json:"balance"`
	StakeAmount POKT             `json:"stakeamount"`
}

type AccountPublicKey string // hex string

func (a AccountPublicKey) Bytes() ([]byte, error) {
	return hex.DecodeString(string(a))
}
