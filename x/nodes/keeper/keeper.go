package keeper

import (
	"fmt"
	log2 "log"

	"github.com/pokt-network/pocket-core/codec"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Implements ValidatorSet interface
var _ types.ValidatorSet = Keeper{}

// Keeper of the staking store
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           *codec.Codec
	AccountKeeper types.AuthKeeper
	PocketKeeper  types.PocketKeeper // todo combine all modules
	Paramstore    sdk.Subspace
	// codespace
	codespace sdk.CodespaceType
	// Cache
	validatorCache *types.Cache
}

// NewKeeper creates a new staking Keeper instance
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, accountKeeper types.AuthKeeper,
	paramstore sdk.Subspace, codespace sdk.CodespaceType) Keeper {
	// ensure staked module accounts are set
	if addr := accountKeeper.GetModuleAddress(types.StakedPoolName); addr == nil {
		log2.Fatal(fmt.Errorf("%s module account has not been set", types.StakedPoolName))
	}
	cache, err := types.New(int(types.ValidatorCacheSize))
	if err != nil {
		log2.Fatal(fmt.Errorf("%d is an invalid size", types.ValidatorCacheSize))
	}

	return Keeper{
		storeKey:       key,
		cdc:            cdc,
		AccountKeeper:  accountKeeper,
		Paramstore:     paramstore.WithKeyTable(ParamKeyTable()),
		codespace:      codespace,
		validatorCache: cache,
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
