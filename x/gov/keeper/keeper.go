package keeper

import (
	"fmt"

	"github.com/pokt-network/pocket-core/codec"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/gov/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the global paramstore
type Keeper struct {
	cdc        *codec.Codec
	key        sdk.StoreKey
	tkey       sdk.StoreKey
	codespace  sdk.CodespaceType
	paramstore sdk.Subspace
	AuthKeeper types.AuthKeeper
	spaces     map[string]sdk.Subspace
}

// NewKeeper constructs a params keeper
func NewKeeper(cdc *codec.Codec, key *sdk.KVStoreKey, tkey *sdk.TransientStoreKey, codespace sdk.CodespaceType, authKeeper types.AuthKeeper, subspaces ...sdk.Subspace) (k Keeper) {
	k = Keeper{
		cdc:        cdc,
		key:        key,
		tkey:       tkey,
		codespace:  codespace,
		AuthKeeper: authKeeper,
		spaces:     make(map[string]sdk.Subspace),
	}
	k.paramstore = sdk.NewSubspace(types.ModuleName).WithKeyTable(types.ParamKeyTable())
	k.paramstore.SetCodec(k.cdc)
	subspaces = append(subspaces, k.paramstore)
	k.AddSubspaces(subspaces...)
	return k
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Ctx) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) UpgradeCodec(ctx sdk.Ctx) {
	if ctx.IsOnUpgradeHeight() {
		k.ConvertState(ctx)
	}
}

func (k Keeper) ConvertState(ctx sdk.Ctx) {
	k.cdc.SetUpgradeOverride(false)
	params := k.GetParams(ctx)
	k.cdc.SetUpgradeOverride(true)
	k.SetParams(ctx, params)
	k.cdc.DisableUpgradeOverride()
}
