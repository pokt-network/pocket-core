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
	storeKey       sdk.StoreKey
	cdc            *codec.Codec
	AccountsKeeper types.AuthKeeper
	POSKeeper      types.PosKeeper
	Paramstore     sdk.Subspace
	// codespace
	codespace sdk.CodespaceType
	// Cache
	ApplicationCache *sdk.Cache
}

// NewKeeper creates a new staking Keeper instance
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, posKeeper types.PosKeeper, supplyKeeper types.AuthKeeper,
	paramstore sdk.Subspace, codespace sdk.CodespaceType) Keeper {

	// ensure staked module accounts are set
	if addr := supplyKeeper.GetModuleAddress(types.StakedPoolName); addr == nil {
		log2.Fatal(fmt.Errorf("%s module account has not been set", types.StakedPoolName))
	}
	cache := sdk.NewCache(int(types.ApplicationCacheSize))

	return Keeper{
		storeKey:         key,
		cdc:              cdc,
		AccountsKeeper:   supplyKeeper,
		POSKeeper:        posKeeper,
		Paramstore:       paramstore.WithKeyTable(ParamKeyTable()),
		codespace:        codespace,
		ApplicationCache: cache,
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
