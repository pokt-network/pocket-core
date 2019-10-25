package pocketcore

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pokt-network/pocket-core/types"
)

// nodeKeeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type nodeKeeper struct {
	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc      *codec.Codec // The wire codec for binary encoding/decoding.
}

// newNodeKeeper creates new instances of the blockchain nodeKeeper
func newNodeKeeper(storeKey sdk.StoreKey, cdc *codec.Codec) nodeKeeper {
	return nodeKeeper{
		storeKey: storeKey,
		cdc:      cdc,
	}
}

// Gets the entire Node metadata struct for an address
func (k nodeKeeper) GetNode(ctx sdk.Context, address string) types.Node {
	store := ctx.KVStore(k.storeKey)
	if !k.ContainsNode(ctx, address) {
		return types.Node{} // todo handle standard error logic
	}
	structEncoding := store.Get([]byte(address))
	var n types.Node
	k.cdc.MustUnmarshalBinaryBare(structEncoding, &n)
	return n
}

// Gets the entire Node metadata struct for an address at a certain block height
func (k nodeKeeper) GetNodeAtHeight(ctx sdk.Context, address string, height int64) types.Application {
	atHeight := ctx.WithBlockHeight(height)
	store := atHeight.KVStore(k.storeKey)
	if !k.ContainsNode(atHeight, address) {
		return types.Application{} // todo handle standard error logic
	}
	structEncoding := store.Get([]byte(address))
	var n types.Application
	k.cdc.MustUnmarshalBinaryBare(structEncoding, &n)
	return n
}

// Sets the entire node for an address
func (k nodeKeeper) SetNode(ctx sdk.Context, address string, n types.Node) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(address), k.cdc.MustMarshalBinaryBare(n))
}

// Deletes the entire metadata struct for an address
func (k nodeKeeper) DeleteNode(ctx sdk.Context, address string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete([]byte(address))
}

// Check if the name is present in the store or not
func (k nodeKeeper) ContainsNode(ctx sdk.Context, address string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(address))
}

func (k nodeKeeper) GetAllNodes(ctx sdk.Context) ([]types.Application, sdk.Error) {
	var nodeList []types.Application
	node := new(types.Application)

	iterator := k.GetNodesIterator(ctx)

	for ; iterator.Valid(); iterator.Next() {
		err := k.cdc.UnmarshalBinaryBare(iterator.Key(), node)
		if err != nil {
			return nil, types.NewError(types.CODENODEUNMARSHALUNSUCCESSFUL, err.Error())
		}
		nodeList = append(nodeList, *node)
	}
	return nodeList, nil
}

// Get an iterator over all names in which the keys are the names and the values are the whois
func (k nodeKeeper) GetNodesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, nil)
}
