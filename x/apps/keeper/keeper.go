package keeper

import (
	"fmt"
	"github.com/pokt-network/pocket-core/codec"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/apps/types"
	"github.com/tendermint/tendermint/libs/log"
	log2 "log"
)

// Implements ApplicationSet interface
var _ types.ApplicationSet = Keeper{}

// Keeper of the staking store
type Keeper struct {
	storeKey      sdk.StoreKey
	Cdc           *codec.Codec
	AccountKeeper types.AuthKeeper
	POSKeeper     types.PosKeeper
	PocketKeeper  types.PocketKeeper
	Paramstore    sdk.Subspace
	// codespace
	codespace sdk.CodespaceType
	// Cache
	ApplicationCache *sdk.Cache
}

// NewKeeper creates a new staking Keeper instance
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, posKeeper types.PosKeeper, supplyKeeper types.AuthKeeper, pocketKeeper types.PocketKeeper,
	paramstore sdk.Subspace, codespace sdk.CodespaceType) Keeper {

	// ensure staked module accounts are set
	if addr := supplyKeeper.GetModuleAddress(types.StakedPoolName); addr == nil {
		log2.Fatal(fmt.Errorf("%s module account has not been set", types.StakedPoolName))
	}
	cache := sdk.NewCache(int(types.ApplicationCacheSize))

	return Keeper{
		storeKey:         key,
		PocketKeeper:     pocketKeeper,
		AccountKeeper:    supplyKeeper,
		POSKeeper:        posKeeper,
		Paramstore:       paramstore.WithKeyTable(ParamKeyTable()),
		codespace:        codespace,
		ApplicationCache: cache,
		Cdc:              cdc,
	}
}

// Logger - returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Ctx) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Codespace - Retrieve the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

func (k Keeper) UpgradeCodec(ctx sdk.Ctx) {
	if ctx.IsOnUpgradeHeight() {
		k.ConvertState(ctx)
	}
}

func (k Keeper) ConvertState(ctx sdk.Ctx) {
	k.Cdc.SetUpgradeOverride(false)
	params := k.GetParams(ctx)
	apps := k.GetAllApplications(ctx)
	k.Cdc.SetUpgradeOverride(true)
	k.SetParams(ctx, params)
	k.SetApplications(ctx, apps)
	k.Cdc.DisableUpgradeOverride()
}
