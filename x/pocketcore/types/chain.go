package types

import (
	"encoding/hex"
	"encoding/json"
	sdk "github.com/pokt-network/posmint/types"
)

type NonNativeChain struct {
	Ticker  string
	Netid   string
	Version string
	Client  string
	Inter   string
}

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

func (c NonNativeChain) Hash() ([]byte, sdk.Error) {
	res, err := c.Bytes()
	if err != nil {
		return nil, err
	}
	return SHA3FromBytes(res), nil
}

func (c NonNativeChain) HashString() (string, sdk.Error) {
	res, err := c.Hash()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(res), nil
}
