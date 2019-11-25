package keeper

import (
	sdk "github.com/pokt-network/posmint/types"
	"github.com/tendermint/tendermint/crypto"
)

func (k Keeper) PubKeyFromString(pubKey string) (crypto.PubKey, error) { // todo
	return nil, nil
}

func (k Keeper) AddressFromPubKeyString(pubKey string) (sdk.ValAddress, error) {
	pk, err := k.PubKeyFromString(pubKey)
	if err != nil {
		return nil, nil
	}
	return sdk.ValAddress(pk.Address()), nil
}
