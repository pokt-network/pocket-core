package types

import (
	"encoding/hex"
	"encoding/json"
	sdk "github.com/pokt-network/posmint/types"
)

// "NonNativeChain" - A strucutre used to identify a non native (external) blockchain on the pocket network
type NonNativeChain struct {
	Ticker  string `json:"ticker"` // market identifier (ETH BTC)
	Netid   string `json:"netid"`
	Version string `json:"version"`
	Client  string `json:"client"`
	Inter   string `json:"interface"`
}

// "Bytes"- Converts the non native chains to bytes
func (c NonNativeChain) Bytes() ([]byte, sdk.Error) {
	// ensure essential fields are not empty
	if c.Ticker == "" || c.Netid == "" || c.Version == "" {
		return nil, NewInvalidChainParamsError(ModuleName)
	}
	// marshal into json bz
	res, er := json.Marshal(c)
	if er != nil {
		return nil, NewJSONMarshalError(ModuleName, er)
	}
	// return the bz
	return res, nil
}

// "ID" - Hashes the bytes of the non native chain
func (c NonNativeChain) Hash() ([]byte, sdk.Error) {
	// get the bz of the chain
	res, err := c.Bytes()
	if err != nil {
		return nil, err
	}
	// return the short hash of the chain
	return ShortHash(res), nil
}

// "HashString" - Returns a hex string of the non native chain addr.
func (c NonNativeChain) HashString() (string, sdk.Error) {
	// get the hash of the chain
	res, err := c.Hash()
	if err != nil {
		return "", err
	}
	// hex encode into a string
	return hex.EncodeToString(res), nil
}
