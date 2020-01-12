package types

import (
	"encoding/hex"
	"encoding/json"
	sdk "github.com/pokt-network/posmint/types"
)

// strucutre used to identify a non native (external) blockchain on the pocket network
type NonNativeChain struct {
	Ticker  string `json:"ticker"`
	Netid   string `json:"netid"`
	Version string `json:"version"`
	Client  string `json:"client"`
	Inter   string `json:"interface"`
}

// converts the non native chains to bytes
func (c NonNativeChain) Bytes() ([]byte, sdk.Error) {
	if c.Ticker == "" || c.Netid == "" || c.Version == "" {
		return nil, NewInvalidChainParamsError(ModuleName)
	}
	res, er := json.Marshal(c)
	if er != nil {
		return nil, NewJSONMarshalError(ModuleName, er)
	}
	return res, nil
}

// hashes the bytes of the non native hcain
func (c NonNativeChain) Hash() ([]byte, sdk.Error) {
	res, err := c.Bytes()
	if err != nil {
		return nil, err
	}
	return Hash(res), nil
}

// returns a hex string of the non native chain addr
func (c NonNativeChain) HashString() (string, sdk.Error) {
	res, err := c.Hash()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(res), nil
}
