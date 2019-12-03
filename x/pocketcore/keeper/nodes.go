package keeper

import (
	"github.com/pokt-network/pocket-core/x/nodes/exported"
	pc "github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
)

func (k Keeper) GetAllNodes(ctx sdk.Context) []exported.ValidatorI {
	return k.posKeeper.GetAllValidators(ctx)
}

func (k Keeper) GetNodeFromPublicKey(ctx sdk.Context, pubKey string) (node exported.ValidatorI, found bool) {
	// get the node at the session context
	pk, err := crypto.NewPublicKey(pubKey)
	if err != nil {
		return nil, false
	}
	return k.GetNode(ctx, pk.Address())
}

func (k Keeper) GetNode(ctx sdk.Context, address sdk.ValAddress) (node exported.ValidatorI, found bool) {
	return k.posKeeper.GetValidator(ctx, address)
}

func (k Keeper) GetNodeChains(ctx sdk.Context, address sdk.ValAddress) (chains map[string]struct{}, found bool) {
	node, found := k.posKeeper.GetValidator(ctx, address)
	if !found {
		return nil, false
	}
	return node.GetChains(), true
}

func (k Keeper) GetNodeServiceURL(ctx sdk.Context, address sdk.ValAddress) (serviceURL string, found bool) {
	node, found := k.posKeeper.GetValidator(ctx, address)
	if !found {
		return "", false
	}
	return node.GetServiceURL(), true
}

// todo create store in pos module for efficiency
func (k Keeper) GetAllNodesForChain(ctx sdk.Context, chain string) (nodes []exported.ValidatorI) {
	return nil
}

// self node is needed to verify that self node is part of a session
func (k Keeper) GetSelfNode(ctx sdk.Context) (node exported.ValidatorI, er sdk.Error) {
	// get the keybase addr list
	keypairs, err := k.keybase.List()
	if err != nil || len(keypairs) < 1 {
		return nil, pc.NewKeybaseError(pc.ModuleName, err)
	}
	// get the node from the world state
	self, found := k.GetNode(ctx, sdk.ValAddress(keypairs[0].GetAddress())) // todo need to verify that this is the validator key we want
	if !found {
		return nil, pc.NewSelfNotFoundError(pc.ModuleName)
	}
	return self, nil
}

func (k Keeper) GetHostedBlockchains() pc.HostedBlockchains {
	return k.hostedBlockchains
}

func (k Keeper) AwardCoinsForRelays(ctx sdk.Context, relays int64, toAddr sdk.ValAddress) {
	k.posKeeper.AwardCoinsTo(ctx, sdk.NewInt(relays), toAddr)
}
