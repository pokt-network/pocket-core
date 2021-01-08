package keeper

import (
	"fmt"
	"github.com/pokt-network/pocket-core/codec"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/auth/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the supply store
type Keeper struct {
	Cdc       *codec.Codec
	storeKey  sdk.StoreKey
	subspace  sdk.Subspace
	permAddrs map[string]types.PermissionsForAddress
}

// NewKeeper creates a new Keeper instance
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, subspace sdk.Subspace, maccPerms map[string][]string) Keeper {
	// set the addresses
	permAddrs := make(map[string]types.PermissionsForAddress)
	for name, perms := range maccPerms {
		permAddrs[name] = types.NewPermissionsForAddress(name, perms)
	}

	return Keeper{
		Cdc:       cdc,
		storeKey:  key,
		subspace:  subspace.WithKeyTable(types.ParamKeyTable()),
		permAddrs: permAddrs,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Ctx) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Codespace returns the keeper's codespace.
func (k Keeper) Codespace() sdk.CodespaceType {
	return types.DefaultCodespace
}

func (k Keeper) UpgradeCodec(ctx sdk.Ctx) {
	if ctx.IsOnUpgradeHeight() {
		k.ConvertState(ctx)
	}
}

func (k Keeper) ConvertState(ctx sdk.Ctx) {
	k.Cdc.SetUpgradeOverride(false)
	params := k.GetParams(ctx)
	accounts := k.GetAllAccounts(ctx)
	supply := k.GetSupply(ctx)
	k.Cdc.SetUpgradeOverride(true)
	k.SetParams(ctx, params)
	k.SetAccounts(ctx, accounts)
	k.SetSupply(ctx, supply)
	k.Cdc.DisableUpgradeOverride()
}
