package keeper

import (
	"github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/nodes/exported"
	pc "github.com/pokt-network/pocket-core/x/pocketcore/types"
)

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

func (k Keeper) GetSelfPrivKey(ctx sdk.Ctx) (crypto.PrivateKey, sdk.Error) {
	// get the private key from the private validator file
	pk, er := k.GetPKFromFile(ctx)
	if er != nil {
		return nil, pc.NewKeybaseError(pc.ModuleName, er)
	}
	return pk, nil
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
func (k Keeper) AwardCoinsForRelays(ctx sdk.Ctx, relays int64, toAddr sdk.Address) sdk.BigInt {
	return k.posKeeper.RewardForRelays(ctx, sdk.NewInt(relays), toAddr)
}

// "BurnCoinsForChallenges" - Executes the burn for challenge function in the nodes module
func (k Keeper) BurnCoinsForChallenges(ctx sdk.Ctx, relays int64, toAddr sdk.Address) {
	k.posKeeper.BurnForChallenge(ctx, sdk.NewInt(relays), toAddr)
}
