package keeper

import (
	"github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
)

// "GetPKFromFile" - Returns the private key object from a file
func (k Keeper) GetPKFromFile(ctx sdk.Ctx) (crypto.PrivateKey, error) {
	// get the Private validator key from the file
	pvKey, err := types.GetPVKeyFile()
	if err != nil {
		return nil, err
	}
	// convert the privKey to a private key object (compatible interface)
	pk, er := crypto.PrivKeyToPrivateKey(pvKey.PrivKey)
	if er != nil {
		return nil, er
	}
	return pk, nil
}

var PkFromAddressMap map[string]crypto.PrivateKey

func InitializePkMap() error {
	PkFromAddressMap = make(map[string]crypto.PrivateKey)
	for _, pvKey := range types.GlobalPVKeyFiles {
		pk, er := crypto.PrivKeyToPrivateKey(pvKey.PrivKey)

		address := sdk.Address(pk.PublicKey().Address())
		if er != nil {
			return er
		}
		PkFromAddressMap[address.String()] = pk
	}
	return nil
}

func (k Keeper) GetPkFromAddress(address *sdk.Address) (crypto.PrivateKey, error) {
	if PkFromAddressMap == nil {
		err := InitializePkMap()
		if err != nil {
			return nil, err
		}
	}
	return PkFromAddressMap[address.String()], nil
}
