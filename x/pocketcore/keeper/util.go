package keeper

import (
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
)

func (k Keeper) GetPKFromFile(ctx sdk.Ctx) (crypto.PrivateKey, error) {
	pvKey, err := types.GetPvKeyFile()
	if err != nil {
		return nil, err
	}
	pk, _ := crypto.PrivKeyToPrivateKey(pvKey.PrivKey)
	return pk, nil
}
