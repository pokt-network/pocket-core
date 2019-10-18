package blockchain

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/pokt-network/pocket-core/types"
)

// NodeKeeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type NodeKeeper struct {
	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the blockchain NodeKeeper
func NewKeeper(coinKeeper bank.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec) NodeKeeper {
	return NodeKeeper{
		storeKey:   storeKey,
		cdc:        cdc,
	}
}

// Gets the entire Node metadata struct for an address
func (k NodeKeeper) GetNode(ctx sdk.Context, address string) types.Node {
	store := ctx.KVStore(k.storeKey)
	if !k.ContainsNode(ctx, address) {
		return types.Node{} // todo handle standard error logic
	}
	structEncoding := store.Get([]byte(address))
	var n types.Node
	k.cdc.MustUnmarshalBinaryBare(structEncoding, &n)
	return n
}

// Sets the entire node for an address
func (k NodeKeeper) SetNode(ctx sdk.Context, address string, n types.Node) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(address), k.cdc.MustMarshalBinaryBare(n))
}

// Deletes the entire metadata struct for an address
func (k NodeKeeper) DeleteNode(ctx sdk.Context, address string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete([]byte(address))
}

// Check if the name is present in the store or not
func (k NodeKeeper) ContainsNode(ctx sdk.Context, address string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(address))
}

func (k NodeKeeper) GetAllNodes(ctx sdk.Context) ([]types.Application, sdk.Error) {
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
func (k NodeKeeper) GetNodesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, nil)
}
