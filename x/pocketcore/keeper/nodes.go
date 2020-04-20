package keeper

import (
	"github.com/pokt-network/pocket-core/x/nodes/exported"
	pc "github.com/pokt-network/pocket-core/x/pocketcore/types"
	sdk "github.com/pokt-network/posmint/types"
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

// "GetSelfNode" - Gets self node (private val key) from the world state
func (k Keeper) GetSelfNode(ctx sdk.Ctx) (node exported.ValidatorI, er sdk.Error) {
	// get the Keybase address list
	kp, err := k.GetPKFromFile(ctx)
	if err != nil {
		er = pc.NewKeybaseError(pc.ModuleName, err)
		return nil, er
	}
	// get the node from the world state
	self, found := k.GetNode(ctx, sdk.Address(kp.PublicKey().Address()))
	if !found {
		er = pc.NewSelfNotFoundError(pc.ModuleName)
		return nil, er
	}
	return self, nil
}

// "AwardCoinsForRelays" - Award coins to nodes for relays completed using the nodes keeper
func (k Keeper) AwardCoinsForRelays(ctx sdk.Ctx, relays int64, toAddr sdk.Address) {
	k.posKeeper.RewardForRelays(ctx, sdk.NewInt(relays), toAddr)
}

// "BurnCoinsForChallenges" - Executes the burn for challenge function in the nodes module
func (k Keeper) BurnCoinsForChallenges(ctx sdk.Ctx, relays int64, toAddr sdk.Address) {
	k.posKeeper.BurnForChallenge(ctx, sdk.NewInt(relays), toAddr)
}
