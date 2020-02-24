package keeper

import (
	"fmt"
	"github.com/pokt-network/pocket-core/x/nodes/exported"
	pc "github.com/pokt-network/pocket-core/x/pocketcore/types"
	sdk "github.com/pokt-network/posmint/types"
)

// get all nodes from the world state
func (k Keeper) GetAllNodes(ctx sdk.Context) []exported.ValidatorI {
	validators := k.posKeeper.GetStakedValidators(ctx)
	ctx.Logger().Info(fmt.Sprintf("GetNodes() = %v", validators))
	return validators
}

// get a node from the world state
func (k Keeper) GetNode(ctx sdk.Context, address sdk.Address) (n exported.ValidatorI, found bool) {
	ctx.Logger().Info(fmt.Sprintf("GetNode(Address = %v) \n", address.String()))
	n = k.posKeeper.Validator(ctx, address)
	if n == nil {
		return n, false
	}
	return n, true
}

// self node is needed to verify that self node is part of a session
func (k Keeper) GetSelfNode(ctx sdk.Context) (node exported.ValidatorI, er sdk.Error) {
	// get the Keybase addr list
	keypairs, err := (k.Keybase).GetCoinbase()
	if err != nil {
		er = pc.NewKeybaseError(pc.ModuleName, err)
		return nil, er
	}
	// get the node from the world state
	self, found := k.GetNode(ctx, sdk.Address(keypairs.GetAddress()))
	if !found {
		er = pc.NewSelfNotFoundError(pc.ModuleName)
		return nil, er
	}
	return self, nil
}

// award coins to nodes for relays completed
func (k Keeper) AwardCoinsForRelays(ctx sdk.Context, relays int64, toAddr sdk.Address) {
	ctx.Logger().Info(fmt.Sprintf("AwardsCoinsForRelays(relays = %v, toAddr = %v) \n", relays, toAddr.String()))
	k.posKeeper.AwardCoinsTo(ctx, sdk.NewInt(relays), toAddr)
}
