package keeper

import (
	"github.com/pokt-network/pocket-core/x/nodes/exported"
	pc "github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
)

// get all nodes from the world state
func (k Keeper) GetAllNodes(ctx sdk.Context) []exported.ValidatorI {
	return k.posKeeper.AllValidators(ctx)
}

// get a node from the world state
func (k Keeper) GetNode(ctx sdk.Context, address sdk.Address) (n exported.ValidatorI, found bool) {
	n = k.posKeeper.Validator(ctx, address)
	if n == nil {
		return n, false
	}
	return n, true
}

// get the node from the public key string
func (k Keeper) GetNodeFromPublicKey(ctx sdk.Context, pubKey string) (node exported.ValidatorI, found bool) {
	// get the node at the session context
	pk, err := crypto.NewPublicKey(pubKey)
	if err != nil {
		return nil, false
	}
	return k.GetNode(ctx, sdk.Address(pk.Address()))
}

//// todo create store in pos module for efficiency
//func (k Keeper) GetAllNodesForChain(ctx sdk.Context, chain string) (nodes []exported.ValidatorI) {
//	return nil
//}

// self node is needed to verify that self node is part of a session
func (k Keeper) GetSelfNode(ctx sdk.Context) (node exported.ValidatorI, er sdk.Error) {
	// get the Keybase addr list
	keypairs, err := (k.Keybase).GetCoinbase()
	if err != nil {
		return nil, pc.NewKeybaseError(pc.ModuleName, err)
	}
	// get the node from the world state
	self, found := k.GetNode(ctx, sdk.Address(keypairs.GetAddress()))
	if !found {
		return nil, pc.NewSelfNotFoundError(pc.ModuleName)
	}
	return self, nil
}

// award coins to nodes for relays completed
func (k Keeper) AwardCoinsForRelays(ctx sdk.Context, relays int64, toAddr sdk.Address) {
	k.posKeeper.AwardCoinsTo(ctx, sdk.NewInt(relays), toAddr)
}
