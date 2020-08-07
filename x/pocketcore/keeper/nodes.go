package keeper

import (
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/nodes/exported"
	pc "github.com/pokt-network/pocket-core/x/pocketcore/types"
)

// "GetAllNodes" - Gets all of the nodes in the state storage
func (k Keeper) GetAllNodes(ctx sdk.Ctx) []exported.ValidatorI {
	validators := k.posKeeper.GetStakedValidators(ctx)
	return validators
}

// "GetNode" - Gets a node from the state storage
func (k Keeper) GetNode(ctx sdk.Ctx, address sdk.Address) (n exported.ValidatorI, found bool) {
	n = k.posKeeper.Validator(ctx, address)
	if n == nil {
		return n, false
	}
	return n, true
}

func (k Keeper) GetSelfAddress(ctx sdk.Ctx) sdk.Address {
	kp, err := k.GetPKFromFile(ctx)
	if err != nil {
		ctx.Logger().Error("Unable to retrieve selfAddress: " + err.Error())
		return nil
	}
	return sdk.Address(kp.PublicKey().Address())
}

// "GetSelfNode" - Gets self node (private val key) from the world state
func (k Keeper) GetSelfNode(ctx sdk.Ctx) (node exported.ValidatorI, er sdk.Error) {
	// get the node from the world state
	self, found := k.GetNode(ctx, k.GetSelfAddress(ctx))
	if !found {
		er = pc.NewSelfNotFoundError(pc.ModuleName)
		return nil, er
	}
	return self, nil
}

// "AwardCoinsForRelays" - Award coins to nodes for relays completed using the nodes keeper
func (k Keeper) AwardCoinsForRelays(ctx sdk.Ctx, relays int64, toAddr sdk.Address) sdk.Int {
	return k.posKeeper.RewardForRelays(ctx, sdk.NewInt(relays), toAddr)
}

// "BurnCoinsForChallenges" - Executes the burn for challenge function in the nodes module
func (k Keeper) BurnCoinsForChallenges(ctx sdk.Ctx, relays int64, toAddr sdk.Address) {
	k.posKeeper.BurnForChallenge(ctx, sdk.NewInt(relays), toAddr)
}
